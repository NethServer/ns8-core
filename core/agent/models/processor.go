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
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Step struct {
	Name   string
	Path   string
	Weight int
}

type Processor struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Steps    []Step `json:"-"`
}

func CreateOneStepProcessor(stepName string) Processor {
	return Processor{Status: "pending", Progress: 0, Steps: []Step{{Name: stepName, Path:"", Weight: 1}}}
}

/*
 * An abstract "Processor" actually processes a Task or an Event. Given a
 * list of base paths, each subdirectory defines a task action or an event
 * handler. If two base paths contain the same subdir name, the processor
 * counts as one.
 */
func ListProcessors(paths []string) []string {
	procDirs := make(map[string]bool)
	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				procDirs[entry.Name()] = true
			}
		}
	}
	processors := make([]string, 0, len(procDirs))
	for dir, _ := range procDirs {
		processors = append(processors, dir)
	}
	return processors
}

func (d Processor) ToJSON() []byte {
	actionJson, _ := json.Marshal(d)
	return actionJson
}

func (d *Processor) SetProgressAtStep(stepIndex int, stepProgress int) error {
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

func (d *Processor) getStepIndex(stepName string) (int, error) {
	for stepIndex, step := range d.Steps {
		if step.Name == stepName {
			return stepIndex, nil
		}
	}
	return -1, errors.New("Step " + stepName + " not found")
}

func (d *Processor) SetStepWeight(stepName string, weight int) error {
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
