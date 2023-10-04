/*
 * Copyright (C) 2023 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"log"
	"os/exec"
	"encoding/json"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/NethServer/ns8-core/core/api-moduled/validation"
	"time"
)

var logger *log.Logger
var identityKey string

// Reference: https://www.man7.org/linux/man-pages/man3/sd-daemon.3.html
const (
	SD_EMERG   = "<0>" /* system is unusable */
	SD_ALERT   = "<1>" /* action must be taken immediately */
	SD_CRIT    = "<2>" /* critical conditions */
	SD_ERR     = "<3>" /* error conditions */
	SD_WARNING = "<4>" /* warning conditions */
	SD_NOTICE  = "<5>" /* normal but significant condition */
	SD_INFO    = "<6>" /* informational */
	SD_DEBUG   = "<7>" /* debug-level messages */
)

func main() {

	identityKey = os.Getenv("JWT_IDKEY")
	if len(identityKey) == 0 {
		identityKey = "uid"
	}

	logger = log.New(os.Stderr, "", 0)

	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter),
		gin.Recovery(),
		gzip.Gzip(gzip.DefaultCompression),
	)

	// Allow cross-origin requests in DebugMode
	// For development only, set in environment: GIN_MODE=debug
	if gin.Mode() == gin.DebugMode {
		corsConf := cors.DefaultConfig()
		corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Accept"}
		corsConf.AllowAllOrigins = true
		router.Use(cors.New(corsConf))
	}

	ijwt := createJwtInstance("./handlers") // XXX hardcoded handlers root path

	api := router.Group("/api")
	api.POST("/login", ijwt.LoginHandler)
	api.Use(ijwt.MiddlewareFunc()) // next API route definitions require the Authorization header
	api.POST("/logout", ijwt.LogoutHandler)
	mapHandlers(api, "./handlers") // XXX hardcoded handlers root path

	router.NoRoute(ijwt.MiddlewareFunc(), func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusNotFound, gin.H{
			"code":		http.StatusNotFound,
			"status": 	"Not found",
		})
	})

	router.Static("/", "public")

	router.Run(":9313") // XXX hardcoded listen address
}

func mapHandlers(routerGroup *gin.RouterGroup, baseHandlerDir string) {

	entries, err := os.ReadDir(baseHandlerDir)
	if err != nil {
		logger.Println(SD_ERR + "mapHandlers:", err)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "login" {
			var handlerDir = baseHandlerDir + "/" + entry.Name() + "/"
			if _, err := os.Stat(handlerDir + "post"); err == nil {
				routerGroup.POST(entry.Name(), func (ginCtx *gin.Context) {
					requestBytes, _ := io.ReadAll(ginCtx.Request.Body)

					//
					// INPUT validation
					//
					if _, err := os.Stat(handlerDir + "validate-input.json") ; err == nil {
						errData, errInfo := validation.ValidatePayload(handlerDir + "validate-input.json", requestBytes)
						if errInfo != nil {
							logger.Println(SD_ERR + "Input validation error", errInfo)
							ginCtx.JSON(http.StatusInternalServerError, gin.H{
								"code":	http.StatusInternalServerError,
								"status":	"Internal server error",
							})
							return
						} else if len(errData) > 0 {
							ginCtx.JSON(http.StatusBadRequest, gin.H{
								"code":		http.StatusBadRequest,
								"status":	"Bad request",
								"error":	errData,
							})
							return
						}
					}

					///
					/// COMMAND execution
					///
					cmd := exec.Command(handlerDir + "post")
					cmd.Stdin = bytes.NewReader(requestBytes)
					cmd.Stderr = os.Stderr
					cmd.Env = append(os.Environ(), // XXX sanitize the environment!!
						"JWT_USER=unknown", // XXX read user from jwt
						"JWT_ROLES=role1,role2", // XXX read roles from jwt
						"JWT_SCOPES=scope1,scope2", // XXX read scopes from jwt
					)
					responseBytes, cerr := cmd.Output()
					if cerr != nil {
						logger.Println(SD_ERR + "Error from", cmd.String() + ":", cerr)
					}

					//
					// OUTPUT validation
					//
					if _, err := os.Stat(handlerDir + "validate-output.json"); err == nil {
						errData, errInfo := validation.ValidatePayload(handlerDir + "validate-output.json", responseBytes)
						if errInfo != nil {
							logger.Println(SD_ERR + "Output validation error", errInfo)
							ginCtx.JSON(http.StatusInternalServerError, gin.H{
								"code":	http.StatusInternalServerError,
								"status":	"Internal server error",
							})
							return
						} else if len(errData) > 0 {
							ginCtx.JSON(http.StatusInternalServerError, gin.H{
								"code":		http.StatusInternalServerError,
								"status":	"Internal server error",
								"error":	errData,
							})
							return
						}
					}

					///
					/// Response output
					///
					var responsePayload gin.H
					jerr := json.Unmarshal(responseBytes, &responsePayload)
					if jerr != nil {
						logger.Println(SD_ERR + "JSON Unmarshal() error:", jerr)
						logger.Println(SD_DEBUG + "Response buffer", string(responseBytes[:]))
						ginCtx.JSON(http.StatusInternalServerError, gin.H{
							"code":	http.StatusInternalServerError,
							"status":	"Internal server error",
						})
						return
					}
					ginCtx.JSON(http.StatusOK, responsePayload)
					return
				})
			}
		}
	}
}

func createJwtInstance(baseHandlerDir string) *jwt.GinJWTMiddleware {
	// define jwt middleware

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if len(jwtSecretKey) == 0 {
		logger.Println(SD_WARNING + "Missing configuration: set JWT_SECRET in the environment.")
	}

	jwtInstance, errDefine := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       	"nethserver",
		Key:         	[]byte(jwtSecretKey),
		Timeout:     	time.Hour, // XXX get JWT timeout
		IdentityKey:	identityKey,
		TokenLookup:  	"header: Authorization, query: jwt",
		TokenHeadName:	"Bearer",
		TimeFunc:     	time.Now,

		Authenticator: func(ginCtx *gin.Context) (interface{}, error) {
			requestBytes, _ := io.ReadAll(ginCtx.Request.Body)
			cmd := exec.Command(baseHandlerDir + "/login/post")
			cmd.Stdin = bytes.NewReader(requestBytes)
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ() // XXX sanitize the environment!!
			responseBytes, cerr := cmd.Output()
			if cerr != nil {
				logger.Println(SD_ERR + "Error from", cmd.String() + ":", cerr)
				return nil, jwt.ErrFailedAuthentication
			}
			var responsePayload gin.H
			jerr := json.Unmarshal(responseBytes, &responsePayload)
			if jerr != nil {
				logger.Println(SD_ERR + "Login response error: ", jerr)
				return nil, jwt.ErrFailedAuthentication
			}
			return responsePayload, nil	// Authentication is successful
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if claims, ok := data.(gin.H) ; ok {
				return jwt.MapClaims(claims)
			}
			logger.Println(SD_CRIT + "PayloadFunc error: login output cannot be converted to jwt claims")
			return nil
		},
	})

	// check middleware errors
	if errDefine != nil {
		logger.Println(SD_ERR + "JWT middleware definition error:", errDefine)
	}

	// init middleware
	errInit := jwtInstance.MiddlewareInit()

	// check error on initialization
	if errInit != nil {
		logger.Println(SD_ERR + "JWT middleware initialization error:", errInit)
	}

	return jwtInstance
}
