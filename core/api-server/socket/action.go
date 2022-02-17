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
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/NethServer/ns8-scratchpad/core/api-server/utils"
)

func Action(socketAction models.SocketAction) {
	// switch action received
	switch socketAction.Action {
	case "logs-stop":
		// decode data payload into specific action
		var logsAction models.LogsStopAction

		// marshal data payload into json string
		jsonStr, err := json.Marshal(socketAction.Data)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[SOCKET] error in json marshal logs-stop action"))
		}

		// convert json string into struct
		if errLogsAction := json.Unmarshal(jsonStr, &logsAction); errLogsAction != nil {
			utils.LogError(errors.Wrap(errLogsAction, "[SOCKET] error in Logs action stop json unmarshal"))
		}

		cmd := Commands[logsAction.Pid]
		if errKill := cmd.Process.Kill(); errKill != nil {
			utils.LogError(errors.Wrap(errKill, "[SOCKET] error in command kill"))
		}

	case "logs-start":
		// decode data payload into specific action
		var logsAction models.LogsStartAction

		// marshal data payload into json string
		jsonStr, err := json.Marshal(socketAction.Data)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[SOCKET] error in json marshal logs-start action"))
		}

		// convert json string into struct
		if errLogsAction := json.Unmarshal(jsonStr, &logsAction); errLogsAction != nil {
			utils.LogError(errors.Wrap(errLogsAction, "[SOCKET] error in Logs action start json unmarshal"))
		}

		// filter logs params
		var mode = ""
		//var filter = ""
		var entity = ""

		// switch mode
		switch logsAction.Mode {
		case "tail":
			mode = "-t"
		case "dump":
			mode = "--limit=" + logsAction.Lines
		}

		// switch entity
		switch logsAction.Entity {
		case "cluster":
			entity = "{job=\"systemd-journal\"} | json | line_format \"{{.nodename}} --> {{.MESSAGE}}\""

		case "node":
			entity = "{nodename=\"" + logsAction.EntityName + "\"} | json | line_format \"{{.nodename}} --> {{.MESSAGE}}\""

		case "module":
			entity = "{job=\"systemd-journal\"} | json | CONTAINER_TAG=\"" + logsAction.EntityName + "\" | line_format \"{{.nodename}} --> {{.MESSAGE}}\""

		}

		// define command
		fmt.Println("/usr/local/bin/logcli", "query", "-q", "--no-labels", mode, entity)
		cmd := exec.Command("/usr/local/bin/logcli", "query", "-q", "--no-labels", mode, entity)

		if logsAction.Mode == "tail" {
			// execute command follow mode
			go func() {
				pid := ""

				// create a pipe for the output of the script
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					utils.LogError(errors.Wrap(err, "[SOCKET] error creating stdout pipe for Cmd"))
					return
				}

				// create scanner to listen to command outputs
				scanner := bufio.NewScanner(cmdReader)
				go func() {
					// foreach command outputs send to websocket
					for scanner.Scan() {
						broadcastToAll(gin.H{"pid": pid, "data": scanner.Text()})
					}
				}()

				// start command
				err = cmd.Start()
				if err != nil {
					return
				}

				// add command to command lists
				pid = strconv.Itoa(cmd.Process.Pid)
				Commands[pid] = cmd

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

				broadcastToAll(gin.H{"pid": "", "data": string(out)})
			}()
		}

	}
}

func broadcastToAll(msg interface{}) {
	// convert interface to json string
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}

	if clientSession, ok := Connections["/ws"]; ok {
		// Broadcast to all sessions
		socketConnection.BroadcastMultiple([]byte(jsonStr), clientSession)
	}
}
