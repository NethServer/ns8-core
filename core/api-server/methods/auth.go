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
	"encoding/base32"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dgryski/dgoogauth"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/redis"
	"github.com/NethServer/ns8-core/core/api-server/response"
	"github.com/NethServer/ns8-core/core/api-server/socket"
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

// OTPVerify godoc
// @Summary After Login verify OTP for 2FA
// @Description logout and remove JWT token
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /2FA/otp-verify [post]
// @Tags /2FA/otp-verify auth
func OTPVerify(c *gin.Context) {
	// get payload
	var jsonOTP models.OTPJson
	if err := c.ShouldBindBodyWith(&jsonOTP, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "request fields malformed",
			Data:    err.Error(),
		}))
		return
	}

	// verify JWT
	if !socket.ValidateAuth(jsonOTP.Token) {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "JWT token invalid",
			Data:    "",
		}))
		return
	}

	// get secret for the user
	secret := getUserSecret(jsonOTP.Username)

	// check secret
	if len(secret) == 0 {
		c.JSON(http.StatusNotFound, structs.Map(response.StatusNotFound{
			Code:    404,
			Message: "User secret not found",
			Data:    "",
		}))
		return
	}

	// set OTP configuration
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	// verifiy OTP
	result, err := otpc.Authenticate(jsonOTP.OTP)
	if err != nil || !result {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "OTP token invalid",
			Data:    "",
		}))
		return
	}

	// set auth token to valid
	if !SetTokenValidation(jsonOTP.Username, jsonOTP.Token) {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "Token validation set error",
			Data:    "",
		}))
		return
	}

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "OTP verified",
		Data:    jsonOTP.Token,
	}))
}

// QRCode godoc
// @Summary Generate QRCode for 2FA
// @Description use this API to generate a valid qr-code for google auth
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /2FA/qrcode [get]
// @Tags /2FA/qrcode auth
func QRCode(c *gin.Context) {
	// generate random secret
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}

	// convert to string
	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	// get claims from token
	claims := jwt.ExtractClaims(c)

	// define issuer
	account := claims["id"].(string)
	issuer := configuration.Config.Issuer

	// set secret for user
	result, setSecret := setUserSecret(account, secretBase32)
	if !result {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "User secret set error",
			Data:    "",
		}))
		return
	}

	// define URL
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}

	// add params
	URL.Path += "/" + issuer + ":" + account
	params := url.Values{}
	params.Add("secret", setSecret)
	params.Add("issuer", issuer)
	params.Add("algorithm", "SHA1")
	params.Add("digits", "6")
	params.Add("period", "30")

	// print url
	URL.RawQuery = params.Encode()

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "QR code string",
		Data:    URL.String(),
	}))
}

// Get2FAStatus godoc
// @Summary Get current 2FA status for user
// @Description get current 2FA status for user
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /2FA [get]
// @Tags /2FA auth
func Get2FAStatus(c *gin.Context) {
	// get claims from token
	claims := jwt.ExtractClaims(c)

	// init redis connection
	redisConnection := redis.Instance()

	// get secret
	status, errRedis2FAGet := redisConnection.HGet(ctx, "user/"+claims["id"].(string), "2fa").Result()

	// handle redis error
	if errRedis2FAGet != nil {
		c.JSON(http.StatusNotFound, structs.Map(response.StatusNotFound{
			Code:    404,
			Message: "2FA not set for this user",
			Data:    nil,
		}))
		return
	}

	// response
	var message = "2FA set for this user"
	if !(status == "1") {
		message = "2FA not set for this user"
	}

	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: message,
		Data:    status == "1",
	}))
}

// Set2FAStatus godoc
// @Summary Set current 2FA status for user
// @Description set current 2FA status for user
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /2FA [post]
// @Tags /2FA auth
func Set2FAStatus(c *gin.Context) {
	// get payload
	var status2FA models.Status2FA
	if err := c.ShouldBindBodyWith(&status2FA, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "request fields malformed",
			Data:    err.Error(),
		}))
		return
	}

	// get claims from token
	claims := jwt.ExtractClaims(c)

	// init redis connection
	redisConnection := redis.Instance()

	// set auth token to valid
	errRedis2FASet := redisConnection.HSet(ctx, "user/"+claims["id"].(string), "2fa", status2FA.Status)

	// revocate all tokens, if true
	var existsSecret = getUserSecret(claims["id"].(string))
	if status2FA.Status && len(existsSecret) == 0 {
		if err := redisConnection.Del(ctx, "user/"+claims["id"].(string)+"/tokens").Err(); err != nil {
			c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
				Code:    403,
				Message: "Error in revocate all user tokens",
				Data:    nil,
			}))
			return
		}
	}

	// check error
	if errRedis2FASet.Err() != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    403,
			Message: "Error in set 2FA for user",
			Data:    nil,
		}))
		return
	}

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "2FA set successfully",
		Data:    "",
	}))

}

// Del2FAStatus godoc
// @Summary Revocate secret for 2FA
// @Description revocate secret for 2FA
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 500 {object} response.StatusInternalServerError{code=int,message=string,data=object}
// @Router /2FA [delete]
// @Tags /2FA auth
func Del2FAStatus(c *gin.Context) {
	// get claims from token
	claims := jwt.ExtractClaims(c)

	// init redis connection
	redisConnection := redis.Instance()

	// revocate secret
	if errRevocate := redisConnection.Del(ctx, "secrets/"+claims["id"].(string), "2fa").Err(); errRevocate != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    403,
			Message: "Error in revocate 2FA for user",
			Data:    nil,
		}))
		return
	}

	// delete tokens
	if errDelete := redisConnection.Del(ctx, "user/"+claims["id"].(string)+"/tokens").Err(); errDelete != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    403,
			Message: "Error in revocate all user tokens",
			Data:    nil,
		}))
		return
	}

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "2FA revocate successfully",
		Data:    "",
	}))
}

func setUserSecret(username string, secret string) (bool, string) {
	// init redis connection
	redisConnection := redis.Instance()

	// check if secret is already set
	createdSecret, _ := redisConnection.HGet(ctx, "secrets/"+username, "2fa").Result()

	// check error
	if len(createdSecret) == 0 {
		// set auth token to valid
		errRedisTokenSet := redisConnection.HSet(ctx, "secrets/"+username, "2fa", secret)

		// check error
		if errRedisTokenSet.Err() != nil {
			return false, ""
		}

		return true, secret
	}

	return true, createdSecret
}

func SetTokenValidation(username string, token string) bool {
	// init redis connection
	redisConnection := redis.Instance()

	// set auth token to valid
	errRedisTokenSet := redisConnection.SAdd(ctx, "user/"+username+"/tokens", token)

	// check error
	if errRedisTokenSet.Err() != nil {
		return false
	}

	return true
}

func getUserSecret(username string) string {
	// init redis connection
	redisConnection := redis.Instance()

	// get secret
	secret, errRedisSecretGet := redisConnection.HGet(ctx, "secrets/"+username, "2fa").Result()

	// handle redis error
	if errRedisSecretGet != nil {
		return ""
	}

	// read from redis
	return secret
}

func CheckTokenValidation(username string, token string) bool {
	// init redis connection
	redisConnection := redis.Instance()

	// get token list
	tokens, errRedisTokenScan := redisConnection.SMembers(ctx, "user/"+username+"/tokens").Result()

	// handle redis error
	if errRedisTokenScan != nil {
		return false
	}

	// loop tokens
	var valid = false
	for _, t := range tokens {
		if t == token {
			valid = true
		}
	}

	return valid
}

func Check2FA(username string) bool {
	// init redis connection
	redisConnection := redis.Instance()

	// get status
	result, errRedisStatusGet := redisConnection.HGet(ctx, "user/"+username, "2fa").Result()

	// handle redis error
	if errRedisStatusGet != nil {
		return false
	}

	return result == "1"

}
