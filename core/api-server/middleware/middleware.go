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
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/time/rate"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/core"
	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/NethServer/ns8-core/core/api-server/audit"
	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/methods"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/response"
	"github.com/NethServer/ns8-core/core/api-server/sessions"
	"github.com/NethServer/ns8-core/core/api-server/utils"
)

// BodyLimit caps the number of bytes read from the request body before
// binding, so unauthenticated routes cannot force large allocations
// ahead of credential validation (see GHSA-3v6g-pgp9-cmm7).
func BodyLimit(limitBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitBytes)
		c.Next()
	}
}

const rateLimiterStaleAfter = 3 * time.Minute
const rateLimiterCleanupInterval = time.Minute

type rateLimiterVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter throttles requests per client IP with a token-bucket limiter,
// so an unauthenticated route cannot be flooded with enough concurrent
// requests to exhaust memory before BodyLimit's per-request cap can help
// (BodyLimit bounds one request's body, not how many requests run at once).
// rps is the sustained rate and burst the number of requests allowed instantly.
func RateLimiter(rps rate.Limit, burst int) gin.HandlerFunc {
	// A zero-capacity bucket can never accumulate the whole token a request
	// needs, so it would reject everything: "0 disables" only holds for rps.
	if burst < 1 {
		burst = 1
	}

	visitors := make(map[string]*rateLimiterVisitor)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(rateLimiterCleanupInterval)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > rateLimiterStaleAfter {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &rateLimiterVisitor{limiter: rate.NewLimiter(rps, burst)}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		limiter := v.limiter
		mu.Unlock()

		if !limiter.Allow() {
			// The rate is configured as a positive integer, so a token always
			// refills within a second and Retry-After, which is expressed in
			// whole seconds, is always 1.
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, gin.H{"code": http.StatusTooManyRequests, "message": "too many requests", "data": nil})
			c.Abort()
			return
		}

		c.Next()
	}
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Otp      string `form:"otp" json:"otp,omitempty"`
}

var jwtMiddleware *jwt.GinJWTMiddleware
var identityKey = "id"

const jwtRealm = "nethserver"

const (
	// AgentTokenTimeout is the token lifetime of cluster, node and
	// module agents: they have no session and no refresh token
	AgentTokenTimeout = 14 * 24 * time.Hour
	// JoinTokenTimeout is the lifetime of a join code token
	JoinTokenTimeout = time.Hour
	// JoinScope marks a token that can only add a node to the cluster
	JoinScope = "join"
	// ctxRefreshToken passes the session refresh token to LoginResponse
	ctxRefreshToken = "session_refresh_token"
)

// actions a join token can run on the cluster agent
var joinActions = []string{"add-node", "list-actions"}

// sessionTokenStore replaces the gin-jwt v3 refresh token store.
// Refresh tokens are issued by the sessions package, with our own
// /api/refresh endpoint: the random tokens that gin-jwt generates on
// every login are discarded. Delete is called by gin-jwt LogoutHandler
// with the refresh token sent by the client: it ends the session.
type sessionTokenStore struct{}

func (sessionTokenStore) Set(ctx context.Context, token string, userData any, expiry time.Time) error {
	return nil
}

func (sessionTokenStore) Get(ctx context.Context, token string) (any, error) {
	return nil, core.ErrRefreshTokenNotFound
}

func (sessionTokenStore) Delete(ctx context.Context, token string) error {
	return sessions.RevokeRefresh(token)
}

func (sessionTokenStore) Cleanup(ctx context.Context) (int, error) {
	return 0, nil
}

func (sessionTokenStore) Count(ctx context.Context) (int, error) {
	return 0, nil
}

// TokenErrorMessage returns the message of a token validation error as
// gin-jwt v2 (golang-jwt v4) worded it, so clients matching on it, like
// the UI check for "Token is expired", keep working.
func TokenErrorMessage(err error) string {
	switch {
	case errors.Is(err, jwt.ErrExpiredToken), errors.Is(err, gojwt.ErrTokenExpired):
		return "Token is expired"
	case errors.Is(err, gojwt.ErrTokenNotValidYet):
		return "Token is not valid yet"
	case errors.Is(err, jwt.ErrInvalidSigningAlgorithm):
		return jwt.ErrInvalidSigningAlgorithm.Error()
	case errors.Is(err, gojwt.ErrSignatureInvalid):
		return gojwt.ErrSignatureInvalid.Error()
	case errors.Is(err, gojwt.ErrTokenMalformed):
		return strings.TrimPrefix(err.Error(), gojwt.ErrTokenMalformed.Error()+": ")
	}
	return err.Error()
}

func InstanceJWT() *jwt.GinJWTMiddleware {
	if jwtMiddleware == nil {
		jwtMiddleware = InitJWT()
	}
	return jwtMiddleware
}

func InitJWT() *jwt.GinJWTMiddleware {
	// define jwt middleware
	authMiddleware, errDefine := jwt.New(&jwt.GinJWTMiddleware{
		Realm:   jwtRealm,
		Key:     []byte(configuration.Config.Secret),
		Timeout: sessions.AccessTimeout,
		TimeoutFunc: func(data any) time.Duration {
			if user, ok := data.(*models.UserAuthorizations); ok {
				if user.Scope == JoinScope {
					return JoinTokenTimeout
				}
				if methods.IsAgentUser(user.Username) {
					return AgentTokenTimeout
				}
			}
			return sessions.AccessTimeout
		},
		IdentityKey:       identityKey,
		RefreshTokenStore: sessionTokenStore{},
		Authenticator: func(c *gin.Context) (any, error) {
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
				return nil, jwt.ErrFailedAuthentication
			}

			otpNeed := methods.Needs2faCheck(username)
			if otpNeed && otpValue == "" {
				return nil, errors.New("missing OTP value")
			}

			otpPassClaim := false
			if otpNeed {
				if !methods.CheckOTP(username, otpValue) {
					err := errors.New("OTP check failed")
					utils.LogError(errors.Wrap(err, "[AUTH] login-fail(OTP) "+generateLogMessage(c, username, otpNeed)))
					return nil, err
				}
				otpPassClaim = true // claim that 2FA is enabled and used
			}

			// check source IP against per-user allowed networks
			allowedNetworks := methods.GetUserNetworks(username)
			if !methods.CheckIPAllowed(c.ClientIP(), allowedNetworks) {
				utils.LogError(errors.New("[AUTH] login-denied (IP not allowed) " + generateLogMessage(c, username, otpNeed)))
				return nil, jwt.ErrFailedAuthentication
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

			// Start a session for administrators: agents use long-lived
			// tokens without refresh
			sessionID := ""
			if !methods.IsAgentUser(username) {
				id, refreshToken, err := sessions.Create(username, otpPassClaim)
				if err != nil {
					utils.LogError(errors.Wrap(err, "[AUTH] session creation failed for user "+username))
					return nil, sessions.ErrUnavailable
				}
				sessionID = id
				c.Set(ctxRefreshToken, refreshToken)
			}

			utils.LogError(errors.New("[AUTH] login-ok " + generateLogMessage(c, username, otpNeed)))

			// return user auth model
			return &models.UserAuthorizations{
				Username:  username,
				OtpPass:   otpPassClaim, // true if OTP passed, false if OTP is not needed
				SessionID: sessionID,
			}, nil

		},
		PayloadFunc: func(data any) gojwt.MapClaims {
			// read current user
			if user, ok := data.(*models.UserAuthorizations); ok {
				// create claims map
				claims := gojwt.MapClaims{
					identityKey: user.Username,
					"role":      "",
					"actions":   []string{},
					"2fa":       user.OtpPass, // OTP value check result: true=success, false=not needed
				}
				if user.SessionID != "" {
					claims["sid"] = user.SessionID
				}
				if user.Scope != "" {
					claims["scope"] = user.Scope
				}
				return claims
			}

			// return claims map
			return gojwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) any {
			// gin-jwt v3 LogoutHandler runs this handler if the request
			// carries a valid token. Logout needs no identity: skip the
			// Redis lookups, as gin-jwt v2 did.
			if c.FullPath() == "/api/logout" {
				return nil
			}

			// handle identity and extract claims
			claims := jwt.ExtractClaims(c)

			// handle authorizations
			userAuthorizations, err := methods.RedisAuthorization(claims[identityKey].(string), c)
			if err != nil {
				utils.LogError(errors.Wrap(err, "[AUTH] error retrieving user authorizations"))
			}

			// create user object
			user := &models.UserAuthorizations{
				Username:        claims[identityKey].(string),
				Role:            userAuthorizations.Role,
				Actions:         userAuthorizations.Actions,
				AllowedNetworks: methods.GetUserNetworks(claims[identityKey].(string)),
			}
			user.SessionID, _ = claims["sid"].(string)
			user.Scope, _ = claims["scope"].(string)

			// return user
			return user
		},
		Authorizer: func(c *gin.Context, data any) bool {
			// extract data payload and check authorizations
			if v, ok := data.(*models.UserAuthorizations); ok {
				// enforce IP allowlist on every request
				if !methods.CheckIPAllowed(c.ClientIP(), v.AllowedNetworks) {
					utils.LogError(errors.New("[AUTH] request denied (IP not allowed) for user " + v.Username))
					return false
				}

				// a join token can only add a node to the cluster
				if v.Scope == JoinScope {
					return joinTokenAllowed(c, v.Actions)
				}

				// bypass auth for GET requests: // TODO
				if c.Request.Method == "GET" {
					return true
				}

				// any user can end their other sessions
				if c.Request.Method == "POST" && c.FullPath() == "/api/logout-others" {
					return true
				}

				// a join code can be created by users allowed to add nodes
				if c.Request.Method == "POST" && c.FullPath() == "/api/cluster/join-token" {
					return actionAllowed(v.Actions, "add-node")
				}

				// bypass for 2FA apis
				if c.Request.Method == "POST" && c.Request.RequestURI == "/api/2FA" {
					return true
				}
				if c.Request.Method == "DELETE" && c.Request.RequestURI == "/api/2FA" {
					return true
				}
				authorizedActions := v.Actions

				// extract task obj
				var jsonTask models.Task
				errJson := c.ShouldBindBodyWith(&jsonTask, binding.JSON)
				if errJson != nil {
					utils.LogError(errors.Wrap(errJson, "[AUTH] error unmarshalling task obj"))
					return false
				}

				// check action authorization
				return actionAllowed(authorizedActions, jsonTask.Action)
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
		LoginResponse: func(c *gin.Context, token *core.Token) {
			response := gin.H{"code": 200, "expire": time.Unix(token.ExpiresAt, 0), "token": token.AccessToken}
			if refreshToken, ok := c.Get(ctxRefreshToken); ok {
				response["refresh_token"] = refreshToken
			}
			c.JSON(200, response)
		},
		LogoutResponse: func(c *gin.Context) {
			// The refresh token, if sent, was revoked by
			// sessionTokenStore. End also the session of a valid access
			// token, so clients that send only the access token are
			// logged out too.
			if claims, ok := c.Get("JWT_PAYLOAD"); ok {
				if mapClaims, ok := claims.(gojwt.MapClaims); ok {
					if sid, _ := mapClaims["sid"].(string); sid != "" {
						sessions.Revoke(sid)
					}
				}
			}
			c.JSON(200, gin.H{"code": 200})
		},
		HTTPStatusMessageFunc: func(c *gin.Context, e error) string {
			return TokenErrorMessage(e)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			// keep the gin-jwt v2 header value
			c.Header("WWW-Authenticate", "JWT realm="+jwtRealm)

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

func generateLogMessage(c *gin.Context, username string, otpNeed bool) string {
	ipSource := c.ClientIP()
	return fmt.Sprintf("user=%s source=%s 2fa=%t", username, ipSource, otpNeed)
}

// actionAllowed tells if action matches one of the authorized action
// patterns
func actionAllowed(authorizedActions []string, action string) bool {
	for _, authorizedAction := range authorizedActions {
		if ok, _ := filepath.Match(authorizedAction, action); ok {
			return true
		}
	}
	return false
}

// joinTokenAllowed limits a join token to the requests of the
// join-cluster action of a new node: run add-node (and list-actions, as
// credentials probe) on the cluster, and read the task status
func joinTokenAllowed(c *gin.Context, authorizedActions []string) bool {
	switch {
	case c.Request.Method == "GET" && c.FullPath() == "/api/cluster/task/:task_id/status":
		return true
	case c.Request.Method == "POST" && c.FullPath() == "/api/cluster/tasks":
		var jsonTask models.Task
		if err := c.ShouldBindBodyWith(&jsonTask, binding.JSON); err != nil {
			return false
		}
		for _, joinAction := range joinActions {
			if jsonTask.Action == joinAction {
				return actionAllowed(authorizedActions, jsonTask.Action)
			}
		}
	}
	utils.LogError(errors.New("[AUTH] request denied (join token) " + c.Request.Method + " " + c.Request.URL.Path))
	return false
}

// ClaimsSessionValid tells if a token belongs to a valid session.
// Agent and join tokens have no session.
func ClaimsSessionValid(claims gojwt.MapClaims) bool {
	username, _ := claims[identityKey].(string)
	if methods.IsAgentUser(username) {
		return true
	}
	if scope, _ := claims["scope"].(string); scope == JoinScope {
		return true
	}
	sid, _ := claims["sid"].(string)
	return sessions.Valid(sid)
}

// SessionCheck rejects tokens whose session ended: logged out, revoked
// or expired. It runs after the JWT middleware.
func SessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !ClaimsSessionValid(jwt.ExtractClaims(c)) {
			InstanceJWT().Unauthorized(c, http.StatusUnauthorized, sessions.ErrInvalid.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshHandler godoc
// @Summary Get a new access token
// @Description exchange the refresh token for a new access token and a new refresh token
// @Produce json
// @Success 200 {object} response.StatusOK{code=int}
// @Failure 401 {object} response.StatusUnauthorized{code=int,message=string,data=object}
// @Router /refresh [post]
// @Tags /refresh auth
func RefreshHandler(c *gin.Context) {
	var request refreshRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		InstanceJWT().Unauthorized(c, http.StatusUnauthorized, "missing refresh_token")
		return
	}

	session, refreshToken, err := sessions.Refresh(request.RefreshToken)
	if err != nil {
		if err != sessions.ErrInvalid {
			utils.LogError(errors.Wrap(err, "[AUTH] session refresh failed"))
		}
		InstanceJWT().Unauthorized(c, http.StatusUnauthorized, sessions.ErrInvalid.Error())
		return
	}

	if !methods.CheckIPAllowed(c.ClientIP(), methods.GetUserNetworks(session.Username)) {
		utils.LogError(errors.New("[AUTH] refresh-denied (IP not allowed) " + generateLogMessage(c, session.Username, session.OtpPass)))
		InstanceJWT().Unauthorized(c, http.StatusUnauthorized, sessions.ErrInvalid.Error())
		return
	}

	token, err := InstanceJWT().TokenGenerator(c.Request.Context(), &models.UserAuthorizations{
		Username:  session.Username,
		OtpPass:   session.OtpPass,
		SessionID: session.ID,
	})
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUTH] token creation failed"))
		InstanceJWT().Unauthorized(c, http.StatusInternalServerError, jwt.ErrFailedTokenCreation.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          200,
		"expire":        time.Unix(token.ExpiresAt, 0),
		"token":         token.AccessToken,
		"refresh_token": refreshToken,
	})
}

// LogoutOthersHandler godoc
// @Summary End the other sessions of the current user
// @Description log out the current user from other browsers and devices
// @Produce json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=object}
// @Router /logout-others [post]
// @Tags /logout-others auth
func LogoutOthersHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username, _ := claims[identityKey].(string)
	sid, _ := claims["sid"].(string)
	if sid == "" {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "the token has no session",
			Data:    nil,
		}))
		return
	}

	count, err := sessions.RevokeOthers(username, sid)
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUTH] session revocation failed for user "+username))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	utils.LogError(fmt.Errorf("[AUTH] sessions-revoked user=%s count=%d", username, count))

	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "other sessions ended",
		Data:    count,
	}))
}

// JoinTokenHandler godoc
// @Summary Create the token of a join code
// @Description the token can only add a node to the cluster and expires after one hour
// @Produce json
// @Success 200 {object} response.StatusOK{code=int}
// @Router /cluster/join-token [post]
// @Tags /cluster auth
func JoinTokenHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username, _ := claims[identityKey].(string)
	otpPass, _ := claims["2fa"].(bool)

	token, err := InstanceJWT().TokenGenerator(c.Request.Context(), &models.UserAuthorizations{
		Username: username,
		OtpPass:  otpPass,
		Scope:    JoinScope,
	})
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUTH] join token creation failed"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	audit.Store(models.Audit{
		ID:        0,
		User:      username,
		Action:    "join-token",
		Data:      "",
		Timestamp: time.Now().UTC(),
	})

	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"expire": time.Unix(token.ExpiresAt, 0),
		"token":  token.AccessToken,
	})
}
