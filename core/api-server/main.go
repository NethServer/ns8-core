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

package main

import (
	"io/ioutil"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/NethServer/ns8-core/core/api-server/audit"
	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/methods"
	"github.com/NethServer/ns8-core/core/api-server/middleware"
	"github.com/NethServer/ns8-core/core/api-server/response"
	"github.com/NethServer/ns8-core/core/api-server/socket"
)

// @title NethServer 8 API
// @version 1.0
// @description NethServer 8 API is used to create tasks across the nodes
// @termsOfService https://nethserver.org/terms/

// @contact.name NethServer Developer Team
// @contact.url https://nethserver.org/support

// @license.name GNU GENERAL PUBLIC LICENSE

// @host localhost:8080
// @schemes http
// @BasePath /api

func main() {
	// init configuration
	configuration.Init()

	// disable log to stdout when running in release mode
	if gin.Mode() == gin.ReleaseMode {
		gin.DefaultWriter = ioutil.Discard
	}

	// init audit file
	audit.Init()

	// init websocket
	socketConnection := socket.Instance()

	// init routers
	//router := gin.Default()
	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/ws"),
		gin.Recovery(),
	)

	// add default compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// cors configuration only in debug mode GIN_MODE=debug (default)
	if gin.Mode() == gin.DebugMode {
		// gin gonic cors conf
		corsConf := cors.DefaultConfig()
		corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Accept"}
		corsConf.AllowAllOrigins = true
		router.Use(cors.New(corsConf))
	}

	// define static file endpoint
	router.Use(static.Serve("/", static.LocalFile(configuration.Config.StaticPath, false)))

	// define api group
	api := router.Group("/api")

	// define login and logout endpoint
	api.POST("/login", middleware.InstanceJWT().LoginHandler)
	api.POST("/logout", middleware.InstanceJWT().LogoutHandler)

	// define refresh endpoint
	api.GET("/refresh_token", middleware.InstanceJWT().RefreshHandler)

	// define JWT middleware
	api.Use(middleware.InstanceJWT().MiddlewareFunc())
	{
		// tasks APIs
		// cluster
		api.GET("/cluster/tasks", methods.GetClusterTasks)
		api.GET("/cluster/task/:task_id/status", methods.GetClusterTaskFiles)
		api.GET("/cluster/task/:task_id/context", methods.GetClusterTaskContext)
		api.POST("/cluster/tasks", methods.CreateClusterTask)

		// node
		api.GET("/nodes", methods.GetNodes)
		api.GET("/node/:node_id/tasks", methods.GetNodeTasks)
		api.GET("/node/:node_id/task/:task_id/status", methods.GetNodeTaskFiles)
		api.GET("/node/:node_id/task/:task_id/context", methods.GetNodeTaskContext)
		api.POST("/node/:node_id/tasks", methods.CreateNodeTask)

		// module
		api.GET("/modules", methods.GetModules)
		api.GET("/module/:module_id/tasks", methods.GetModuleTasks)
		api.GET("/module/:module_id/task/:task_id/status", methods.GetModuleTaskFiles)
		api.GET("/module/:module_id/task/:task_id/context", methods.GetModuleTaskContext)
		api.POST("/module/:module_id/tasks", methods.CreateModuleTask)

		// audit APIs
		api.GET("/audit", methods.GetAudits)
		api.GET("/audit/users", methods.GetAuditsUsers)
		api.GET("/audit/actions", methods.GetAuditsActions)
	}

	// define websocket endpoint
	ws := router.Group("/ws")
	ws.GET("", func(c *gin.Context) {
		if socket.ValidateAuth(c.Query("jwt")) {
			socketConnection.HandleRequest(c.Writer, c.Request)
		}
	})

	// start events
	methods.ListenTaskEvents()

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, structs.Map(response.StatusNotFound{
			Code:    404,
			Message: "API not found",
			Data:    nil,
		}))
	})

	// handle API docs route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run server
	router.Run(configuration.Config.ListenAddress)
}
