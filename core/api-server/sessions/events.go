/*
 * Copyright (C) 2026 Nethesis S.r.l.
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
 */

package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/redis"
	"github.com/NethServer/ns8-core/core/api-server/utils"
)

// CredentialsChangedEvent is published by cluster actions when a user
// password is changed by an administrator, 2FA is reset, or the user is
// removed
const CredentialsChangedEvent = "cluster/event/user-credentials-changed"

type credentialsChanged struct {
	User      string  `json:"user"`
	NotBefore float64 `json:"not_before"`
}

// ListenRevocations ends the sessions of users whose credentials have
// changed. It blocks: run it in its own goroutine.
func ListenRevocations() {
	ctx := context.Background()
	pubsub := redis.Instance().Subscribe(ctx, CredentialsChangedEvent)
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var event credentialsChanged
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil || event.User == "" {
			utils.LogError(errors.New("[AUTH] invalid " + CredentialsChangedEvent + " payload"))
			continue
		}
		notBefore := time.Now()
		if event.NotBefore > 0 {
			notBefore = time.Unix(0, int64(event.NotBefore*1e9))
		}
		count, err := RevokeUser(event.User, notBefore)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUTH] session revocation failed for user "+event.User))
			continue
		}
		if count > 0 {
			utils.LogError(fmt.Errorf("[AUTH] sessions-revoked user=%s count=%d", event.User, count))
		}
	}
}
