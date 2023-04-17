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

package models

import (
	"time"
	"os"
	"sort"
)

type Task struct {
	ID        string      `json:"id" structs:"id"`
	Action    string      `json:"action" structs:"action"`
	Data      interface{} `json:"data" structs:"data"`
	Extra     interface{} `json:"extra" structs:"extra"`
	Queue     string      `json:"queue" structs:"queue"`
	User      string      `json:"user" structs:"user"`
	Timestamp time.Time   `json:"timestamp" structs:"timestamp"`
	Parent    string      `json:"parent" structs:"parent"`
}

const (
	STEP_VALIDATE_INPUT = "validate-input.json"
	STEP_VALIDATE_OUTPUT = "validate-output.json"
)

func CreateTaskProcessor(actionName string, actionPaths []string) Processor {

	actionSteps := make(map[string]string)

	var inputSchema = ""
	var outputSchema = ""

	// Squash the action dirs in a single list
	for _, path := range actionPaths {
		actionPath := path + "/" + actionName
		entries, err := os.ReadDir(actionPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			} else if entry.Name() == STEP_VALIDATE_INPUT {
				inputSchema = actionPath + "/" + STEP_VALIDATE_INPUT
			} else if entry.Name() == STEP_VALIDATE_OUTPUT {
				outputSchema = actionPath + "/" + STEP_VALIDATE_OUTPUT
			} else {
				actionSteps[entry.Name()] = actionPath + "/" + entry.Name()
			}
		}
	}

	actionStepsKeys := make([]string, len(actionSteps))
	i := 0
	for k, _ := range actionSteps {
		actionStepsKeys[i] = k
		i += 1
	}

	// Executable steps under the actionPaths are launched sequentially
	// in alphabetical order:
	sort.Strings(actionStepsKeys)

	steps := make([]Step, len(actionStepsKeys))
	for i, file := range actionStepsKeys {
		steps[i] = Step{Name: file, Path: actionSteps[file], Weight: 1}
	}

	if inputSchema != "" {
		steps = append([]Step{{Name: STEP_VALIDATE_INPUT, Path: inputSchema, Weight: 0}}, steps...)
	}

	if outputSchema != "" {
		steps = append(steps, Step{Name: STEP_VALIDATE_OUTPUT, Path: outputSchema, Weight: 0})
	}

	return Processor{Status: "pending", Progress: 0, Steps: steps}
}
