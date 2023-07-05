/*
 * Copyright (C) 2023 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package main

import (
	"log"
	"os"
	"strings"
)

// Read the environment state file from the current working directory.
func readStateFile() []string {
	env := make([]string, 0)
	content, err := os.ReadFile("./environment")
	if err != nil {
		log.Printf(SD_ERR+"Cannot read ./environment file: %s", err)
		return env
	}
	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			continue
		}
		if strings.Contains(line, "=") == false {
			continue
		}
		env = append(env, line)
	}
	return env
}
