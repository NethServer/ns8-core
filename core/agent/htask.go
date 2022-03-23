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

func runAction(actionCtx context.Context, task *models.Task) {

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

	// Get additional environment variables from the filesystem
	environment := readEnvironmentFile()

	for stepIndex, step := range actionDescriptor.Steps {
		lastStep = step.Name

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
				case "unset-env":
					for i, envVar := range environment {
						if strings.HasPrefix(envVar, record[1]+"=") {
							environment = append(environment[:i], environment[i+1:]...)
						}
					}
				case "dump-env":
					dumpToFile(environment)
					rdb.HSet(ctx, agentPrefix+"/environment", exportToRedis(environment)...).Result()
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

		actionDescriptor.SetProgressAtStep(stepIndex, 100)
		publishStatus(rdb, progressChannel, actionDescriptor)

		environment = dedupEnv(environment)
	}

	if exitCode == 0 {
		dumpToFile(environment)
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
	}
	log.Printf("%s/task/%s: action \"%s\" status is \"%s\" (%d) at step %s", agentPrefix, task.ID, task.Action, actionDescriptor.Status, exitCode, lastStep)
}

func listenActionsAsync(brpopCtx context.Context, rdb *redis.Client, complete chan int) {
	defer func() { complete <- 1 }()
	var workersRegistry sync.WaitGroup
	taskCancelFunctions := make(map[string]context.CancelFunc)

	for { // Action listen loop
		var task models.Task

		// Pop the task from the agent tasks queue
		popResult, popErr := rdb.BRPop(brpopCtx, pollingDuration, agentPrefix+"/tasks").Result()
		if popErr == redis.Nil {
			continue
		} else if brpopCtx.Err() != nil {
			break
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
				runListActions(&task)
			case "cancel-task":
				runCancelTask(&task, taskCancelFunctions)
			default:
				runAction(taskCtx, &task)
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

func publishStatus(client redis.Cmdable, progressChannel string, actionDescriptor models.Processor) {
	err := client.Publish(ctx, progressChannel, actionDescriptor.ToJSON()).Err()
	if err != nil {
		log.Printf(SD_ERR+"Failed to publish the action status on channel %s", progressChannel)
	}
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
