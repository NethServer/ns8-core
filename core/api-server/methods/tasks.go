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
	"time"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	guuid "github.com/google/uuid"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/NethServer/ns8-scratchpad/core/api-server/redis"
	"github.com/NethServer/ns8-scratchpad/core/api-server/response"
	"github.com/NethServer/ns8-scratchpad/core/api-server/socket"
)

func getAllTasks(c *gin.Context, entityName string) {
	// define all task object
	var allTasks []gin.H

	// init redis connection
	redisConnection := redis.Instance()

	// inspect redis tasks keys: KEYS "node/*/tasks"
	keysArray, errRedisKeys := redisConnection.Keys(ctx, entityName+"/*/tasks").Result()

	// handle redis error
	if errRedisKeys != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting all tasks from redis queue",
			Data:    errRedisKeys.Error(),
		}))
		return
	}

	// loop keys array
	for _, k := range keysArray {
		// inspect redis tasks queue: LRANGE <queue> 0 -1
		tasksArray, errRedisRange := redisConnection.LRange(ctx, k, 0, -1).Result()

		// handle redis error
		if errRedisRange != nil {
			c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
				Code:    400,
				Message: "error getting tasks from redis queue",
				Data:    errRedisRange.Error(),
			}))
			return
		}

		// define tasks array
		var tasks []models.Task

		// loop tasks array
		for _, t := range tasksArray {
			// convert to go struct
			var task models.Task
			errJson := json.Unmarshal([]byte(t), &task)
			if errJson != nil {
				c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
					Code:    400,
					Message: "error converting json task to struct",
					Data:    errJson.Error(),
				}))
				return
			}

			// append to slice
			tasks = append(tasks, task)
		}

		allTasks = append(allTasks, gin.H{"tasks": tasks, "queue": k})
	}

	// return all tasks
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    allTasks,
	}))
}

func getTasks(c *gin.Context, queueName string) {
	// define tasks array
	var tasks []models.Task

	// init redis connection
	redisConnection := redis.Instance()

	// inspect redis tasks queue: LRANGE <queue> 0 -1
	tasksArray, errRedis := redisConnection.LRange(ctx, queueName, 0, -1).Result()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting tasks from redis queue",
			Data:    errRedis.Error(),
		}))
		return
	}

	// loop tasks array
	for _, t := range tasksArray {
		// convert to go struct
		var task models.Task
		errJson := json.Unmarshal([]byte(t), &task)
		if errJson != nil {
			c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
				Code:    400,
				Message: "error converting json task to struct",
				Data:    errJson.Error(),
			}))
			return
		}

		// append to slice
		tasks = append(tasks, task)
	}

	// return tasks
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"tasks": tasks, "queue": queueName},
	}))
}

func getTaskFile(c *gin.Context, filePath string, file string) {
	// init redis connection
	redisConnection := redis.Instance()

	// switch file type
	switch file {
	case "error", "output", "exit_code":
		// read redis attribute
		r, errRedis := redisConnection.Get(ctx, filePath).Result()

		// handle redis error
		if errRedis != nil {
			c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
				Code:    400,
				Message: "error getting file from redis",
				Data:    errRedis.Error(),
			}))
			return
		}

		// return file response
		c.JSON(http.StatusOK, structs.Map(response.StatusOK{
			Code:    200,
			Message: "success",
			Data:    gin.H{"content": r, "file": filePath},
		}))
		return

	default:
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "invalid file type. must be output, error or exit_code",
			Data:    nil,
		}))
		return
	}
}

func createTask(c *gin.Context, queueName string) {
	// bind json body
	var jsonTask models.Task
	if err := c.ShouldBindBodyWith(&jsonTask, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "request fields malformed",
			Data:    err.Error(),
		}))
		return
	}

	// extract user info
	info := jwt.ExtractClaims(c)

	// create task object
	var task models.Task
	task.ID = guuid.New().String()
	task.Action = jsonTask.Action
	task.Data = jsonTask.Data
	task.Queue = queueName
	task.User = info["id"].(string)
	task.Timestamp = time.Now()

	// init redis connection
	redisConnection := redis.Instance()

	// convert json struct to json string
	stringTask, errString := json.Marshal(task)
	if errString != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error converting task struct to json",
			Data:    errString.Error(),
		}))
		return
	}

	// set task to redis: LPUSH <queue> <payload>
	errRedis := redisConnection.LPush(ctx, queueName, string(stringTask)).Err()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error queuing tasks to redis",
			Data:    errRedis.Error(),
		}))
		return
	}

	// return status created
	c.JSON(http.StatusCreated, structs.Map(response.StatusOK{
		Code:    201,
		Message: "task queued successfully",
		Data:    task,
	}))
}

func ListenTaskEvents() {
	// init redis connection
	redisConnection := redis.Instance()

	// get socket instance
	socketConnection := socket.Instance()

	// subscribe to progress channel and listen to new messages: PSUBSCRIBE progress/task/*
	progress := redisConnection.PSubscribe(ctx, "progress/task/*")

	// get the channel to use
	channel := progress.Channel()

	// start routine to listen to new messages
	go func() {
		// iterate any messages sent on the channel
		for msg := range channel {
			t := &models.TaskProgress{}

			// unmarshal the data into the task progress
			err := t.UnmarshalBinary([]byte(msg.Payload))
			if err != nil {
				panic(err)
			}

			// identify all sessions associated with task id
			if clientSession, ok := socket.Connections["/ws/tasks/progress"]; ok {
				// Broadcast to all sessions
				socketConnection.BroadcastMultiple([]byte(msg.Payload), clientSession)
			}
		}
	}()
}

// GetClusterTasks godoc
// @Summary Get the list of current cluster tasks
// @Description get cluster tasks
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/cluster [get]
// @Tags /tasks cluster
func GetClusterTasks(c *gin.Context) {
	// define queue name
	queueName := "cluster/tasks"

	// get tasks
	getTasks(c, queueName)
}

// GetClusterTaskFiles godoc
// @Summary Get the output, error or exit code of a cluster task
// @Description get task statuses (output, error, exit_code)
// @Produce  json
// @Param task_id path string true "Task ID"
// @Param file path string true "Must be: `output`, `error`, `exit_code`"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/cluster/{task_id}/{file} [get]
// @Tags /tasks cluster
func GetClusterTaskFiles(c *gin.Context) {
	// get task id
	taskID := c.Param("task_id")

	// get file type
	file := c.Param("file")

	// define queue name
	filePath := "cluster/task/" + taskID + "/" + file

	// get result of file
	getTaskFile(c, filePath, file)
}

// GetNodeTasks godoc
// @Summary Get the list of current node tasks
// @Description get node tasks
// @Produce  json
// @Param node_id path string true "Node ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/node/{node_id} [get]
// @Tags /tasks node
func GetNodeTasks(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// define queue name
	queueName := "node/" + nodeID + "/tasks"

	// get tasks
	getTasks(c, queueName)
}

// GetNodeTaskFiles godoc
// @Summary Get the output, error or exit code of a node task
// @Description get task statuses (output, error, exit_code)
// @Produce  json
// @Param node_id path string true "Node ID"
// @Param task_id path string true "Task ID"
// @Param file path string true "Must be: `output`, `error`, `exit_code`"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/node/{node_id}/{task_id}/{file} [get]
// @Tags /tasks node
func GetNodeTaskFiles(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// get task id
	taskID := c.Param("task_id")

	// get file type
	file := c.Param("file")

	// define queue name
	filePath := "node/" + nodeID + "/task/" + taskID + "/" + file

	// get result of file
	getTaskFile(c, filePath, file)
}

// GetModuleTasks godoc
// @Summary Get the list of current module tasks
// @Description get module tasks
// @Produce  json
// @Param module_id path string true "Module ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/module/{module_id} [get]
// @Tags /tasks module
func GetModuleTasks(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// define queue name
	queueName := "module/" + moduleID + "/tasks"

	// get tasks
	getTasks(c, queueName)
}

// GetModuleTaskFiles godoc
// @Summary Get the output, error or exit code of a module task
// @Description get task statuses (output, error, exit_code)
// @Produce  json
// @Param module_id path string true "Module ID"
// @Param task_id path string true "Task ID"
// @Param file path string true "Must be: `output`, `error`, `exit_code`"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/module/{module_id}/{task_id}/{file} [get]
// @Tags /tasks module
func GetModuleTaskFiles(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// get task id
	taskID := c.Param("task_id")

	// get file type
	file := c.Param("file")

	// define queue name
	filePath := "module/" + moduleID + "/task/" + taskID + "/" + file

	// get result of file
	getTaskFile(c, filePath, file)
}

// CreateClusterTask godoc
// @Summary Create and queue a cluster task
// @Description create a cluster task
// @Produce  json
// @Success 200 {object} response.StatusCreated{code=int,message=string,data=models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/cluster [post]
// @Tags /tasks cluster
func CreateClusterTask(c *gin.Context) {
	// define queue name
	queueName := "cluster/tasks"

	// create task
	createTask(c, queueName)
}

// CreateNodeTask godoc
// @Summary Create and queue a node task
// @Description create a node task
// @Produce  json
// @Param node_id path string true "Node ID"
// @Success 200 {object} response.StatusCreated{code=int,message=string,data=models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/node/{node_id} [post]
// @Tags /tasks node
func CreateNodeTask(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// define queue name
	queueName := "node/" + nodeID + "/tasks"

	// create task
	createTask(c, queueName)
}

// CreateModuleTask godoc
// @Summary Create and queue a module task
// @Description create a module task
// @Produce  json
// @Param module_id path string true "Module ID"
// @Success 200 {object} response.StatusCreated{code=int,message=string,data=models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /tasks/module/{module_id} [post]
// @Tags /tasks module
func CreateModuleTask(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// define queue name
	queueName := "module/" + moduleID + "/tasks"

	// create task
	createTask(c, queueName)
}
