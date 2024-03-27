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
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package main

import (
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/socket"
	"github.com/NethServer/ns8-core/core/api-server/utils"
	"github.com/pkg/errors"

	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// define command flags
var (
	modeFlag       = "tail"
	linesFlag      = "25"
	fromFlag       = ""
	toFlag         = ""
	entityFlag     = "cluster"
	entityNameFlag = ""
	searchFlag     = ""
	timezone       = ""
	instance       = ""
)

var RootCmd = &cobra.Command{
	Use:   "api-server-logs",
	Short: "api-server-logs is a wrapper of api-server websocket logs",
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show api-server-logs version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " " + "0.0.1")
	},
}

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "get logs of a specific entity. default: cluster",
	Run: func(cmd *cobra.Command, args []string) {
		Logs()
	},
}

func Execute() {
	// add commands
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(LogsCmd)

	// define logs command flags
	LogsCmd.Flags().StringVarP(&modeFlag, "mode", "m", "tail", "get logs in a specific mode: tail or dump")
	LogsCmd.Flags().StringVarP(&linesFlag, "lines", "l", "25", "get logs for a specific lines in dump mode")
	LogsCmd.Flags().StringVarP(&fromFlag, "from", "f", "", "get logs from a specific date. ISO8601 format")
	LogsCmd.Flags().StringVarP(&toFlag, "to", "t", "", "get logs to a specific date. ISO8601 format")
	LogsCmd.Flags().StringVarP(&entityFlag, "entity", "e", "cluster", "get logs for a specific entity: cluster, node, module")
	LogsCmd.Flags().StringVarP(&entityNameFlag, "name", "n", "", "get logs for a specific entity name. used in node or module")
	LogsCmd.Flags().StringVarP(&searchFlag, "search", "s", "", "get logs for a specific search string")
	LogsCmd.Flags().StringVarP(&timezone, "timezone", "z", "", "get logs in a specific timezone")
	LogsCmd.Flags().StringVarP(&instance, "instance", "i", "", "search for logs in a specific instance. (Example: loki1, loki2, ...)")

	// check errors on cmd execution
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Logs() {
	// set wait group for commands
	var wg sync.WaitGroup
	wg.Add(1)

	// define payload
	payload := `
  {
   "action": "logs-start",
   "payload": {
      "id": "` + uuid.New().String() + `",
      "mode": "` + modeFlag + `",
      "lines": "` + linesFlag + `",
      "filter": "` + searchFlag + `",
      "from": "` + fromFlag + `",
      "to": "` + toFlag + `",
      "entity": "` + entityFlag + `",
      "entity_name": "` + entityNameFlag + `",
      "timezone": "` + timezone + `",
      "instance": "` + instance + `"
   }
  }
  `

	// init command to execute
	var action models.SocketAction
	if errAction := json.Unmarshal([]byte(payload), &action); errAction != nil {
		utils.LogError(errors.Wrap(errAction, "[LOG-CLI] error in Action json unmarshal"))
	}

	// execute command
	socket.Action(action, nil, &wg)

	// wait for command exit
	wg.Wait()
}

func main() {
	// init cli
	Execute()
}
