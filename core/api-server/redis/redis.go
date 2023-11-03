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
	"context"
	"strings"
	"strconv"
	"fmt"
	"github.com/NethServer/ns8-core/core/api-server/configuration"

	"github.com/go-redis/redis/v8"
)

func Instance() *redis.Client {
	// init client
	return redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     configuration.Config.RedisAddress,
		Username: configuration.Config.RedisUser,
		Password: configuration.Config.RedisPassword,
		DB:       0, // redis database
		OnConnect: setClientNameCallback,
	})
}

/*
 * Find the leader host address from Redis
 */
func GetLeaderHostAddress(ctx context.Context, rdb *redis.Client) string {
	var err error
	var leaderNodeId string
	var leaderHostAddress string

	leaderNodeId, err = rdb.HGet(ctx, "cluster/environment", "NODE_ID").Result()
	if err != nil {
		return ""
	}

	leaderHostAddress, err = rdb.HGet(ctx, "node/" + leaderNodeId + "/vpn", "endpoint").Result()
	if err != nil {
		return ""
	}

	separatorPosition := strings.Index(leaderHostAddress, ":")
	if separatorPosition < 1 {
		return ""
	}

	return leaderHostAddress[0:separatorPosition]
}

/*
 * Check the client with name "searchName" has an idle time less than "limit"
 */
func CheckClientIdle(ctx context.Context, rdb *redis.Client, searchName string, limit int) error {
	// Consider only "normal" clients: agents are among them
	val, err := rdb.Do(ctx, "CLIENT", "LIST", "TYPE", "normal").Result()
	if err != nil {
		return fmt.Errorf("Redis CLIENT LIST command failed (%v)", err)
	}
	// Split the CLIENT LIST output by lines and spaces
	for _, line := range(strings.Split(val.(string), "\n")) {
		var clientName string
		var clientIdle int
		var convError error
		// Seek idle= and name= fields among the fields list
		for _, field := range(strings.Split(line, " ")) {
			if strings.HasPrefix(field, "idle=") {
				clientIdle, convError = strconv.Atoi(field[5:])
				if convError != nil {
					return fmt.Errorf("Failed to parse Redis response (%v)", convError)
				}
			} else if strings.HasPrefix(field, "name=") {
				clientName = field[5:]
			}
		}
		if searchName == clientName {
			// agent found: let's check its idle time
			if clientIdle <= limit {
				return nil
			}
		}
	}
	return fmt.Errorf("Client name %s not found", searchName)
}

func setClientNameCallback (ctx context.Context, cn *redis.Conn) error {
	return cn.ClientSetName(ctx, configuration.Config.RedisUser).Err()
}
