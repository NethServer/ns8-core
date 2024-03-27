/*
 * Copyright (C) 2023 Nethesis S.r.l.
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
	"syscall"

	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/olahol/melody"
	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/middleware"
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

		// Lookup the action Pid among the child Pids spawned by the session
		tPid, err := strconv.Atoi(logsAction.Pid)
		if err != nil {
			break
		}
		iChilds, keyExists := s.Get("childs")
		if keyExists {
			aChilds := iChilds.([]int)
			for idx, xPid := range aChilds {
				// If a match is found, the Pid belongs to the session and
				// it can be terminated
				if tPid == xPid {
					// negative pid terminates the whole process group (PG)
					syscall.Kill(-tPid, syscall.SIGTERM)
					// update the session storage, removing the idx-element from the array
					s.Set("childs", append(aChilds[:idx], aChilds[idx+1:]...))
					writeSocketResponse(s, "logs-stop", gin.H{"id": logsAction.Id, "pid": logsAction.Pid, "message": "logs follow stopped"})
					break
				}
			}
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
		var streamSelector = ""
		var logqlPipeline = ""
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

		if len(logsAction.Instance) > 0 {
			envs = append(envs, "LOKI_INSTANCE="+logsAction.Instance)
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
			filter = ` |= "` + strings.ReplaceAll(logsAction.Filter, `"`, `\"`) + `"`
		} else {
			filter = ``
		}

		// switch entity
		switch logsAction.Entity {
		default:
			streamSelector = `{node_id=~".+"}`

		case "node":
			streamSelector = `{node_id="`+ logsAction.EntityName +`"}`

		case "module":
			streamSelector = `{module_id="`+ logsAction.EntityName +`"}`
		}

		logqlPipeline = ` | json syslog_id="SYSLOG_IDENTIFIER", message="MESSAGE" | line_format "[{{.node_id}}:{{.module_id}}:{{.syslog_id}}] {{.message}}"`

		// Compose and append the query strings to logcli arguments
		args = append(args, streamSelector + logqlPipeline + filter)

		// define command
		cmd := exec.Command("/usr/local/bin/logcli", args...)

		// add envs
		cmd.Env = append(os.Environ(), envs...)

		// Run cmd in a new process group (PG) to easily send termination signals
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

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
							writeSocketResponse(s, "logs-start", gin.H{"id": logsAction.Id, "pid": pid, "message": scannerStdOut.Text()})
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

				if s != nil {
					// In a Melody session, store the command pid so it
					// can be killed if connection is closed
					var aChilds []int
					iChilds, keyExists := s.Get("childs")
					if keyExists {
						aChilds = append(iChilds.([]int), cmd.Process.Pid)
					} else {
						aChilds = []int{cmd.Process.Pid}
					}
					s.Set("childs", aChilds)
					// send feedback for command received
					writeSocketResponse(s, "logs-start", gin.H{"id": logsAction.Id, "pid": strconv.Itoa(cmd.Process.Pid), "message": ""})
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
					writeSocketResponse(s, "logs-start", gin.H{"id": logsAction.Id, "pid": "", "message": logsStringsOut})
				} else {
					fmt.Println(logsStringsOut)
					defer wg.Done()
				}
			}()
		}

	case "authorize":
		authPayload, ok := socketAction.Payload.(map[string]interface{})
		if ! ok {
			utils.LogError(errors.New("Authorize payload is corrupt"))
		}

		// Check "jwt" attribute is a string
		token, ok := authPayload["jwt"].(string)
		if ! ok {
			utils.LogError(errors.New("Unknown authorize payload"))
		}

		// Parse "jwt" string and check it is a valid JWT token
		oJwt, err := middleware.InstanceJWT().ParseTokenString(token)
		if err != nil {
			utils.LogError(errors.Wrap(err, "Websocket auth error"))
			oErrorMsg := map[string]string{
				"type": "authorize-error",
				"payload": err.Error(),
			}
			jErrorMsg, _ := json.Marshal(oErrorMsg)
			s.Write(jErrorMsg)
			s.CloseWithMsg(melody.FormatCloseMessage(1000, "Bye"))
			break
		}

		// Authentication is successful: store JWT claims to filter Melody
		// sessions:
		mClaims := jwt.ExtractClaimsFromToken(oJwt)

		exp := int64(mClaims["exp"].(float64))

		// Store the session expire timestamp. In case of conversion
		// error, exp is zero.
		s.Set("exp", exp)

		// Do not send back any message to a successfully authorized
		// session, just keep the socket open.
	}
}

func reverse(ss []string) []string {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}

	return ss
}
