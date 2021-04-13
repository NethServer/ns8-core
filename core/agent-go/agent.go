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

	"github.com/NethServer/ns8-scratchpad/api-server/models"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var agentPrefix = os.Getenv("AGENT_PREFIX")
var actionPaths = os.Args[1:]
var rdb *redis.Client

var pollingDuration = 5 * time.Second
var taskExpireDuration = 24 * time.Hour

func getActionSteps(actionName string) ([]string, bool) {
	dirFound := false
	actionSteps := make([]string, 0, 4)

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
			actionSteps = append(actionSteps, actionPath+"/"+entry.Name())
		}
	}

	// Executable steps under the actionPaths are launched sequentially
	// in alphabetical order:
	sort.Strings(actionSteps)

	return actionSteps, dirFound
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

	for _, step := range actionSteps {
		cmd := exec.Command(step)
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

		// After the select{} async procedures already consumed their input from pipes. It is safe to Wait(). 
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
		log.Fatal("[FATAL] AGENT_PREFIX is not set in the environment")
	}
	if len(actionPaths) == 0 {
		log.Fatal("[FATAL] Action path command arguments are not set")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// TODO notify systemd unit startup

	queueName := agentPrefix + "/tasks"

	for {
		var task models.Task

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
