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
	"encoding/json"
	"log"
	"sync"

	"github.com/NethServer/ns8-core/core/agent/models"
	"github.com/go-redis/redis/v8"
)

func runListActions(rdb *redis.Client, task *models.Task) {
	// Redis key names where the action response is stored:
	progressChannel := "progress/" + agentPrefix + "/task/" + task.ID
	outputKey := "task/" + agentPrefix + "/" + task.ID + "/output"
	errorKey := "task/" + agentPrefix + "/" + task.ID + "/error"
	exitCodeKey := "task/" + agentPrefix + "/" + task.ID + "/exit_code"

	actionDescriptor := models.CreateOneStepProcessor(task.Action)
	publishStatus(rdb, progressChannel, actionDescriptor)

	actionDescriptor.Status = "running"
	publishStatus(rdb, progressChannel, actionDescriptor)

	// read the action list from disk
	actions := models.ListProcessors(actionPaths)

	// append agent builtin actions
	actions = append(actions, "list-actions", "cancel-task")

	actionOutput, _ := json.Marshal(actions)
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

func runCancelTask(rdb *redis.Client, task *models.Task, cancelFuncMap map[string]context.CancelFunc, tcMu *sync.Mutex) {
	// Redis key names where the action response is stored:
	progressChannel := "progress/" + agentPrefix + "/task/" + task.ID
	outputKey := "task/" + agentPrefix + "/" + task.ID + "/output"
	errorKey := "task/" + agentPrefix + "/" + task.ID + "/error"
	exitCodeKey := "task/" + agentPrefix + "/" + task.ID + "/exit_code"

	actionError := ""
	actionOutput := ""

	actionDescriptor := models.CreateOneStepProcessor(task.Action)
	publishStatus(rdb, progressChannel, actionDescriptor)

	type CancelRequest struct {
		Task    string
		Timeout int
	}

	var request CancelRequest
	var exitCode int

	payload, _ := json.Marshal(task.Data)
	jsonErr := json.Unmarshal(payload, &request)

	// STEP 1. builtin-validation.json
	lastStep := "builtin-validation.json"
	if jsonErr == nil && len(request.Task) > 0 {
		actionDescriptor.Status = "running"
		publishStatus(rdb, progressChannel, actionDescriptor)

		// STEP 2. task-lookup
		lastStep = "task-lookup"
		tcMu.Lock()
		cancelFunc, hasTask := cancelFuncMap[request.Task]
		tcMu.Unlock()
		if hasTask == true {
			log.Printf("task/%s/%s: %s/%s is starting", agentPrefix, task.ID, task.Action, lastStep)
			cancelFunc() // STEP 3. task cancellation
			lastStep = "cancellation"
			exitCode = 0
			actionDescriptor.SetProgressAtStep(0, 100)
			actionDescriptor.Status = "completed"
		} else {
			exitCode = 2
			actionDescriptor.Status = "validation-failed"
			actionError = "Action cancel-task validation-failed at step task-lookup: task ID not found"
			log.Printf(SD_ERR + actionError)
		}
	} else {
		exitCode = 10
		actionDescriptor.Status = "validation-failed"
		actionError = "Action cancel-task validation-failed at step builtin-validation.json: invalid task ID or timeout value"
		log.Printf(SD_ERR + actionError)
	}

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
	log.Printf("task/%s/%s: action \"%s\" status is \"%s\" (%d) at step %s", agentPrefix, task.ID, task.Action, actionDescriptor.Status, exitCode, lastStep)
}
