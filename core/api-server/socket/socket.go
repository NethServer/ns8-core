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
	"syscall"

	"github.com/olahol/melody"
	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/utils"
)

var socketConnection *melody.Melody
var muClock *utils.MuClock

func Instance() *melody.Melody {
	if socketConnection == nil {
		muClock = new(utils.MuClock)
		muClock.Sync()
		socketConnection = melody.New()
		socketConnection.HandleDisconnect(onDisconnect)
		socketConnection.HandleMessage(onMessage)
		socketConnection.HandlePong(onPong)
	}
	return socketConnection
}

/*
 * Check if the session is still valid every time the ping-pong websocket
 * message is received. If a session did not send the authorize message
 * forcibly close it.
 */
func onPong(s *melody.Session) {
	muClock.Sync()
	if ! ValidSessionFilter(s) {
		oErrorMsg := map[string]string{
			"type": "authorize-error",
			"payload": "Token has expired",
		}
		jErrorMsg, _ := json.Marshal(oErrorMsg)
		s.Write(jErrorMsg)
		s.CloseWithMsg(melody.FormatCloseMessage(1000, "Bye"))
	}
}

/*
 * onDisconnect
 * - Terminate any process spawned by the session
 */
func onDisconnect(s *melody.Session) {
	iChilds, keyExists := s.Get("childs")
	if ! keyExists {
		return // nothing to do
	}
	for _, pid := range iChilds.([]int) {
		// negative pid terminates the whole process group (PG)
		syscall.Kill(-pid, syscall.SIGTERM)
	}
	s.UnSet("childs")
}

func onMessage(s *melody.Session, msg []byte) {
	// get action received
	var socketAction models.SocketAction
	if errSocketAction := json.Unmarshal([]byte(msg), &socketAction); errSocketAction != nil {
		utils.LogError(errors.Wrap(errSocketAction, "[SOCKET] error in Action json unmarshal"))
	}

	// switch action received
	Action(socketAction, s, nil)
}

func ValidSessionFilter(s *melody.Session) bool {
	if muClock.Now() > getSessionExpireTimestamp(s) {
		return false // Session is expired
	}
	return true
}

func getSessionExpireTimestamp(s *melody.Session) int64 {
	exp, ok := s.Get("exp")
	if ! ok {
		return 0
	}
	texp, _ := exp.(int64)
	return texp
}
