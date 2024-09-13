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
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

// Reference: https://www.man7.org/linux/man-pages/man3/sd-daemon.3.html
const (
	SD_EMERG   = "<0>" /* system is unusable */
	SD_ALERT   = "<1>" /* action must be taken immediately */
	SD_CRIT    = "<2>" /* critical conditions */
	SD_ERR     = "<3>" /* error conditions */
	SD_WARNING = "<4>" /* warning conditions */
	SD_NOTICE  = "<5>" /* normal but significant condition */
	SD_INFO    = "<6>" /* informational */
	SD_DEBUG   = "<7>" /* debug-level messages */
)

var ctx = context.Background()
var agentPrefix string
var actionPaths flagStringSlice
var eventPaths flagStringSlice

var pollingDuration = 5000 * time.Millisecond
var taskExpireDuration = 8 * time.Hour

// Command arguments --actionsdir and --eventsdir can be repeated multiple
// times. Each item is inserted into a []string.
type flagStringSlice []string

func (v *flagStringSlice) String() string {
	return "" // doesn't matter
}

func (v *flagStringSlice) Set(raw string) error {
	if len(raw) > 0 {
		*v = append(*v, raw)
	}
	return nil
}

func setClientNameCallback(ctx context.Context, cn *redis.Conn) error {
	return cn.ClientSetName(ctx, agentPrefix).Err()
}

func main() {

	// Parse command-line arguments
	flag.Var(&actionPaths, "actionsdir", "Directory path with actions definition")
	flag.Var(&eventPaths, "eventsdir", "Directory path with events definition")
	flag.StringVar(&agentPrefix, "agentid", "", "Agent ID and default Redis username")
	flag.Parse()

	log.SetFlags(0)
	if agentPrefix == "" {
		log.Fatal(SD_EMERG + "The agent prefix argument is not set")
		// exit(1) log.Fatal terminates the process with an exit code
	}
	if len(actionPaths) == 0 {
		log.Fatal(SD_EMERG + "Action path command arguments are not set")
		// exit(1) log.Fatal terminates the process with an exit code
	}

	// Override default length of polling interval (5000ms)
	ePollingDuration := os.Getenv("AGENT_POLLING_INTERVAL")
	if ePollingDuration != "" {
		oValue, convError := time.ParseDuration(ePollingDuration)
		if convError == nil {
			pollingDuration = oValue
		}
	}

	var signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGUSR1)
	var actionsCtx, cancelActions = context.WithCancel(ctx)
	var eventsCtx, cancelEvents = context.WithCancel(ctx)

	go func() {
		xsig := <-signalChannel
		signal.Stop(signalChannel)
		log.Printf(SD_WARNING+"Signal \"%v\" caught: shutdown started.", xsig)
		cancelActions()
	}()

	actionsChannel := make(chan int, 1)
	eventsChannel := make(chan int, 1)

	go listenActionsAsync(actionsCtx, actionsChannel)
	go listenEventsAsync(eventsCtx, eventsChannel)

	// wait completion of action goroutines
	<-actionsChannel

	cancelEvents()

	// wait completion of event goroutines
	<-eventsChannel
}
