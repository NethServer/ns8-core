/*
 * Copyright (C) 2023 Nethesis S.r.l.
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
	outputKey := "task/" + agentPrefix + "/" + task.ID + "/output"
	errorKey := "task/" + agentPrefix + "/" + task.ID + "/error"
	exitCodeKey := "task/" + agentPrefix + "/" + task.ID + "/exit_code"

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

	for stepIndex, step := range actionDescriptor.Steps {
		lastStep = step.Name

		// Sync environment variables from the filesystem. Reading the
		// environment file on every step allows to load changes produced
		// by child actions.
		environment := readStateFile()

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
			"AGENT_TASK_USER="+task.User,
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

		if err := cmd.Start(); err != nil {
			log.Printf(SD_WARNING + "Action %s skipped step %s: %v", task.Action, step.Name, err)
			actionDescriptor.SetProgressAtStep(stepIndex, 100)
			publishStatus(rdb, progressChannel, actionDescriptor)
			continue
		}
		log.Printf("task/%s/%s: %s/%s is starting", agentPrefix, task.ID, task.Action, step.Name)

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
				log.Printf(SD_WARNING+"task/%s/%s: Sending TERM signal to action \"%s\" at step %s", agentPrefix, task.ID, task.Action, lastStep)
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

		actionDescriptor.SetProgressAtStep(stepIndex, 100)
		publishStatus(rdb, progressChannel, actionDescriptor)
	}

	if exitCode == 0 && actionDescriptor.Status == "running" {
		actionDescriptor.Status = "completed"
	}

	_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// Use a single command pipeline to preserve data integrity.

		// Publish the final environment copy in a Redis HASH key
		pipe.Del(ctx, agentPrefix + "/environment")
		pipe.HSet(ctx, agentPrefix + "/environment", exportToRedis(readStateFile())...)

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
	log.Printf("task/%s/%s: action \"%s\" status is \"%s\" (%d) at step %s", agentPrefix, task.ID, task.Action, actionDescriptor.Status, exitCode, lastStep)
}

func listenActionsAsync(actionsCtx context.Context, complete chan int) {
	defer func() { complete <- 1 }()
	var workersRegistry sync.WaitGroup
	brpopCtx, cancelBrpop := context.WithCancel(ctx)
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
		MaxRetries:       10,
		MinRetryBackoff:   100 * time.Millisecond,
		MaxRetryBackoff:   5000 * time.Millisecond,
	})

	var tcMu sync.Mutex
	taskCh := make(chan models.Task)
	go readTasks(rdb, brpopCtx, taskCh)
	MAINLOOP:
		for {
			select {
			case <-actionsCtx.Done():
				go func() {
					// SIGUSR1 received: wait until all workers complete
					workersRegistry.Wait()
					// Then stop reading events
					cancelBrpop()
				}()
				// Prevents entering this branch further
				actionsCtx = ctx
			case task, stillListening := <-taskCh:
				if stillListening == false {
					workersRegistry.Wait()
					break MAINLOOP
				}
				// Create a cancelable context for the task and
				// store its cancel function in a safe map
				taskCtx, taskCancelFunction := context.WithCancel(ctx)
				tcMu.Lock()
				taskCancelFunctions[task.ID] = taskCancelFunction
				tcMu.Unlock()

				// Start the task worker
				go func(task models.Task) {
					workersRegistry.Add(1)
					defer workersRegistry.Done()
					switch task.Action {
					case "list-actions":
						runListActions(rdb, &task)
					case "cancel-task":
						runCancelTask(rdb, &task, taskCancelFunctions, &tcMu)
					default:
						runAction(rdb, taskCtx, &task)
					}
					tcMu.Lock()
					delete(taskCancelFunctions, task.ID)
					tcMu.Unlock()
				}(task)
			}
		}
}

func readTasks(rdb *redis.Client, brpopCtx context.Context, taskCh chan models.Task) {
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
		setErr := rdb.Set(ctx, "task/" + agentPrefix + "/" + task.ID + "/context", obscureTaskInput(popResult[1]), taskExpireDuration).Err()
		if setErr != nil {
			log.Print(SD_ERR+"Context set error: ", setErr)
		}

		taskCh <- task
	}
	close(taskCh)
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

func obscureTaskInput(jsonStr string) string {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		log.Println(SD_ERR+"Error unmarshalling JSON:", err)
		return jsonStr
	}

	recursiveObscureSensitiveKeys(jsonData)

	updatedJson, err := json.Marshal(jsonData)
	if err != nil {
		log.Println(SD_ERR+"Error marshalling JSON:", err)
		return jsonStr
	}

	return string(updatedJson)
}

func isSensitive(target string) bool {
	sensitiveList := []string{"password", "secret", "token"}
	ltarget := strings.ToLower(target)
	for _, sensitive := range sensitiveList {
		if strings.HasSuffix(ltarget, sensitive) {
			return true
		}
	}
	return false
}

func recursiveObscureSensitiveKeys(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		// It's an object, so iterate through its key-value pairs
		for key, value := range v {
			// Recursively update the value
			recursiveObscureSensitiveKeys(value)

			// Check for sensitive keys
			if isSensitive(key) {
				v[key] = "XXX" // replace the secret value
			}
		}

	case []interface{}:
		// It's an array, so iterate through its elements
		for _, item := range v {
			// Recursively update the element
			recursiveObscureSensitiveKeys(item)
		}
	}
}
