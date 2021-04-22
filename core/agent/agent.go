/*
 * Copyright (C) 2021 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var agentPrefix = os.Args[1]
var actionPaths = os.Args[2:]
var rdb *redis.Client

var pollingDuration = 5 * time.Second
var taskExpireDuration = 24 * time.Hour

func getActionSteps(actionName string) ([]string, bool) {
	dirFound := false
	actionSteps := make(map[string]string)

	// Squash the action dirs in a single list
	for _, path := range actionPaths {
		actionPath := path + "/" + actionName
		entries, err := os.ReadDir(actionPath)
		if err != nil {
			continue
		}
		dirFound = true
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			actionSteps[entry.Name()] = actionPath + "/" + entry.Name()
		}
	}

	actionStepsKeys := make([]string, 0, 8)
	for k, _ := range actionSteps {
		actionStepsKeys = append(actionStepsKeys, k)
	}

	// Executable steps under the actionPaths are launched sequentially
	// in alphabetical order:
	sort.Strings(actionStepsKeys)

	actionStepsValues := make([]string, 0, 8)
	for _, k := range actionStepsKeys {
		actionStepsValues = append(actionStepsValues, actionSteps[k])
	}

	return actionStepsValues, dirFound
}

func prepareActionEnvironment() []string {
	env := os.Environ()
	key := agentPrefix + "/environment"
	redisHash, err := rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return env
	}
	if err != nil {
		log.Printf("[ERROR] Could not fetch Redis key %s: %v", key, err)
		return env
	}
	for key, value := range redisHash {
		env = append(env, key+"="+value)
	}
	return env
}

func runAction(task *models.Task) {

	// Redis key names where the action response is stored:
	progressChannel := "progress/task/" + task.ID
	outputKey := agentPrefix + "/task/" + task.ID + "/output"
	errorKey := agentPrefix + "/task/" + task.ID + "/error"
	exitCodeKey := agentPrefix + "/task/" + task.ID + "/exit_code"

	rdb.Publish(ctx, progressChannel, `{"status":"pending","progress":0}`)

	// Action response payloads
	var actionOutput string = ""
	var actionError string = ""
	var exitCode int = 0

	actionSteps, actionDefined := getActionSteps(task.Action)

	if !actionDefined {
		// If the action is not defined our exit code is returned
		// More info in man systemd.exec and from `systemd-analyze exit-status` output
		actionOutput = ""
		actionError = ""
		exitCode = 8
		log.Printf("[ERROR] Action %s is not defined", task.Action)
	}

	// Get additional environment variables from Redis DB and
	// other runtime sources
	environment := prepareActionEnvironment()

	for _, step := range actionSteps {
		cmd := exec.Command(step)
		cmd.Env = environment
		cmd.Stdin = bytes.NewBufferString(task.Data)

		stderrReader, _ := cmd.StderrPipe()
		stdoutReader, _ := cmd.StdoutPipe()

		stderrChannel := make(chan string)
		go func() {
			// Copy the command stderr to our stderr stream with a pipe tee
			r := io.TeeReader(stderrReader, os.Stderr)
			qb, _ := io.ReadAll(r)
			stderrChannel <- string(qb)
			close(stderrChannel)
		}()

		stdoutChannel := make(chan string)
		go func() {
			qb, _ := io.ReadAll(stdoutReader)
			stdoutChannel <- string(qb)
			close(stdoutChannel)
		}()

		if err := cmd.Start(); err != nil {
			log.Printf("[ERROR] Action %s startup error at step %s: %v", task.Action, step, err)
			break
		}

		// Block until stderr and stdout are closed
		select {
		case buf := <-stderrChannel:
			actionError = buf
		case buf := <-stdoutChannel:
			actionOutput = buf
		}

		// It is safe to Wait() after the select{} finishes to consume the pipes input.
		if err := cmd.Wait(); err != nil {
			exitCode = cmd.ProcessState.ExitCode()
			log.Printf("[ERROR] Action %s terminated with errors at step %s: %v", task.Action, step, err)
			break
		}
	}

	// Publish the action response
	rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, outputKey, actionOutput, taskExpireDuration)
		pipe.Set(ctx, errorKey, actionError, taskExpireDuration)
		pipe.Set(ctx, exitCodeKey, exitCode, taskExpireDuration)
		pipe.Publish(ctx, progressChannel, `{"status":"finished","progress":100}`)
		return nil
	})
}

func main() {
	log.SetFlags(0)
	if agentPrefix == "" {
		log.Fatal("[FATAL] The agent prefix argument is not set")
	}
	if len(actionPaths) == 0 {
		log.Fatal("[FATAL] Action path command arguments are not set")
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "127.0.0.1:6379"
	}

	// If we have a REDIS_PASSWORD the default redis username is the agentPrefix string
	// The default user name can be overridden by the REDIS_USERNAME environment variable
	redisUsername := ""
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword != "" {
		redisUsername = os.Getenv("REDIS_USERNAME")
		if redisUsername == "" {
			redisUsername = agentPrefix
		}
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Username: redisUsername,
		Password: redisPassword,
		DB:       0,
	})

	// TODO notify systemd unit startup

	queueName := agentPrefix + "/tasks"

	for {
		var task models.Task
		rdb.Do(ctx, "CLIENT", "SETNAME", agentPrefix).Result()

		result, err := rdb.BLPop(ctx, pollingDuration, queueName).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Print("[ERROR] Task queue pop error: ", err)
			time.Sleep(pollingDuration)
			continue
		}

		if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
			log.Print("[ERROR] Task ignored for decoding error: ", err)
			continue
		}

		// run the Action required by the Task payload
		go runAction(&task)
	}
}
