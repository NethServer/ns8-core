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
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package methods

import (
	"context"

	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/NethServer/ns8-scratchpad/core/api-server/redis"
)

var ctx = context.Background()

func RedisAuth(username string, password string) error {
	// init redis connection
	redisConnection := redis.Instance()

	// execute redis auth: AUTH <username> <password>
	pipe := redisConnection.Pipeline()
	_, _ = pipe.AuthACL(ctx, username, password).Result()

	// check authentication
	_, errRedisAuth := pipe.Exec(ctx)
	if errRedisAuth != nil {
		return errRedisAuth
	}

	return nil
}

func GetUserAuthorizations(username string) (models.UserAuthorizations, error) {
	userAuthorizationsRedis := models.UserAuthorizations{
		Username: "edoardo",
		Role:     "admin",
		Actions:  []string{"create-users", "delete-users"},
	}

	// read auths from redis
	// userAuthorizationsRedis, err := ReadUserAuthorizationsRedis()
	// if err != nil {
	// 	return userAuthorizations, err
	// }

	// return auths
	return userAuthorizationsRedis, nil
}
