/*
 * Copyright (C) 2022 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package main

import (
	"log"
	"os"
	"strings"
	"sync"
)

var muxEnvironment sync.Mutex

// Read the environment file from the current working directory.
func readEnvironmentFile() []string {
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

// Persist the environment state file to the current working directory
// Systemd unit file can pickup it with EnvironmentFile option
func writeEnvironmentFile(env []string) {
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
}

func setEnvironmentVariable(varname string, varvalue string) {
	muxEnvironment.Lock()
	defer func() { muxEnvironment.Unlock() }()
	environment := readEnvironmentFile()
	environment = append(environment, varname + "=" + varvalue)
	writeEnvironmentFile(environment)
}

func unsetEnvironmentVariable(varname string) {
	muxEnvironment.Lock()
	defer func() { muxEnvironment.Unlock() }()
	environment := readEnvironmentFile()
	for i, envVar := range environment {
		if strings.HasPrefix(envVar, varname + "=") {
			// var exists at index i: slice and splice at that index:
			environment = append(environment[:i], environment[i+1:]...)
			// the array could have duplicates, continue the iteration
		}
	}
	writeEnvironmentFile(environment)
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
