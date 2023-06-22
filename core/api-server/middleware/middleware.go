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

package middleware

import (
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/NethServer/ns8-core/core/api-server/audit"
	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/methods"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/response"
	"github.com/NethServer/ns8-core/core/api-server/utils"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Otp string `form:"otp" json:"otp,omitempty"`
}

var jwtMiddleware *jwt.GinJWTMiddleware
var identityKey = "id"

func InstanceJWT() *jwt.GinJWTMiddleware {
	if jwtMiddleware == nil {
		jwtMiddleware := InitJWT()
		return jwtMiddleware
	}
	return jwtMiddleware
}

func InitJWT() *jwt.GinJWTMiddleware {
	// define jwt middleware
	authMiddleware, errDefine := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "nethserver",
		Key:         []byte(configuration.Config.Secret),
		Timeout:     time.Hour * 24 * 14, // 2 weeks
		IdentityKey: identityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// check login credentials exists
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// set login credentials
			username := loginVals.Username
			password := loginVals.Password
			otpValue := loginVals.Otp

			// check if redis is running
			if !methods.RedisLive() {
				return nil, errors.New("redis is not running")
			}

			// check login with redis
			err := methods.RedisAuthentication(username, password)
			if err != nil {
				utils.LogError(errors.Wrap(err, "[AUTH] redis authentication failed for user "+username))

				// store login action
				auditData := models.Audit{
					ID:        0,
					User:      username,
					Action:    "login-fail",
					Data:      "",
					Timestamp: time.Now().UTC(),
				}
				audit.Store(auditData)

				return nil, jwt.ErrFailedAuthentication
			}

			otpNeed := methods.Needs2faCheck(username)
			if otpNeed && otpValue == "" {
				return nil, errors.New("missing OTP value")
			}

			otpPassClaim := false
			if otpNeed {
				if ! methods.CheckOTP(username, otpValue) {
					err := errors.New("OTP check failed")
					// store login action
					auditData := models.Audit{
						ID:        0,
						User:      username,
						Action:    "login-fail",
						Data:      "",
						Timestamp: time.Now().UTC(),
					}
					audit.Store(auditData)
					return nil, err
				}
				otpPassClaim = true // claim that 2FA is enabled and used
			}

			// Login is successful. Middleware returns a JWT.
			// store login action
			auditData := models.Audit{
				ID:        0,
				User:      username,
				Action:    "login-ok",
				Data:      "",
				Timestamp: time.Now().UTC(),
			}
			audit.Store(auditData)

			// return user auth model
			return &models.UserAuthorizations{
				Username: username,
				OtpPass: otpPassClaim, // true if OTP passed, false if OTP is not needed
			}, nil

		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// read current user
			if user, ok := data.(*models.UserAuthorizations); ok {
				// create claims map
				return jwt.MapClaims{
					identityKey: user.Username,
					"role":      "",
					"actions":   []string{},
					"2fa":       user.OtpPass, // OTP value check result: true=success, false=not needed
				}
			}

			// return claims map
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			// handle identity and extract claims
			claims := jwt.ExtractClaims(c)

			// handle authorizations
			userAuthorizations, err := methods.RedisAuthorization(claims[identityKey].(string), c)
			if err != nil {
				utils.LogError(errors.Wrap(err, "[AUTH] error retrieving user authorizations"))
			}

			// create user object
			user := &models.UserAuthorizations{
				Username: claims[identityKey].(string),
				Role:     userAuthorizations.Role,
				Actions:  userAuthorizations.Actions,
			}

			// return user
			return user
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// bypass auth for GET requests: // TODO
			if c.Request.Method == "GET" {
				return true
			}

			// bypass for 2FA apis
			if c.Request.Method == "POST" && c.Request.RequestURI == "/api/2FA" {
				return true
			}
			if c.Request.Method == "DELETE" && c.Request.RequestURI == "/api/2FA" {
				return true
			}

			// extract data payload and check authorizations
			if v, ok := data.(*models.UserAuthorizations); ok {
				authorizedActions := v.Actions

				// extract task obj
				var jsonTask models.Task
				errJson := c.ShouldBindBodyWith(&jsonTask, binding.JSON)
				if errJson != nil {
					utils.LogError(errors.Wrap(errJson, "[AUTH] error unmarshalling task obj"))
					return false
				}

				// check action authorization
				actionAllowed := false
				for _, authorizedAction := range authorizedActions {
					actionAllowed, _ = filepath.Match(authorizedAction, jsonTask.Action)
					if actionAllowed {
						break
					}
				}

				// store auth action
				auditData := models.Audit{
					ID:        0,
					User:      data.(*models.UserAuthorizations).Username,
					Action:    "auth-ok",
					Data:      "",
					Timestamp: time.Now().UTC(),
				}
				audit.Store(auditData)

				return actionAllowed
			}

			// store auth action
			auditData := models.Audit{
				ID:        0,
				User:      data.(*models.UserAuthorizations).Username,
				Action:    "auth-fail",
				Data:      "",
				Timestamp: time.Now().UTC(),
			}
			audit.Store(auditData)

			// not authorized
			return false
		},
		LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
			c.JSON(200, gin.H{"code": 200, "expire": t, "token": token})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(200, gin.H{"code": 200})
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if message == "redis is not running" {
				c.JSON(503, structs.Map(response.StatusServiceUnavailable{
					Code:    503,
					Message: message,
					Data:    nil,
				}))
				return
			}

			c.JSON(code, structs.Map(response.StatusUnauthorized{
				Code:    code,
				Message: message,
				Data:    nil,
			}))
			return
		},
		TokenLookup:   "header: Authorization, query: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	// check middleware errors
	if errDefine != nil {
		utils.LogError(errors.Wrap(errDefine, "[AUTH] middleware definition error"))
	}

	// init middleware
	errInit := authMiddleware.MiddlewareInit()

	// check error on initialization
	if errInit != nil {
		utils.LogError(errors.Wrap(errInit, "[AUTH] middleware initialization error"))
	}

	// return object
	return authMiddleware
}
