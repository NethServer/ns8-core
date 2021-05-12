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
	"github.com/gin-gonic/gin"

	docs "github.com/NethServer/ns8-scratchpad/core/api-server/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/NethServer/ns8-scratchpad/core/api-server/configuration"
	"github.com/NethServer/ns8-scratchpad/core/api-server/methods"
	"github.com/NethServer/ns8-scratchpad/core/api-server/middleware"
	"github.com/NethServer/ns8-scratchpad/core/api-server/response"
	"github.com/NethServer/ns8-scratchpad/core/api-server/socket"
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

	// init websocket
	socket := socket.Instance()

	// init routers
	router := gin.Default()

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
		tasks := api.Group("/tasks")
		{
			// cluster
			tasks.GET("/cluster", methods.GetClusterTasks)
			tasks.GET("/cluster/:task_id/status", methods.GetClusterTaskFiles)
			tasks.POST("/cluster", methods.CreateClusterTask)

			// node
			tasks.GET("/node/:node_id", methods.GetNodeTasks)
			tasks.GET("/node/:node_id/:task_id/status", methods.GetNodeTaskFiles)
			tasks.POST("/node/:node_id", methods.CreateNodeTask)

			// module
			tasks.GET("/module/:module_id", methods.GetModuleTasks)
			tasks.GET("/module/:module_id/:task_id/status", methods.GetModuleTaskFiles)
			tasks.POST("/module/:module_id", methods.CreateModuleTask)
		}
	}

	// define websocket group
	ws := router.Group("/ws")

	// task ws APIS
	task := ws.Group("/tasks")
	{
		task.GET("/progress", func(c *gin.Context) {
			socket.HandleRequest(c.Writer, c.Request)
		})
	}

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
	docs.SwaggerInfo.Host = configuration.Config.ListenAddress
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run server
	router.Run(configuration.Config.ListenAddress)
}
