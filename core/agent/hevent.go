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
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/NethServer/ns8-core/core/agent/models"
	"github.com/go-redis/redis/v8"
)

func listenEventsAsync(ctx context.Context, complete chan int) {
	// Connect with default credentials to listen event channels with no
	// restrictions.
	redisAddress := os.Getenv("REDIS_REPLICA_ADDRESS")
	if redisAddress == "" {
		// Fallback to leader address
		redisAddress = os.Getenv("REDIS_ADDRESS")
	}
	if redisAddress == "" {
		// Fallback to local default
		redisAddress = "127.0.0.1:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:      redisAddress,
		Username:  "default",
		Password:  agentPrefix,
		DB:        0,
		OnConnect: setClientNameCallback,
		MaxRetries:       10,
		MinRetryBackoff:   100 * time.Millisecond,
		MaxRetryBackoff:   5000 * time.Millisecond,
	})

	pubsub := rdb.PSubscribe(ctx, "*/event/*")

	var wg sync.WaitGroup
	csyn := make(chan int, 1)

	go func() {
		for msg := range pubsub.Channel(redis.WithChannelHealthCheckInterval(pollingDuration)) {
			if before, after, found := strings.Cut(msg.Channel, "/event/"); found {
				go runEvent(&wg, &models.Event{Source: before, Payload: msg.Payload, Name: after})
			}
		}
		csyn <- 1
	}()

	select {
	case <-csyn:
	case <-ctx.Done():
	}

	pubsub.Close()
	wg.Wait() // wait for running event handlers completion
	complete <- 1
}

func runEvent(wg *sync.WaitGroup, event *models.Event) {
	wg.Add(1)
	defer wg.Done()

	var exitCode int = 0
	var lastStep string = ""

	handler := models.CreateEventHandler(event.Name, eventPaths) // global variable
	if len(handler.Steps) == 0 {
		return // if event handler is not defined, skip silently this run
	}

	// Get additional environment variables from the filesystem
	environment := readStateFile()
	for stepIndex, step := range handler.Steps {
		lastStep = step.Name

		handler.Status = "running"

		cmd := exec.Command(step.Path)
		cmd.Env = append(append(os.Environ(), environment...),
			"AGENT_COMFD=2", // Commands are redirected to stderr
			"AGENT_EVENT_NAME="+event.Name,
			"AGENT_EVENT_SOURCE="+event.Source,
		)
		cmd.Stdin = strings.NewReader(event.Payload)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Printf(SD_WARNING + "Handler of %s/event/%s skipped step %s: %v", event.Source, event.Name, step.Name, err)
			handler.SetProgressAtStep(stepIndex, 100)
			continue
		}
		log.Printf(SD_DEBUG+"Handler of %s/event/%s is starting step %s", event.Source, event.Name, step.Name)

		if err := cmd.Wait(); err != nil {
			log.Print(SD_ERR, err)
			exitCode = cmd.ProcessState.ExitCode()
			if handler.Status == "running" {
				handler.Status = "aborted"
			}
			break
		}

		handler.SetProgressAtStep(stepIndex, 100)
	}

	if exitCode == 0 {
		handler.Status = "completed"
	}

	log.Printf("Handler of %s/event/%s exited with status \"%s\" (%d) at step %s", event.Source, event.Name, handler.Status, exitCode, lastStep)
}
