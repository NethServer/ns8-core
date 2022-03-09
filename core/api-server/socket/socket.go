/*
 * Copyright (C) 2021 Nethesis S.r.l.
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
	"encoding/json"
	"fmt"
	"os/exec"

	jwt "github.com/golang-jwt/jwt"
	"github.com/olahol/melody"

	"github.com/pkg/errors"

	"github.com/NethServer/ns8-scratchpad/core/api-server/configuration"
	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/NethServer/ns8-scratchpad/core/api-server/utils"
)

var socketConnection *melody.Melody

var Connections map[string][]*melody.Session
var Commands map[string]map[string]*exec.Cmd

func Instance() *melody.Melody {
	if socketConnection == nil {
		socketConnection := InitSocketConnection()
		return socketConnection
	}
	return socketConnection
}

func InitSocketConnection() *melody.Melody {
	// init socket
	socketConnection := melody.New()

	// assign handlers
	socketConnection.HandleConnect(OnConnect)
	socketConnection.HandleDisconnect(OnDisconnect)
	socketConnection.HandleMessage(OnMessage)

	// init connection obj
	Connections = make(map[string][]*melody.Session)

	// init commands obj
	Commands = make(map[string]map[string]*exec.Cmd)

	return socketConnection
}

func ValidateAuth(tokenString string) bool {
	// convert token string and validate it
	if tokenString != "" {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// validate the alg
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// return secret
			return []byte(configuration.Config.Secret), nil
		})

		if err != nil {
			utils.LogError(errors.Wrap(err, "[SOCKET] error in JWT token validation"))
			return false
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["id"] != nil {
				return true
			}
		} else {
			utils.LogError(errors.Wrap(err, "[SOCKET] error in JWT token claims"))
			return false
		}
	}
	return false
}

func OnConnect(s *melody.Session) {
	// URL Path is the unique connection ID and append all sessions which are connected to that URL
	Connections[s.Request.URL.Path] = append(Connections[s.Request.URL.Path], s)
}

func OnDisconnect(s *melody.Session) {
	// reassign existing connections
	sessions := Connections[s.Request.URL.Path]
	var newSessions []*melody.Session
	for _, existingSession := range sessions {
		if s != existingSession {
			newSessions = append(newSessions, existingSession)
		}
	}
	Connections[s.Request.URL.Path] = newSessions

	// kill running processes
	for pid := range Commands[s.Request.Header["Sec-Websocket-Key"][0]] {
		cmd := Commands[s.Request.Header["Sec-Websocket-Key"][0]][pid]
		cmd.Process.Kill()
	}
}

func OnMessage(s *melody.Session, msg []byte) {
	// get action received
	var socketAction models.SocketAction
	if errSocketAction := json.Unmarshal([]byte(msg), &socketAction); errSocketAction != nil {
		utils.LogError(errors.Wrap(errSocketAction, "[SOCKET] error in Action json unmarshal"))
	}

	// switch action received
	Action(socketAction, s)
}
