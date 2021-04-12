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

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/NethServer/ns8-scratchpad/api-server/configuration"
	"github.com/NethServer/ns8-scratchpad/api-server/methods"
)

func main() {
	// init configuration
	configuration.Init()

	// disable log to stdout when running in release mode
	if gin.Mode() == gin.ReleaseMode {
		gin.DefaultWriter = ioutil.Discard
	}

	// init routers
	router := gin.Default()

	// add default compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// cors configuration only in debug mode GIN_MODE=debug (default)
	if gin.Mode() == gin.DebugMode {
		corsConf := cors.DefaultConfig()
		corsConf.AllowHeaders = []string{"Authorization", "Content-Type"}
		corsConf.AllowAllOrigins = true
		router.Use(cors.New(corsConf))
	}

	// define api group
	api := router.Group("/api")

	// cluster APIs
	cluster := api.Group("/cluster")
	{
		cluster.GET("/tasks", methods.GetClusterTasks)
		cluster.POST("/tasks", methods.CreateClusterTask)
	}
	// node APIs
	node := api.Group("/node")
	{
		node.GET("/tasks", methods.GetAllNodeTasks)
		node.GET("/:node_id/tasks", methods.GetNodeTasks)
		node.POST("/:node_id/tasks", methods.CreateNodeTask)
	}
	// module APIs
	module := api.Group("/module")
	{
		module.GET("/tasks", methods.GetAllModuleTasks)
		module.GET("/:module_id/tasks", methods.GetModuleTasks)
		module.POST("/:module_id/tasks", methods.CreateModuleTask)
	}

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})

	// run server
	router.Run(configuration.Config.ListenAddress)
}
