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
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"

	"github.com/NethServer/ns8-scratchpad/api-server/models"
	"github.com/NethServer/ns8-scratchpad/api-server/redis"
)

func getTasks(c *gin.Context, queueName string) {
	// define tasks array
	var tasks []models.Task

	// init redis connection
	redisConnection := redis.Instance()

	// inspect redis tasks queue: LRANGE <queue> 0 -1
	tasksArray, errRedis := redisConnection.LRange(queueName, 0, -1).Result()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error getting tasks from redis queue", "status": errRedis.Error()})
		return
	}

	// loop tasks array
	for _, t := range tasksArray {
		// convert to go struct
		var task models.Task
		errJson := json.Unmarshal([]byte(t), &task)
		if errJson != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error converting json task to struct", "status": errJson.Error()})
			return
		}

		// append to slice
		tasks = append(tasks, task)
	}

	// return status created
	c.JSON(http.StatusOK, gin.H{"tasks": tasks, "queue": queueName})
}

func createTask(c *gin.Context, queueName string) {
	// bind json body
	var jsonTask models.Task
	if err := c.BindJSON(&jsonTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "status": err.Error()})
		return
	}

	// create task object
	var task models.Task
	task.ID = guuid.New().String()
	task.Action = jsonTask.Action
	task.Data = jsonTask.Data

	// init redis connection
	redisConnection := redis.Instance()

	// convert json struct to json string
	stringTask, errString := json.Marshal(task)
	if errString != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error converting task struct to json", "status": errString.Error()})
		return
	}

	// set task to redis: LPUSH <queue> <payload>
	errRedis := redisConnection.LPush(queueName, string(stringTask)).Err()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error queuing tasks to redis", "status": errRedis.Error()})
		return
	}

	// return status created
	c.JSON(http.StatusCreated, gin.H{"message": "task queued successfully", "task": task, "queue": queueName})
}

func GetClusterTasks(c *gin.Context) {
	// define queue name
	queueName := "cluster/tasks"

	// get tasks
	getTasks(c, queueName)
}

func GetNodeTasks(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// define queue name
	queueName := "node/" + nodeID + "/tasks"

	// get tasks
	getTasks(c, queueName)
}

func GetModuleTasks(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// define queue name
	queueName := "module/" + moduleID + "/tasks"

	// get tasks
	getTasks(c, queueName)
}

func CreateClusterTask(c *gin.Context) {
	// define queue name
	queueName := "cluster/tasks"

	// create task
	createTask(c, queueName)
}

func CreateNodeTask(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// define queue name
	queueName := "node/" + nodeID + "/tasks"

	// create task
	createTask(c, queueName)
}

func CreateModuleTask(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// define queue name
	queueName := "module/" + moduleID + "/tasks"

	// create task
	createTask(c, queueName)
}

// func DeleteTask(c *gin.Context) {
// 	// get tasks param
// 	taskID := c.Param("task_id")

// 	// init redis connection
// 	redisConnection := redis.Instance()

// 	// delete task from redis: DEL tasks/:task_id
// 	errRedis := redisConnection.Del("tasks/" + taskID).Err()

// 	// handle redis error
// 	if errRedis != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "error deleting task in redis", "status": errRedis.Error()})
// 		return
// 	}

// 	// return status created
// 	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})

// }
