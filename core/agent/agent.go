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
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
	"path"

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
	env := make([]string, 0)
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

// NOTE: This function is a copy of dedupEnvCase in go/exec module
func dedupEnv(env []string) []string {
	out := make([]string, 0, len(env))
	saw := make(map[string]int, len(env)) // key => index into out
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			out = append(out, kv)
			continue
		}
		k := kv[:eq]
		if dupIdx, isDup := saw[k]; isDup {
			out[dupIdx] = kv
			continue
		}
		saw[k] = len(out)
		out = append(out, kv)
	}
	return out
}

// Transform the string array representing an environment
// to a list of values for the redis.HSet variadic function
func exportToRedis(env []string) []interface{} {
	out := make([]interface{}, len(env)*2)
	for i, kv := range env {
		j := i * 2
		eq := strings.Index(kv, "=")
		if eq < 0 {
			out[j] = kv
			out[j+1] = ""
		} else {
			out[j] = kv[:eq]
			out[j+1] = kv[eq+1:]
		}
	}
	return out
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
	environment := dedupEnv(prepareActionEnvironment())

	for _, step := range actionSteps {
		// Create a pipe to read control commands from action steps
		comReadFd, comWriteFd, _ := os.Pipe()

		cmd := exec.Command(step)
		cmd.Env = append(append(os.Environ(), environment...),
			"AGENT_COMFD=3", // 3 is the additional FD number where the action step can write its commands for us
			"AGENT_TASK_ID=" + task.ID,
			"AGENT_TASK_ACTION=" + task.Action,
		)
		cmd.Stdin = strings.NewReader(task.Data)
		cmd.ExtraFiles = []*os.File{comWriteFd}

		stdoutReader, _ := cmd.StdoutPipe()
		// Copy the command stderr to our stderr stream with a pipe tee
		stderrReader, _ := cmd.StderrPipe()
		stderrReader = io.NopCloser(io.TeeReader(stderrReader, os.Stderr))

		go func() {
			bytes, _ := io.ReadAll(stderrReader)
			actionError = string(bytes)
		}()

		go func() {
			bytes, _ := io.ReadAll(stdoutReader)
			actionOutput = string(bytes)
		}()

		go func() {
			// Read the commands FD line by line, parse and execute the command
			rdr := csv.NewReader(comReadFd)
			rdr.Comma = ' '
			rdr.FieldsPerRecord = 3
			for {
				record, err := rdr.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("[ERROR] Parse error: %v", err)
					continue
				}
				switch cmd := record[0]; cmd {
				case "set-env":
					environment = append(environment, record[1]+"="+record[2])
				default:
					log.Printf("[ERROR] Unknown command %s", cmd)
				}
			}
		}()

		log.Printf("%s/task/%s: %s/%s is starting", agentPrefix, task.ID, task.Action, path.Base(step))
		if err := cmd.Start(); err != nil {
			log.Printf("[ERROR] Action %s startup error at step %s: %v", task.Action, step, err)
			break
		}

		// standard FD are closed by Start. Our extra FD for commands must be closed manually
		// otherwise it blocks our thread.
		comWriteFd.Close()

		// It is safe to Wait() after the select{} finishes to consume the pipes input.
		if err := cmd.Wait(); err != nil {
			exitCode = cmd.ProcessState.ExitCode()
			log.Printf("[ERROR] Action %s terminated with errors at step %s: %v", task.Action, step, err)
			break
		}

		environment = dedupEnv(environment)
	}

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		// Persist the environment
		if exitCode == 0 && len(environment) > 0 {
			pipe.HSet(ctx, agentPrefix+"/environment", exportToRedis(environment)...)
		}
		// Publish the action response
		pipe.Set(ctx, outputKey, actionOutput, taskExpireDuration)
		pipe.Set(ctx, errorKey, actionError, taskExpireDuration)
		pipe.Set(ctx, exitCodeKey, exitCode, taskExpireDuration)
		pipe.Publish(ctx, progressChannel, `{"status":"finished","progress":100}`)
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Redis command failed: ", err)
	}
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
	// The default user name can be overridden by the REDIS_USER environment variable
	redisUsername := ""
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword != "" {
		redisUsername = os.Getenv("REDIS_USER")
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

		result, err := rdb.BRPop(ctx, pollingDuration, queueName).Result()
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
