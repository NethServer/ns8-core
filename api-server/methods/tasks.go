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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"

	"github.com/NethServer/ns8-scratchpad/api-server/models"
	"github.com/NethServer/ns8-scratchpad/api-server/redis"
)

func GetTask(c *gin.Context) {
	// get tasks param
	taskID := c.Param("task_id")

	// init redis connection
	redisConnection := redis.Instance()

	// get task from redis: HGET tasks/:task_id PAYLOAD
	stringTask, errRedis := redisConnection.HGet("tasks/"+taskID, "PAYLOAD").Result()

	// check error
	if errRedis != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no task found"})
		return
	}

	// convert task to struct
	var task models.Task
	errTask := json.Unmarshal([]byte(stringTask), &task)
	if errTask != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "error converting to struct", "status": errTask.Error()})
		return
	}

	// return saved task
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	// bind json body
	var jsonTask models.Task
	if err := c.BindJSON(&jsonTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "status": err.Error()})
		return
	}

	fmt.Println(jsonTask)

	// create task object
	var task models.Task
	task = jsonTask
	task.ID = guuid.New().String()

	// init redis connection
	redisConnection := redis.Instance()

	// convert json struct to json string
	stringTask, errString := json.Marshal(task)
	if errString != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error converting task struct to json", "status": errString.Error()})
		return
	}

	// set task to redis: HSET tasks/:task_id PAYLOAD <payload>
	errRedis := redisConnection.HSet("tasks/"+task.ID, "PAYLOAD", stringTask).Err()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error saving tasks to redis", "status": errRedis.Error()})
		return
	}

	// publish new task in tasks-channel: PUBLISH tasks-channel tasks/:task_id
	errRedisPub := redisConnection.Publish("tasks-channel", "tasks/"+task.ID).Err()

	// handle redis error
	if errRedisPub != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error publish on tasks channel in redis", "status": errRedisPub.Error()})
		return
	}

	// return status created
	c.JSON(http.StatusCreated, gin.H{"message": "task created successfully", "task_id": task.ID})
}

func DeleteTask(c *gin.Context) {
	// get tasks param
	taskID := c.Param("task_id")

	// init redis connection
	redisConnection := redis.Instance()

	// delete task from redis: DEL tasks/:task_id
	errRedis := redisConnection.Del("tasks/" + taskID).Err()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error deleting task in redis", "status": errRedis.Error()})
		return
	}

	// return status created
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})

}
