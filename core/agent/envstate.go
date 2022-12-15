/*
 * Copyright (C) 2022 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package main

import (
	"os"
	"strings"
	"sync"
	"log"
)

var muxEnvironment sync.Mutex

// Read the environment state file from the current working directory.
func readStateFile() []string {
	muxEnvironment.Lock()
	defer func() { muxEnvironment.Unlock() }()
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
func writeStateFile(env []string) {
	muxEnvironment.Lock()
	defer func() { muxEnvironment.Unlock() }()
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