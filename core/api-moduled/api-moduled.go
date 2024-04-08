/*
 * Copyright (C) 2023 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/NethServer/ns8-core/core/api-moduled/validation"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var logger *log.Logger

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
	viper.SetEnvPrefix("AMLD")
	viper.SetDefault("handler_dir", "./handlers/")
	viper.SetDefault("public_dir", "./public/")
	viper.SetDefault("bind_address", ":9313")
	viper.SetDefault("id_key", "uid")
	viper.SetDefault("jwt_secret", "")
	viper.SetDefault("jwt_timeout", time.Hour*4)
	viper.SetDefault("jwt_token_lookup", "header: Authorization")
	viper.SetDefault("jwt_realm", "api-moduled")
	viper.SetDefault("export_env", "")
	viper.AutomaticEnv()

	logger = log.New(os.Stderr, "", 0)

	if len(viper.GetString("jwt_secret")) == 0 {
		logger.Println(SD_WARNING + "AMLD_JWT_SECRET environment variable is empty! JWT tokens are unsecure.")
	}

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

	ijwt := createJwtInstance(viper.GetString("handler_dir"))

	api := router.Group("/api")
	api.GET("/auth", basicAuth)
	api.POST("/login", ijwt.LoginHandler)
	api.Use(ijwt.MiddlewareFunc()) // next API route definitions require the Authorization header
	api.POST("/logout", ijwt.LogoutHandler)
	api.POST("/:handler", apiPostHandler)

	router.NoRoute(ijwt.MiddlewareFunc(), func(ginCtx *gin.Context) {
		ginCtx.JSON(http.StatusNotFound, gin.H{
			"code":   http.StatusNotFound,
			"status": "Not found",
		})
	})

	router.Use(static.Serve("/", static.LocalFile(viper.GetString("public_dir"), false)))

	router.Run(viper.GetString("bind_address"))
}

func prepareEnvironment(ginCtx *gin.Context) []string {
	claims := jwt.ExtractClaims(ginCtx)
	jclaims, _ := json.Marshal(claims)
	jwt_id, _ := claims[viper.GetString("id_key")].(string)
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"JWT_ID=" + jwt_id,
		"JWT_CLAIMS=" + string(jclaims),
	}
	for _, varname := range strings.Split(viper.GetString("export_env"), " ") {
		env = append(env, varname+"="+os.Getenv(varname))
	}
	return env
}

func apiPostHandler(ginCtx *gin.Context) {

	pHandler := path.Base(ginCtx.Param("handler"))
	handlerDir := path.Clean(viper.GetString("handler_dir")+"/"+pHandler) + "/"

	// Check if handler directory exists, return 404 otherwise:
	if _, err := os.Stat(handlerDir); err != nil {
		ginCtx.Error(err)
		ginCtx.JSON(http.StatusNotFound, gin.H{
			"code":   http.StatusNotFound,
			"status": "Not found",
		})
		return
	}

	requestBytes, _ := io.ReadAll(ginCtx.Request.Body)

	//
	// INPUT validation
	//
	if _, err := os.Stat(handlerDir + "validate-input.json"); err == nil {
		errData, errInfo := validation.ValidatePayload(handlerDir+"validate-input.json", requestBytes)
		if errInfo != nil {
			logger.Println(SD_ERR+"Input validation error:", errInfo)
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Internal server error",
			})
			return
		} else if len(errData) > 0 {
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"code":   http.StatusBadRequest,
				"status": "Bad request",
				"error":  errData,
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
	cmd.Env = prepareEnvironment(ginCtx)
	responseBytes, cerr := cmd.Output()
	if cerr != nil {
		logger.Println(SD_ERR+"Error from", cmd.String()+":", cerr)
	}

	//
	// OUTPUT validation
	//
	if _, err := os.Stat(handlerDir + "validate-output.json"); err == nil {
		errData, errInfo := validation.ValidatePayload(handlerDir+"validate-output.json", responseBytes)
		if errInfo != nil {
			logger.Println(SD_ERR+"Output validation error:", errInfo)
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Internal server error",
			})
			return
		} else if len(errData) > 0 {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"code":   http.StatusInternalServerError,
				"status": "Internal server error",
				"error":  errData,
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
		logger.Println(SD_ERR+"JSON Unmarshal() error:", jerr)
		logger.Println(SD_DEBUG+"Response buffer", string(responseBytes[:]))
		ginCtx.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"status": "Internal server error",
		})
		return
	}
	ginCtx.JSON(http.StatusOK, responsePayload)
}

func createJwtInstance(baseHandlerDir string) *jwt.GinJWTMiddleware {

	jwtInstance, errDefine := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         viper.GetString("jwt_realm"),
		Key:           []byte(viper.GetString("jwt_secret")),
		Timeout:       viper.GetDuration("jwt_timeout"),
		IdentityKey:   viper.GetString("id_key"),
		TokenLookup:   viper.GetString("jwt_token_lookup"),
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,

		Authenticator: func(ginCtx *gin.Context) (interface{}, error) {
			requestBytes, _ := io.ReadAll(ginCtx.Request.Body)

			//
			// INPUT validation
			//
			if _, err := os.Stat(baseHandlerDir + "/login/validate-input.json"); err == nil {
				errData, errInfo := validation.ValidatePayload(baseHandlerDir+"/login/validate-input.json", requestBytes)
				if errInfo != nil {
					logger.Println(SD_ERR+"Input validation error:", errInfo)
					return nil, jwt.ErrFailedAuthentication
				} else if len(errData) > 0 {
					logger.Println(SD_ERR+"Input validation error:", errData)
					return nil, jwt.ErrFailedAuthentication
				}
			}

			///
			/// Run login
			///
			responseBytes, cerr, exitCode, cmdString := runLogin(baseHandlerDir, requestBytes, ginCtx)

			///
			/// Check error
			///
			if cerr != nil {
				// Consider exit code 2-7 as a bad login attempt that must
				// be logged properly. The first logged IP should be
				// ignored if the service runs behind a proxy
				if exitCode > 1 && exitCode < 8 {
					logger.Printf(SD_NOTICE+"Bad login attempt! Exit code %d; remote address: %s\n", exitCode, ginCtx.ClientIP())
					return nil, jwt.ErrFailedAuthentication
				} else {
					logger.Printf(SD_ERR+"Error from %s: exit code %d\n", cmdString, exitCode)
					return nil, errors.New("internal error")
				}
			}
			var responsePayload gin.H
			jerr := json.Unmarshal(responseBytes, &responsePayload)
			if jerr != nil {
				logger.Println(SD_ERR+"Login response error: ", jerr)
				return nil, errors.New("internal error")
			}
			ginCtx.Set("JWT_PAYLOAD", jwt.MapClaims(responsePayload))
			return responsePayload, nil // Authentication is successful
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if claims, ok := data.(gin.H); ok {
				return jwt.MapClaims(claims)
			}
			logger.Println(SD_CRIT + "PayloadFunc error: login output cannot be converted to jwt claims")
			return nil
		},
		LoginResponse: func(ginCtx *gin.Context, code int, token string, texp time.Time) {
			ginCtx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"token":   token,
				"expire":  texp.Format(time.RFC3339),
				"message": "login succeeded",
				"claims":  jwt.ExtractClaims(ginCtx),
			})
		},
		Authorizator: func(data interface{}, ginCtx *gin.Context) bool {
			claims := jwt.ExtractClaims(ginCtx)

			if ginCtx.FullPath() == "/api/logout" {
				return true // logout is always allowed
			} else if !strings.HasPrefix(ginCtx.FullPath(), "/api/") {
				return true // static paths are always allowed
			}

			scopeValue, scopeClaimExists := claims["scope"]
			if !scopeClaimExists {
				// The token scope is unlimited: auth ok.
				return true
			}

			if scope, ok := scopeValue.([]interface{}); ok {
				// Check if the wanted handler is in the scope
				for _, elemValue := range scope {
					elem, ok := elemValue.(string)
					if ok && ginCtx.Param("handler") == elem {
						// The token scope matches request path: auth ok.
						return true
					}
				}
			} else {
				logger.Printf(SD_CRIT+"Authorizator: invalid token scope type \"%T\". Value: %#v\n", ginCtx.FullPath(), scope, scope)
			}

			return false
		},
	})

	// check middleware errors
	if errDefine != nil {
		logger.Println(SD_ERR+"JWT middleware definition error:", errDefine)
	}

	// init middleware
	errInit := jwtInstance.MiddlewareInit()

	// check error on initialization
	if errInit != nil {
		logger.Println(SD_ERR+"JWT middleware initialization error:", errInit)
	}

	return jwtInstance
}

func runLogin(baseHandlerDir string, requestBytes []byte, ginCtx *gin.Context) ([]byte, error, int, string) {
	cmd := exec.Command(baseHandlerDir + "/login/post")
	cmd.Stdin = bytes.NewReader(requestBytes)
	cmd.Stderr = os.Stderr
	cmd.Env = prepareEnvironment(ginCtx)
	responseBytes, cerr := cmd.Output()

	return responseBytes, cerr, cmd.ProcessState.ExitCode(), cmd.String()
}

func basicAuth(ginCtx *gin.Context) {
	// get basic auth credentials
	user, password, hasAuth := ginCtx.Request.BasicAuth()

	if !hasAuth {
		logger.Println(SD_ERR + "Missing or invalid auth header")
		ginCtx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "incorrect Username or Password",
		})
		return
	}

	// create request
	requestBytes := []byte(`{"username":"` + user + `","password":"` + password + `"}`)

	// execute login
	_, cerr, _, _ := runLogin(viper.GetString("handler_dir"), requestBytes, ginCtx)

	// check error
	if cerr != nil {
		ginCtx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "incorrect Username or Password",
		})
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "login succeeded",
	})

}
