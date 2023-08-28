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
	"context"
	"encoding/json"
	"syscall"
	"time"

	"github.com/olahol/melody"
	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/redis"
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
		relayRedisProgressEvents(socketConnection)
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
	if ! validSessionFilter(s) {
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

func validSessionFilter(s *melody.Session) bool {
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

/*
 * Subscribe Redis progress/* channels and relay events from those
 * channels to authenticated websocket sessions
 */
func relayRedisProgressEvents(socketConnection *melody.Melody) {
	ctx := context.Background()

	// init redis connection
	redisConnection := redis.Instance()

	// subscribe to progress channels and listen to new messages:
	progress := redisConnection.PSubscribe(ctx, "progress/*")

	// get the channel to use
	channel := progress.Channel()

	// start routine to listen to new messages
	go func() {
		// iterate any messages sent on the channel
		for msg := range channel {
			// create event object
			event := &models.Event{}
			event.Name = msg.Channel
			event.Timestamp = time.Now().UTC()
			event.Type = "task"

			if err := json.Unmarshal([]byte(msg.Payload), &event.Payload); err != nil {
				utils.LogError(errors.Wrap(err, "[SOCKET] error converting task payload to event"))
			}

			// marshal model to json string
			taskJSON, errJSON := json.Marshal(event)
			if errJSON != nil {
				utils.LogError(errors.Wrap(errJSON, "[SOCKET] error converting task object to string"))
			}

			socketConnection.BroadcastFilter(taskJSON, validSessionFilter)
		}
	}()
}

func writeSocketResponse(s *melody.Session, name string, msg interface{}) {
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
		return
	}

	s.Write(actionJSON)
}
