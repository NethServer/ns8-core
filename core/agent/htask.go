/*
 * Copyright (C) 2022 Nethesis S.r.l.
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
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/NethServer/ns8-core/core/agent/models"
	"github.com/NethServer/ns8-core/core/agent/validation"
	"github.com/go-redis/redis/v8"
)

func runAction(rdb *redis.Client, actionCtx context.Context, task *models.Task) {

	// Redis key names where the action response is stored:
	progressChannel := "progress/" + agentPrefix + "/task/" + task.ID
	outputKey := agentPrefix + "/task/" + task.ID + "/output"
	errorKey := agentPrefix + "/task/" + task.ID + "/error"
	exitCodeKey := agentPrefix + "/task/" + task.ID + "/exit_code"

	// Action response payloads
	var actionOutput string = ""
	var actionError string = ""
	var exitCode int = 0

	lastStep := "{unknown}"

	actionDescriptor := models.CreateTaskProcessor(task.Action, actionPaths)

	publishStatus(rdb, progressChannel, actionDescriptor) // publish pending status

	if len(actionDescriptor.Steps) == 0 {
		// If the action is not defined our exit code is returned
		// More info in man systemd.exec and from `systemd-analyze exit-status` output
		actionOutput = ""
		actionError = fmt.Sprintf("Action %s is not defined\n", task.Action)
		exitCode = 8
		actionDescriptor.Status = "aborted"
		log.Printf(SD_ERR+"Action %s is not defined", task.Action)
	}

	// Read initial environment file contents
	environment := readStateFile()

	// Write the environment state to disk if any change occurs:
	var isStateWriteNeeded bool = false

	for stepIndex, step := range actionDescriptor.Steps {
		lastStep = step.Name

		// Sync environment variables from the filesystem. Reading the
		// environment file on every step allows to load changes produced
		// by child actions.
		environment = readStateFile()
		isStateWriteNeeded = false // reset the flag on each cycle

		// Special treatment for builtin validation steps
		if step.Name == models.STEP_VALIDATE_INPUT {
			errorList, validateFault := validation.ValidateGoStruct(step.Path, task.Data)
			if validateFault != nil {
				errorMessage := fmt.Sprintf("JSON Schema input validation aborted at step %s: %v\n", step.Path, validateFault)
				log.Printf(SD_ERR + errorMessage)
				actionError += errorMessage
				actionDescriptor.Status = "aborted"
				exitCode = 1
				break
			} else if len(errorList) > 0 {
				actionError += fmt.Sprintf("Validation errors: %v\n", errorList)
				errorsBuf, _ := validation.ToJSON(errorList)
				actionOutput += string(errorsBuf)
				exitCode = 10
				actionDescriptor.Status = "validation-failed"
				log.Printf(SD_ERR+"Action %s %s at step %s: %v", task.Action, actionDescriptor.Status, step.Path, errorList)
				break
			}
			continue

		} else if step.Name == models.STEP_VALIDATE_OUTPUT {
			errorList, validateFault := validation.ValidateJsonString(step.Path, []byte(actionOutput))
			if validateFault != nil {
				errorMessage := fmt.Sprintf("JSON Schema output validation aborted at step %s: %v\n", step.Path, validateFault)
				log.Printf(SD_ERR + errorMessage)
				actionError += errorMessage
				actionDescriptor.Status = "aborted"
				exitCode = 1
				break
			} else if len(errorList) > 0 {
				actionError += fmt.Sprintf("Validation errors: %v\n", errorList)
				exitCode = 1
				actionDescriptor.Status = "aborted"
				log.Printf(SD_ERR+"Action %s %s at step %s: %v", task.Action, actionDescriptor.Status, step.Path, errorList)
				break
			}
			continue

		} else if actionDescriptor.Status == "pending" {
			actionDescriptor.Status = "running"
			publishStatus(rdb, progressChannel, actionDescriptor) // publish running status
		}

		// Create a pipe to read control commands from action steps
		comReadFd, comWriteFd, _ := os.Pipe()
		comReadLock := make(chan int, 3)

		cmd := exec.Command(step.Path)
		cmd.Env = append(append(os.Environ(), environment...),
			"AGENT_COMFD=3", // 3 is the additional FD number where the action step can write its commands for us
			"AGENT_TASK_ID="+task.ID,
			"AGENT_TASK_ACTION="+task.Action,
		)
		inputData, _ := json.Marshal(task.Data)
		cmd.Stdin = strings.NewReader(string(inputData))
		cmd.ExtraFiles = []*os.File{comWriteFd}
		// Run cmd in a new progress group (PG) to easily send termination signals
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

		stdoutReader, _ := cmd.StdoutPipe()
		// Copy the command stderr to our stderr stream with a pipe tee
		stderrReader, _ := cmd.StderrPipe()
		stderrReader = io.NopCloser(io.TeeReader(stderrReader, os.Stderr))

		go func() {
			bytes, _ := io.ReadAll(stderrReader)
			actionError += string(bytes)
			comReadLock <- 1
		}()

		go func() {
			bytes, _ := io.ReadAll(stdoutReader)
			actionOutput += string(bytes)
			comReadLock <- 1
		}()

		go func() {
			// Read the commands FD line by line, parse and execute the command
			rdr := csv.NewReader(comReadFd)
			rdr.Comma = ' '
			rdr.FieldsPerRecord = -1
			for {
				record, err := rdr.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf(SD_ERR+"Parse error: %v", err)
					continue
				}
				switch cmd := record[0]; cmd {
				case "set-env":
					environment = append(environment, record[1]+"="+record[2])
					isStateWriteNeeded = true
					log.Printf(SD_WARNING + "Command set-env is deprecated")
				case "unset-env":
					for i, envVar := range environment {
						if strings.HasPrefix(envVar, record[1]+"=") {
							// var exists at index i: slice and splice at that index:
							environment = append(environment[:i], environment[i+1:]...)
							isStateWriteNeeded = true
							// the array could have duplicates, continue the iteration
						}
					}
					log.Printf(SD_WARNING + "Command unset-env is deprecated")
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
						continue
					}
					err = actionDescriptor.SetStepWeight(record[1], weight)
					if err != nil {
						log.Printf(SD_ERR+"set-weight command failed: %v", err)
						continue
					}
				case "set-progress":
					var progress int
					var err error
					progress, err = strconv.Atoi(record[1])
					if err != nil {
						log.Printf(SD_ERR+"set-progress command failed: %v", err)
						continue
					}
					err = actionDescriptor.SetProgressAtStep(stepIndex, progress)
					if err != nil {
						log.Printf(SD_ERR+"set-progress command failed: %v", err)
						continue
					}
					publishStatus(rdb, progressChannel, actionDescriptor)
				default:
					log.Printf(SD_ERR+"Unknown command %s", cmd)
				}
			}
			comReadLock <- 1
		}()

		log.Printf("%s/task/%s: %s/%s is starting", agentPrefix, task.ID, task.Action, step.Name)
		if err := cmd.Start(); err != nil {
			exitCode = 9
			actionError += fmt.Sprintf("Action %s startup error at step %s: %v", task.Action, step, err)
			actionDescriptor.Status = "aborted"
			log.Print(SD_ERR + actionError)
			break
		}

		// standard FDs are closed by Start. Our extra FD for commands must be closed manually
		// otherwise it blocks our thread.
		comWriteFd.Close()

		// Block until the three coroutines (stdout, stderr, comfd) finish,
		// or the action is canceled
		doneChan := actionCtx.Done()
		for chanCount := 0; chanCount < 3; {
			select {
			case <-comReadLock:
				chanCount++
			case <-doneChan:
				// Just send a TERM signal to the running step. It then
				// returns an exit code and the whole action is aborted.
				log.Printf(SD_WARNING+"%s/task/%s: Sending TERM signal to action \"%s\" at step %s", agentPrefix, task.ID, task.Action, lastStep)
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
					log.Print(SD_ERR+"Kill failed: ", err)
				}
				doneChan = make(chan struct{}) // stop doneChan from being select()'ed
			case <-time.After(pollingDuration):
				publishStatus(rdb, progressChannel, actionDescriptor)
			}
		}
		if err := cmd.Wait(); err != nil {
			exitCode = cmd.ProcessState.ExitCode()
			if actionDescriptor.Status == "running" {
				actionDescriptor.Status = "aborted"
			}
			break
		}

		if actionDescriptor.Status == "validation-failed" {
			exitCode = 10 // if no exit code was returned, validation-failed forces exit code 10
			log.Printf(SD_WARNING+"Action \"%s\" validation-failed at step %s. Exit code is not set! Forced to 10", task.Action, step.Path)
			break
		}

		// Write the environment state file with changes from the successful step:
		if exitCode == 0 && isStateWriteNeeded {
			writeStateFile(dedupEnv(environment))
		}

		actionDescriptor.SetProgressAtStep(stepIndex, 100)
		publishStatus(rdb, progressChannel, actionDescriptor)
	}

	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// Use a single command pipeline to preserve data integrity.
		if exitCode == 0 {
			if actionDescriptor.Status == "running" {
				actionDescriptor.Status = "completed"
			}
			// Action completed successfully. Publish the environment in a Redis HASH key.
			// Erase the Redis key to honor unset-env:
			pipe.Del(ctx, agentPrefix+"/environment")
			// Re-add the environment as a Redis HASH entry:
			if len(environment) > 0 {
				pipe.HSet(ctx, agentPrefix+"/environment", exportToRedis(environment)...)
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
	}
	log.Printf("%s/task/%s: action \"%s\" status is \"%s\" (%d) at step %s", agentPrefix, task.ID, task.Action, actionDescriptor.Status, exitCode, lastStep)
}

func listenActionsAsync(brpopCtx context.Context, complete chan int) {
	defer func() { complete <- 1 }()
	var workersRegistry sync.WaitGroup
	taskCancelFunctions := make(map[string]context.CancelFunc)

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
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "127.0.0.1:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:      redisAddress,
		Username:  redisUsername,
		Password:  redisPassword,
		DB:        0,
		OnConnect: setClientNameCallback,
	})

	// Ignore the credential error on agent startup
	//
	// Initialize to a well-known error condition, to avoid log pollution
	// and false alarms. When an agent is created its credentials might be
	// still not stored in the leader node Redis instance:
	var lastPopErr error = errors.New("WRONGPASS invalid username-password pair or user is disabled.")

	for { // Action listen loop
		var task models.Task

		// Pop the task from the agent tasks queue
		popResult, popErr := rdb.BRPop(brpopCtx, pollingDuration, agentPrefix+"/tasks").Result()
		if popErr == redis.Nil {
			// poll timeout, it's ok: start a new cycle
			lastPopErr = nil // Clear the last error buffer
			continue
		} else if brpopCtx.Err() != nil {
			break
		} else if popErr != nil {
			// Avoid error log repetitions: print the error only if it is
			// different from the value in the last error buffer
			if lastPopErr == nil || popErr.Error() != lastPopErr.Error() {
				log.Print(SD_ERR+"Task queue pop error: ", popErr)
				lastPopErr = popErr // set the last error buffer
			}
			time.Sleep(pollingDuration)
			continue
		}
		// Task popped from the queue

		lastPopErr = nil // Clear the last error buffer

		if err := json.Unmarshal([]byte(popResult[1]), &task); err != nil {
			log.Print(SD_ERR+"Task ignored for decoding error: ", err)
			continue
		}

		// Store the task as context
		setErr := rdb.Set(ctx, agentPrefix+"/task/"+task.ID+"/context", popResult[1], taskExpireDuration).Err()
		if setErr != nil {
			log.Print(SD_ERR+"Context set error: ", setErr)
		}

		taskCtx, taskCancelFunction := context.WithCancel(ctx)
		taskCancelFunctions[task.ID] = taskCancelFunction

		// run the Action required by the Task payload
		go func(task models.Task) {
			defer func() {
				workersRegistry.Done()
				delete(taskCancelFunctions, task.ID)
			}()
			workersRegistry.Add(1)
			switch task.Action {
			case "list-actions":
				runListActions(rdb, &task)
			case "cancel-task":
				runCancelTask(rdb, &task, taskCancelFunctions)
			default:
				runAction(rdb, taskCtx, &task)
			}
		}(task)
	}
	workersRegistry.Wait() // block until Action coroutines finish
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

func publishStatus(client redis.Cmdable, progressChannel string, actionDescriptor models.Processor) {
	err := client.Publish(ctx, progressChannel, actionDescriptor.ToJSON()).Err()
	if err != nil {
		log.Printf(SD_ERR+"Failed to publish the action status on channel %s", progressChannel)
	}
}
