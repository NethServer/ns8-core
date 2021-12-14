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
	"errors"
	"fmt"
	"os"
	"sort"
)

const (
	STEP_VALIDATE_INPUT = "validate-input.json"
	STEP_VALIDATE_OUTPUT = "validate-output.json"
)

type Step struct {
	Name   string
	Path   string
	Weight int
}

type Descriptor struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Steps    []Step `json:"-"`
}

func (d Descriptor) ToJSON() []byte {
	actionJson, _ := json.Marshal(d)
	return actionJson
}

func (d *Descriptor) SetProgressAtStep(stepIndex int, stepProgress int) error {
	if stepProgress < 0 || stepProgress > 100 {
		return errors.New(fmt.Sprintf("Invalid step progress value: %d", stepProgress))
	}
	currentProgress := 0
	targetProgress := 0
	for s := 0; s < len(d.Steps); s++ {
		targetProgress += 100 * d.Steps[s].Weight
		if s == (stepIndex - 1) {
			// weighed progress of previous steps
			currentProgress = targetProgress
		}
	}
	currentProgress += stepProgress * d.Steps[stepIndex].Weight // add last step progress
	if targetProgress > 0 {
		d.Progress = currentProgress * 100 / targetProgress
	}
	return nil
}

func (d *Descriptor) getStepIndex(stepName string) (int, error) {
	for stepIndex, step := range d.Steps {
		if step.Name == stepName {
			return stepIndex, nil
		}
	}
	return -1, errors.New("Step " + stepName + " not found")
}

func (d *Descriptor) SetStepWeight(stepName string, weight int) error {
	if weight < 0 {
		return errors.New(fmt.Sprintf("Invalid step weight value: %d", weight))
	}
	stepIndex, err := d.getStepIndex(stepName)
	if err != nil {
		return err
	}
	d.Steps[stepIndex].Weight = weight
	return nil
}

func ListActions(actionPaths []string) []string {
	actionDirs := make(map[string]bool)
	for _, path := range actionPaths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				actionDirs[entry.Name()] = true
			}
		}
	}
	actions := make([]string, 0, len(actionDirs))
	for dir, _ := range actionDirs {
		actions = append(actions, dir)
	}
	return actions
}

func CreateBuiltin(actionName string) Descriptor {
	return Descriptor{Status: "pending", Progress: 0, Steps: []Step{{Name: actionName, Path:"", Weight: 1}}}
}

func Create(actionName string, actionPaths []string) Descriptor {

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

	return Descriptor{Status: "pending", Progress: 0, Steps: steps}
}
