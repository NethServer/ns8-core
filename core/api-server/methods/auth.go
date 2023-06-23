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
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dgryski/dgoogauth"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/redis"
	"github.com/NethServer/ns8-core/core/api-server/response"
	"github.com/NethServer/ns8-core/core/api-server/utils"
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

func RedisAuthentication(username string, password string) error {
	if username == "default" {
		return errors.New("User is disabled")
	}

	// init redis connection
	redisClient := redis.Instance()

	// connection is only to check credentials: close it at the end
	defer redisClient.Close()

	// execute redis auth: AUTH <username> <password> and return its error, if any
	return redisClient.Conn(ctx).AuthACL(ctx, username, password).Err()
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

// CheckOTP godoc
// @Summary Check if OTP validates for the given username
func CheckOTP(username string, otp string) (bool) {
	secret := getUserSecret(username)
	return checkOtpBySecret(secret, otp)
}

func checkOtpBySecret(secret string, otp string) (bool) {
	// set OTP configuration
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	// verify OTP
	result, err := otpc.Authenticate(otp)
	if err == nil && result {
		return true
	}
	if err != nil {
		utils.LogError(err)
	}
	return false
}

// Set2FAStatus godoc
// @Summary Set the 2FA Redis flag, if OTP is verified
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /2FA [post]
// @Tags /2FA auth
func Set2FAStatus(c *gin.Context) {
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

	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)

	// check secret
	if len(jsonOTP.Secret) == 0 {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "Unexpected secret length",
			Data:    "",
		}))
		return
	}

	// Check if OTP is valid
	if ! checkOtpBySecret(jsonOTP.Secret, jsonOTP.OTP) {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "OTP token verification failed",
			Data:    "",
		}))
		return
	}

	// set the user 2FA secret: the next login request requires OTP.
	err := setUserSecret(username, jsonOTP.Secret)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "2FA set successfully",
		Data:    nil,
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
		fmt.Println("Failed to generate random secret for QRCode", err)
	}

	// convert to string
	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	// get claims from token
	claims := jwt.ExtractClaims(c)

	// define issuer
	account := claims["id"].(string)
	issuer := configuration.Config.Issuer

	// define URL
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		fmt.Println("Failed to parse URL for QRCode", err)
	}

	// add params
	URL.Path += "/" + issuer + ":" + account
	params := url.Values{}
	params.Add("secret", secretBase32)
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
		Data:    gin.H{"url": URL.String(), "key": secretBase32},
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
	username := claims["id"].(string)

	secret := getUserSecret(username)

	if secret == "" {
		c.JSON(http.StatusOK, structs.Map(response.StatusOK{
			Code:    200,
			Message: "2FA not set for this user",
			Data:    false,
		}))
		return
	} else {
		c.JSON(http.StatusOK, structs.Map(response.StatusOK{
			Code:    200,
			Message: "2FA set for this user",
			Data:    true,
		}))
		return
	}
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
	username := claims["id"].(string)

	// init redis connection
	redisConnection := redis.Instance()

	// set 2FA to disabled
	if err := redisConnection.HSet(ctx, "user/" + username, "2fa", "").Err(); err != nil {
		utils.LogError(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "2FA revocate successfully",
		Data:    "",
	}))
}

// BasicAuthModule godoc
// @Summary Use Basic HTTP auth for module
// @Description for a specific action, check basic auth for module
// @Produce  json
// @Header 200 {string} Authorization "Basic base64(username:password)"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Failure 401 {object} response.StatusUnauthorized{code=int,message=string,data=object}
// @Failure 403 {object} response.StatusForbidden{code=int,message=string,data=object}
// @Router /module/{module_id}/http-basic/{action} [get]
// @Tags /module basic_auth
func BasicAuthModule(c *gin.Context) {
	// get http basic credentials
	username, password, _ := c.Request.BasicAuth()

	// init redis connection
	redisConnection := redis.Instance()

	// close redis connection
	defer redisConnection.Close()

	// check authentication
	err := RedisAuthentication(username, password)
	if err != nil {
		utils.LogError(errors.Wrap(err, "[BASIC AUTH] redis authentication failed for user "+username))

		// response
		c.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		c.JSON(http.StatusUnauthorized, structs.Map(response.StatusUnauthorized{
			Code:    401,
			Message: "basic auth failed. username or password invalid",
			Data:    "",
		}))
		return
	}

	// check authorizations
	pathGet := "module/" + c.Param("module_id")
	pathScan := "module/" + c.Param("module_id") + "/roles/"

	// get roles of current user: HGET roles/<username>.<entity> -> <role>
	role, errRedisRoleGet := redisConnection.HGet(ctx, "roles/"+username, pathGet).Result()

	// handle redis error
	if errRedisRoleGet != nil {
		// response
		c.JSON(http.StatusForbidden, structs.Map(response.StatusForbidden{
			Code:    403,
			Message: "basic auth failed. module not found",
			Data:    "",
		}))
		return
	}

	// get action for current role and entity: SMEMBERS <entity>/<reference>/roles/<role>
	actions, errRedisRoleScan := redisConnection.SMembers(ctx, pathScan+role).Result()

	// handle redis error
	if errRedisRoleScan != nil {
		// response
		c.JSON(http.StatusForbidden, structs.Map(response.StatusForbidden{
			Code:    403,
			Message: "basic auth failed. action not found",
			Data:    "",
		}))
		return
	}

	// verify wildcard actions
	actionAllowed := false
	for _, authorizedAction := range actions {
		actionAllowed, _ = filepath.Match(authorizedAction, c.Param("action"))
		if actionAllowed {
			break
		}
	}

	if !actionAllowed {
		// response
		c.JSON(http.StatusForbidden, structs.Map(response.StatusForbidden{
			Code:    403,
			Message: "basic auth failed. action not allowed",
			Data:    "",
		}))
		return
	}

	// response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "basic auth ok",
		Data:    actions,
	}))

}

func setUserSecret(username string, secret string) (error) {
	// init redis connection
	redisConnection := redis.Instance()
	err := redisConnection.HSet(ctx, "user/" + username, "2fa", secret).Err()
	if err != nil {
		return err
	}
	return nil
}

func getUserSecret(username string) string {
	redisConnection := redis.Instance()
	result, err := redisConnection.HGet(ctx, "user/" + username, "2fa").Result()
	if err != nil {
		return ""
	}
	if len(result) != 32 {
		return ""
	}
	return result
}

func Needs2faCheck(username string) bool {
	return getUserSecret(username) != ""
}
