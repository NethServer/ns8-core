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

package redis

import (
	"github.com/NethServer/ns8-scratchpad/core/api-server/configuration"

	"github.com/go-redis/redis/v8"
)

var redisConnection *redis.Client

func Instance() *redis.Client {
	// init client
	redisConnection = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     configuration.Config.RedisAddress,
		Username: configuration.Config.RedisUser,
		Password: configuration.Config.RedisPassword,
		DB:       0, // redis database
		OnConnect: setClientNameCallback,
	})

	return redisConnection
}

func setClientNameCallback (ctx context.Context, cn *redis.Conn) error {
	return cn.ClientSetName(ctx, configuration.Config.RedisUser).Err()
}
