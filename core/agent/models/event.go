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
	"os"
	"sort"
)

type Event struct {
	Source    string
	Name      string
	Payload   string
}

func CreateEventHandler(eventName string, basePaths []string) Processor {

	handlerSteps := make(map[string]string)

	// Squash the action dirs in a single list
	for _, path := range basePaths {
		handlerPath := path + "/" + eventName
		entries, err := os.ReadDir(handlerPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			} else {
				handlerSteps[entry.Name()] = handlerPath + "/" + entry.Name()
			}
		}
	}

	handlerStepsKeys := make([]string, len(handlerSteps))
	i := 0
	for k, _ := range handlerSteps {
		handlerStepsKeys[i] = k
		i += 1
	}

	// Executable steps under the basePaths are launched sequentially
	// in alphabetical order:
	sort.Strings(handlerStepsKeys)

	steps := make([]Step, len(handlerStepsKeys))
	for i, file := range handlerStepsKeys {
		steps[i] = Step{Name: file, Path: handlerSteps[file], Weight: 1}
	}

	return Processor{Status: "pending", Progress: 0, Steps: steps}
}

