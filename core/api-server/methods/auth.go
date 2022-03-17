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
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/redis"
)

var ctx = context.Background()

func RedisLive() bool {
	// init redis connection
	redisConnection := redis.Instance()

	// check ping command
	_, err := redisConnection.Ping(ctx).Result()
	if err != nil {
		return false
	}

	return true
}

// RedisAuthentication godoc
// @Summary Login and get JWT token
// @Description login and get JWT token
// @Produce  json
// @Param user body response.LoginRequestJWT false "The user to login"
// @Success 200 {object} response.LoginResponseJWT
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /login [post]
// @Tags /login auth
func RedisAuthentication(username string, password string) error {
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

	// close redis connection
	redisConnection.Close()

	return nil
}

// RedisAuthorization godoc
// @Summary Login and remove JWT token
// @Description logout and remove JWT token
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /logout [post]
// @Tags /logout auth
func RedisAuthorization(username string, c *gin.Context) (models.UserAuthorizations, error) {
	var userAuthorizationsRedis models.UserAuthorizations

	// bypass auth for GET requests: // TODO
	if c.Request.Method == "GET" {
		return userAuthorizationsRedis, nil
	}

	// init redis connection
	redisConnection := redis.Instance()

	// define path
	var pathGet string
	var pathScan string
	parts := strings.Split(c.FullPath(), "/")
	entity := parts[2]

	// switch entity
	switch entity {
	case "cluster":
		pathGet = "cluster"
		pathScan = "cluster/roles/"
	case "node":
		pathGet = "node/" + c.Param("node_id")
		pathScan = "node/" + c.Param("node_id") + "/roles/"
	case "module":
		pathGet = "module/" + c.Param("module_id")
		pathScan = "module/" + c.Param("module_id") + "/roles/"
	}

	// get roles of current user: HGET roles/<username>.<entity> -> <role>
	role, errRedisRoleGet := redisConnection.HGet(ctx, "roles/"+username, pathGet).Result()

	// handle redis error
	if errRedisRoleGet != nil {
		return userAuthorizationsRedis, errRedisRoleGet
	}

	// get action for current role and entity: SMEMBERS <entity>/<reference>/roles/<role>
	actions, errRedisRoleScan := redisConnection.SMembers(ctx, pathScan+role).Result()

	// handle redis error
	if errRedisRoleScan != nil {
		return userAuthorizationsRedis, errRedisRoleScan
	}

	// compose user authorizations
	userAuthorizationsRedis.Username = username
	userAuthorizationsRedis.Role = role
	userAuthorizationsRedis.Actions = actions

	// close redis connection
	redisConnection.Close()

	// return auths
	return userAuthorizationsRedis, nil
}
