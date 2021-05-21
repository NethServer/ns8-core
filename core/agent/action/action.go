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

package action

import (
	"encoding/json"
	"os"
	"sort"
)

type Step struct {
	Name   string
	Path   string
	Weight int
}

type Descriptor struct {
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	CurrentStep int    `json:"-"`
	Steps       []Step `json:"-"`
}

func (d Descriptor) ToJSON() []byte {
	actionJson, _ := json.Marshal(d)
	return actionJson
}

func (d *Descriptor) SetProgressAtStep(stepIndex int, stepProgress int) {
	currentProgress := 0
	targetProgress := 0
	for s := 0; s < len(d.Steps); s++ {
		targetProgress += 100 * d.Steps[s].Weight
		if s == (stepIndex - 1) {
			// weighted progress of previous steps
			currentProgress = targetProgress
		}
	}
	currentProgress += stepProgress * d.Steps[stepIndex].Weight // add last step progress
	if targetProgress > 0 {
		d.Progress = currentProgress * 100 / targetProgress
	}
	return
}

func Create(actionName string, actionPaths []string) Descriptor {

	actionSteps := make(map[string]string)

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
			}
			actionSteps[entry.Name()] = actionPath + "/" + entry.Name()
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

	return Descriptor{Status: "pending", Progress: 0, CurrentStep: 0, Steps: steps}
}
