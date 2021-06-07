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
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/NethServer/ns8-scratchpad/core/agent/action"
	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/go-redis/redis/v8"
)

// Reference: https://www.man7.org/linux/man-pages/man3/sd-daemon.3.html
const (
	SD_EMERG   = "<0>" /* system is unusable */
	SD_ALERT   = "<1>" /* action must be taken immediately */
	SD_CRIT    = "<2>" /* critical conditions */
	SD_ERR     = "<3>" /* error conditions */
	SD_WARNING = "<4>" /* warning conditions */
	SD_NOTICE  = "<5>" /* normal but significant condition */
	SD_INFO    = "<6>" /* informational */
	SD_DEBUG   = "<7>" /* debug-level messages */
)

var ctx = context.Background()
var agentPrefix = os.Args[1]
var actionPaths = os.Args[2:]
var rdb *redis.Client

var pollingDuration = 5 * time.Second
var taskExpireDuration = 24 * time.Hour

func prepareActionEnvironment() []string {
	env := make([]string, 0)
	key := agentPrefix + "/environment"
	redisHash, err := rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return env
	}
	if err != nil {
		log.Printf(SD_ERR+"Could not fetch Redis key %s: %v", key, err)
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

// Persist the environment to the current working directory
// Systemd unit file can pickup it with EnvironmentFile option
func dumpToFile(env []string) {
	path, _ := os.Getwd()
	env = dedupEnv(env)
	f, err := os.Create("./environment")
	if err != nil {
		log.Printf(SD_ERR+"Can't write %s/environment file: %s", path, err)
		return
	}
	for _, line := range env {
		f.WriteString(line + "\n")
	}
	f.Close()
	log.Printf("Wrote %s/environment file", path)
}

func publishStatus(client redis.Cmdable, progressChannel string, actionDescriptor action.Descriptor) {
	err := client.Publish(ctx, progressChannel, actionDescriptor.ToJSON()).Err()
	if err != nil {
		log.Printf(SD_ERR+"Failed to publish the action status on channel %s", progressChannel)
	}
}


func runListActions(task *models.Task) {
	// Redis key names where the action response is stored:
	progressChannel := "progress/" + agentPrefix + "/" + task.ID
	outputKey := agentPrefix + "/task/" + task.ID + "/output"
	errorKey := agentPrefix + "/task/" + task.ID + "/error"
	exitCodeKey := agentPrefix + "/task/" + task.ID + "/exit_code"

	actionDescriptor := action.CreateBuiltin(task.Action)
	publishStatus(rdb, progressChannel, actionDescriptor)

	actionDescriptor.Status = "running"
	publishStatus(rdb, progressChannel, actionDescriptor)

	actionOutput, _ := json.Marshal(action.ListActions(actionPaths))
	actionError := ""
	exitCode := 0

	actionDescriptor.SetProgressAtStep(0, 100)
	actionDescriptor.Status = "completed"

	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// Publish the action response
		pipe.Set(ctx, outputKey, actionOutput, taskExpireDuration)
		pipe.Set(ctx, errorKey, actionError, taskExpireDuration)
		pipe.Set(ctx, exitCodeKey, exitCode, taskExpireDuration)
		publishStatus(pipe, progressChannel, actionDescriptor)
		return nil
	})
	if err != nil {
		log.Print(SD_ERR+"Redis command failed: ", err)
	}
}

func runAction(task *models.Task) {

	// Redis key names where the action response is stored:
	progressChannel := "progress/" + agentPrefix + "/task/" + task.ID
	outputKey := agentPrefix + "/task/" + task.ID + "/output"
	errorKey := agentPrefix + "/task/" + task.ID + "/error"
	exitCodeKey := agentPrefix + "/task/" + task.ID + "/exit_code"

	// Action response payloads
	var actionOutput string = ""
	var actionError string = ""
	var exitCode int = 0

	actionDescriptor := action.Create(task.Action, actionPaths)

	publishStatus(rdb, progressChannel, actionDescriptor)

	if len(actionDescriptor.Steps) == 0 {
		// If the action is not defined our exit code is returned
		// More info in man systemd.exec and from `systemd-analyze exit-status` output
		actionOutput = ""
		actionError = fmt.Sprintf("Action %s is not defined\n", task.Action)
		exitCode = 8
		actionDescriptor.Status = "aborted"
		log.Printf(SD_ERR+"Action %s is not defined", task.Action)
	}

	// Get additional environment variables from Redis DB and
	// other runtime sources
	environment := dedupEnv(prepareActionEnvironment())

	for stepIndex, step := range actionDescriptor.Steps {
		// Create a pipe to read control commands from action steps
		comReadFd, comWriteFd, _ := os.Pipe()

		cmd := exec.Command(step.Path)
		cmd.Env = append(append(os.Environ(), environment...),
			"AGENT_COMFD=3", // 3 is the additional FD number where the action step can write its commands for us
			"AGENT_TASK_ID="+task.ID,
			"AGENT_TASK_ACTION="+task.Action,
		)
		inputData, _ := json.Marshal(task.Data)
		cmd.Stdin = strings.NewReader(string(inputData))
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
			rdr.FieldsPerRecord = -1
			for {
				record, err := rdr.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf(SD_ERR+"Parse error: %v", err)
					continue
				}
				switch cmd := record[0]; cmd {
				case "set-env":
					environment = append(environment, record[1]+"="+record[2])
				case "dump-env":
					rdb.HSet(ctx, agentPrefix+"/environment", exportToRedis(environment)...).Result()
					dumpToFile(environment)
				case "set-status":
					if record[1] == "validation-failed" {
						actionDescriptor.Status = "validation-failed"
					} else {
						log.Printf(SD_ERR+"set-status command failed: unknown status \"%s\"", record[1])
					}
				case "set-weight":
					var weight int
					var err error
					weight, err = strconv.Atoi(record[2])
					if err != nil {
						log.Printf(SD_ERR+"set-weight command failed: %v", err)
						break
					}
					err = actionDescriptor.SetStepWeight(record[1], weight)
					if err != nil {
						log.Printf(SD_ERR+"set-weight command failed: %v", err)
						break
					}
					publishStatus(rdb, progressChannel, actionDescriptor)
				case "set-progress":
					var progress int
					var err error
					progress, err = strconv.Atoi(record[1])
					if err != nil {
						log.Printf(SD_ERR+"set-progress command failed: %v", err)
						break
					}
					err = actionDescriptor.SetProgressAtStep(stepIndex, progress)
					if err != nil {
						log.Printf(SD_ERR+"set-progress command failed: %v", err)
						break
					}
					publishStatus(rdb, progressChannel, actionDescriptor)
				default:
					log.Printf(SD_ERR+"Unknown command %s", cmd)
				}
			}
		}()

		actionDescriptor.Status = "running"

		log.Printf("%s/task/%s: %s/%s is starting", agentPrefix, task.ID, task.Action, step.Name)
		if err := cmd.Start(); err != nil {
			exitCode = 9
			actionError = fmt.Sprintf("Action %s startup error at step %s: %v", task.Action, step, err)
			actionDescriptor.Status = "aborted"
			log.Print(SD_ERR + actionError)
			break
		}

		// standard FDs are closed by Start. Our extra FD for commands must be closed manually
		// otherwise it blocks our thread.
		comWriteFd.Close()

		if err := cmd.Wait(); err != nil {
			exitCode = cmd.ProcessState.ExitCode()
			if actionDescriptor.Status == "running" {
				actionDescriptor.Status = "aborted"
			}
			log.Printf(SD_ERR+"Action %s %s at step %s: %v", task.Action, actionDescriptor.Status, step.Path, err)
			break
		} else if actionDescriptor.Status == "validation-failed" {
			exitCode = 10 // validation-error forces a non-zero exit code
			log.Printf(SD_ERR+"Action %s %s at step %s: %v", task.Action, actionDescriptor.Status, step.Path, err)
			break
		}

		actionDescriptor.SetProgressAtStep(stepIndex, 100)
		publishStatus(rdb, progressChannel, actionDescriptor)

		environment = dedupEnv(environment)
	}

	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if exitCode == 0 {
			if actionDescriptor.Status == "running" {
				actionDescriptor.Status = "completed"
			}
			// Persist the environment
			if len(environment) > 0 {
				pipe.HSet(ctx, agentPrefix+"/environment", exportToRedis(environment)...)
			} else {
				pipe.Del(ctx, agentPrefix+"/environment")
			}
		}
		// Publish the action response
		pipe.Set(ctx, outputKey, actionOutput, taskExpireDuration)
		pipe.Set(ctx, errorKey, actionError, taskExpireDuration)
		pipe.Set(ctx, exitCodeKey, exitCode, taskExpireDuration)
		publishStatus(pipe, progressChannel, actionDescriptor)
		return nil
	})
	if err != nil {
		log.Print(SD_ERR+"Redis command failed: ", err)
	} else {
		if exitCode == 0 {
			dumpToFile(environment)
		}
	}
}

func setClientNameCallback (ctx context.Context, cn *redis.Conn) error {
	return cn.ClientSetName(ctx, agentPrefix).Err()
}

func main() {
	log.SetFlags(0)
	if agentPrefix == "" {
		log.Fatal(SD_EMERG+"The agent prefix argument is not set")
		// exit(1) log.Fatal terminates the process with an exit code
	}
	if len(actionPaths) == 0 {
		log.Fatal(SD_EMERG+"Action path command arguments are not set")
		// exit(1) log.Fatal terminates the process with an exit code
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
		OnConnect: setClientNameCallback,
	})

	queueName := agentPrefix + "/tasks"

	for {
		var task models.Task

		// Pop the task from the agent tasks queue
		popResult, popErr := rdb.BRPop(ctx, pollingDuration, queueName).Result()
		if popErr == redis.Nil {
			continue
		} else if popErr != nil {
			log.Print(SD_ERR+"Task queue pop error: ", popErr)
			time.Sleep(pollingDuration)
			continue
		}

		if err := json.Unmarshal([]byte(popResult[1]), &task); err != nil {
			log.Print(SD_ERR+"Task ignored for decoding error: ", err)
			continue
		}

		// Store the task as context
		setErr := rdb.Set(ctx, agentPrefix + "/task/" + task.ID + "/context", popResult[1], taskExpireDuration).Err()
		if setErr != nil {
			log.Print(SD_ERR+"Context set error: ", setErr)
		}

		// run the Action required by the Task payload
		switch task.Action {
		case "list-actions":
			go runListActions(&task)
		default:
			go runAction(&task)
		}
	}
}
