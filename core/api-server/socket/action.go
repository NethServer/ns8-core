/*
 * Copyright (C) 2022 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of NethServer project.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with NethServer. If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package socket

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/utils"
)

func Action(socketAction models.SocketAction, s *melody.Session, wg *sync.WaitGroup) {
	// switch action received
	switch socketAction.Action {
	case "logs-stop":
		// decode data payload into specific action
		var logsAction models.LogsStopAction

		// marshal data payload into json string
		jsonStr, err := json.Marshal(socketAction.Payload)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[SOCKET] error in json marshal logs-stop action"))
		}

		// convert json string into struct
		if errLogsAction := json.Unmarshal(jsonStr, &logsAction); errLogsAction != nil {
			utils.LogError(errors.Wrap(errLogsAction, "[SOCKET] error in Logs action stop json unmarshal"))
		}

		cmd := Commands[s.Request.Header["Sec-Websocket-Key"][0]][logsAction.Pid]

		// check if specific pid exists and command is valid
		if cmd != nil {
			if errKill := cmd.Process.Kill(); errKill != nil {
				utils.LogError(errors.Wrap(errKill, "[SOCKET] error in command kill"))
			}

			broadcastToAll("logs-stop", gin.H{"id": logsAction.Id, "pid": logsAction.Pid, "message": "logs follow stopped"})
		}

	case "logs-start":
		// decode data payload into specific action
		var logsAction models.LogsStartAction

		// marshal data payload into json string
		jsonStr, err := json.Marshal(socketAction.Payload)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[SOCKET] error in json marshal logs-start action"))
		}

		// convert json string into struct
		if errLogsAction := json.Unmarshal(jsonStr, &logsAction); errLogsAction != nil {
			utils.LogError(errors.Wrap(errLogsAction, "[SOCKET] error in Logs action start json unmarshal"))
		}

		// filter logs params
		var mode = ""
		var filter = ""
		var entity = ""
		var from = ""
		var to = ""
		var timezone = "UTC"

		// define basic args and envs
		args := []string{"query", "-q", "--no-labels"}
		envs := []string{}

		// add timezone
		if len(logsAction.TimeZone) > 0 {
			timezone = logsAction.TimeZone
			args = append(args, "--timezone=Local")
			envs = append(envs, "TZ="+timezone)
		}

		// check date
		if len(logsAction.From) > 0 {
			from = "--from=" + logsAction.From
			args = append(args, from)
		}

		if len(logsAction.To) > 0 {
			to = "--to=" + logsAction.To
			args = append(args, to)
		}

		// switch mode
		switch logsAction.Mode {
		case "tail":
			mode = "-t"
			args = append(args, "--forward")
		case "dump":
			mode = "--limit=" + logsAction.Lines
		}
		args = append(args, mode)

		// check filter
		if len(logsAction.Filter) > 0 {
			filter = "json | line_format \"{{.SYSLOG_IDENTIFIER}} {{.nodename}} {{.MESSAGE}}\" |= `" + logsAction.Filter + "`"
		} else {
			filter = "json"
		}

		// switch entity
		switch logsAction.Entity {
		case "cluster":
			entity = "{job=\"systemd-journal\"} | " + filter + " | line_format \"{{.nodename}} {{.MESSAGE}}\""

		case "node":
			entity = "{nodename=\"" + logsAction.EntityName + "\"} | " + filter + " | line_format \"{{.MESSAGE}}\""

		case "module":
			entity = "{job=\"systemd-journal\"} | " + filter + " | SYSLOG_IDENTIFIER=~\"" + logsAction.EntityName + "(/.+)?\" | line_format \"{{.SYSLOG_IDENTIFIER}} {{.MESSAGE}}\""

		}
		args = append(args, entity)

		// define command
		fmt.Println("TZ="+timezone+" /usr/local/bin/logcli", args)
		cmd := exec.Command("/usr/local/bin/logcli", args...)

		// add envs
		cmd.Env = os.Environ()
		for e := 0; e < len(envs); e++ {
			cmd.Env = append(cmd.Env, envs[e])
		}

		if logsAction.Mode == "tail" {
			// execute command follow mode
			go func() {
				pid := ""

				// create a pipe for the output of the script
				stdout, errStdOut := cmd.StdoutPipe()
				if errStdOut != nil {
					return
				}

				// create scanner to listen to command outputs
				scannerStdOut := bufio.NewScanner(stdout)
				go func() {
					// foreach command outputs send to websocket
					for scannerStdOut.Scan() {
						if s != nil {
							broadcastToAll("logs-start", gin.H{"id": logsAction.Id, "pid": pid, "message": scannerStdOut.Text()})
						} else {
							fmt.Println(scannerStdOut.Text())
						}
					}
				}()

				// start command
				err = cmd.Start()
				if err != nil {
					return
				}

				// add command to command lists
				if s != nil {
					pid = strconv.Itoa(cmd.Process.Pid)
					if Commands[s.Request.Header["Sec-Websocket-Key"][0]] == nil {
						Commands[s.Request.Header["Sec-Websocket-Key"][0]] = make(map[string]*exec.Cmd)
					}
					Commands[s.Request.Header["Sec-Websocket-Key"][0]][pid] = cmd

					// send feedback for command received
					broadcastToAll("logs-start", gin.H{"id": logsAction.Id, "pid": pid, "message": ""})
				}

				// use Wait to avoid defunct process when killed
				err = cmd.Wait()
				if err != nil {
					return
				}
			}()
		}

		if logsAction.Mode == "dump" {
			go func() {
				out, err := cmd.Output()
				if err != nil {
					utils.LogError(errors.Wrap(err, "[SOCKET] error executing Cmd for dump"))
				}

				// reverse logs orders
				logsStrings := strings.Split(string(out), "\n")
				logsStringsR := reverse(logsStrings)
				logsStringsOut := strings.Join(logsStringsR[:], "\n")

				if s != nil {
					broadcastToAll("logs-start", gin.H{"id": logsAction.Id, "pid": "", "message": logsStringsOut})
				} else {
					fmt.Println(logsStringsOut)
					defer wg.Done()
				}
			}()
		}

	}
}

func broadcastToAll(name string, msg interface{}) {
	// create event object
	event := &models.Event{}
	event.Name = name
	event.Timestamp = time.Now().UTC()
	event.Payload = msg
	event.Type = "action"

	// convert interface to json string
	actionJSON, err := json.Marshal(event)
	if err != nil {
		utils.LogError(errors.Wrap(err, "[SOCKET] error converting interface msg to broadcast"))
	}

	if clientSession, ok := Connections["/ws"]; ok {
		// Broadcast to all sessions
		socketConnection.BroadcastMultiple(actionJSON, clientSession)
	}
}

func reverse(ss []string) []string {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}

	return ss
}
