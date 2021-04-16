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
	"github.com/olahol/melody"
	//"github.com/NethServer/ns8-scratchpad/core/api-server/utils"
)

var socketConnection *melody.Melody

var Connections map[string][]*melody.Session

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

	return socketConnection
}

func OnConnect(s *melody.Session) {
	// URL Path is the unique connection ID and append all sessions which are connected to that URL
	Connections[s.Request.URL.Path] = append(Connections[s.Request.URL.Path], s)
}

func OnDisconnect(s *melody.Session) {
	sessions := Connections[s.Request.URL.Path]
	var newSessions []*melody.Session
	for _, existingSession := range sessions {
		if s != existingSession {
			newSessions = append(newSessions, existingSession)
		}
	}
	Connections[s.Request.URL.Path] = newSessions
}

func OnMessage(s *melody.Session, msg []byte) {
	socketConnection.Broadcast(msg)
}
