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
	"strconv"
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

		// close redis connection
		redisConnection.Close()

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

	// close redis connection
	redisConnection.Close()

	// return tasks
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"tasks": tasks, "queue": queueName},
	}))
}

func getTaskFile(c *gin.Context, filePath string) {
	// init redis connection
	redisConnection := redis.Instance()

	// read redis statuses
	errorC, errRedisError := redisConnection.Get(ctx, filePath+"/error").Result()
	outputC, errRedisOutput := redisConnection.Get(ctx, filePath+"/output").Result()
	exitCodeC, errRedisExitCode := redisConnection.Get(ctx, filePath+"/exit_code").Result()

	// handle redis error for error status
	if errRedisError != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting error status from redis",
			Data:    errRedisError.Error(),
		}))
		return
	}

	// handle redis error for output status
	if errRedisOutput != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting output status from redis",
			Data:    errRedisOutput.Error(),
		}))
		return
	}

	// handle redis error for exit code status
	if errRedisExitCode != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting exit code status from redis",
			Data:    errRedisExitCode.Error(),
		}))
		return
	}

	// close redis connection
	redisConnection.Close()

	// decode the exit_status and handle the error
	exitCodeInt, errValueExitCode := strconv.Atoi(exitCodeC)
	if errValueExitCode != nil {
		c.JSON(http.StatusInternalServerError, structs.Map(response.StatusInternalServerError{
			Code:    500,
			Message: "error in exit code value",
			Data:    errValueExitCode.Error(),
		}))
		return
	}

	var outputData interface{}
	// We expect a JSON-encoded output. If the unmarshal fails, return the output as-is.
	if errOutputDecode := json.Unmarshal([]byte(outputC), &outputData); errOutputDecode != nil {
		outputData = outputC
	}

	// return file response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"error": errorC, "output": outputData, "exit_code": exitCodeInt, "file": filePath},
	}))
	return
}

func getTaskContext(c *gin.Context, filePath string) {
	// init redis connection
	redisConnection := redis.Instance()

	// read redis statuses
	contextC, errRedisContext := redisConnection.Get(ctx, filePath+"/context").Result()

	// handle redis error for context status
	if errRedisContext != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting context status from redis",
			Data:    errRedisContext.Error(),
		}))
		return
	}

	// close redis connection
	redisConnection.Close()

	// return file response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"context": contextC, "file": filePath},
	}))
	return
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
	task.Parent = ""

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

	// close redis connection
	redisConnection.Close()

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
	progress := redisConnection.PSubscribe(ctx, "progress/*")

	// get the channel to use
	channel := progress.Channel()

	// start routine to listen to new messages
	go func() {
		// iterate any messages sent on the channel
		for msg := range channel {
			// create event object
			event := &models.Event{}
			event.Name = msg.Channel
			event.Timestamp = time.Now()
			event.Type = "task"

			if err := json.Unmarshal([]byte(msg.Payload), &event.Payload); err != nil {
				panic(err)
			}

			// marshal model to json string
			taskJSON, errJSON := json.Marshal(event)
			if errJSON != nil {
				panic(errJSON)
			}

			// identify all sessions associated with task id
			if clientSession, ok := socket.Connections["/ws"]; ok {
				// Broadcast to all sessions
				socketConnection.BroadcastMultiple(taskJSON, clientSession)
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
// @Router /cluster/tasks [get]
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
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /cluster/task/{task_id}/status [get]
// @Tags /tasks cluster
func GetClusterTaskFiles(c *gin.Context) {
	// get task id
	taskID := c.Param("task_id")

	// define queue name
	filePath := "cluster/task/" + taskID

	// get result of file
	getTaskFile(c, filePath)
}

// GetClusterTaskContext godoc
// @Summary Get the context of a cluster task
// @Description get task statuses (context)
// @Produce  json
// @Param task_id path string true "Task ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /cluster/task/{task_id}/context [get]
// @Tags /tasks cluster
func GetClusterTaskContext(c *gin.Context) {
	// get task id
	taskID := c.Param("task_id")

	// define queue name
	filePath := "cluster/task/" + taskID

	// get result of file
	getTaskContext(c, filePath)
}

// GetNodeTasks godoc
// @Summary Get the list of current node tasks
// @Description get node tasks
// @Produce  json
// @Param node_id path string true "Node ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /node/{node_id}/tasks [get]
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
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /node/{node_id}/task/{task_id}/status [get]
// @Tags /tasks node
func GetNodeTaskFiles(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// get task id
	taskID := c.Param("task_id")

	// define queue name
	filePath := "node/" + nodeID + "/task/" + taskID

	// get result of file
	getTaskFile(c, filePath)
}

// GetNodeTaskContext godoc
// @Summary Get the context of a node task
// @Description get task statuses (context)
// @Produce  json
// @Param node_id path string true "Node ID"
// @Param task_id path string true "Task ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /node/{node_id}/task/{task_id}/context [get]
// @Tags /tasks node
func GetNodeTaskContext(c *gin.Context) {
	// get param
	nodeID := c.Param("node_id")

	// get task id
	taskID := c.Param("task_id")

	// define queue name
	filePath := "node/" + nodeID + "/task/" + taskID

	// get result of file
	getTaskContext(c, filePath)
}

// GetModuleTasks godoc
// @Summary Get the list of current module tasks
// @Description get module tasks
// @Produce  json
// @Param module_id path string true "Module ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /module/{module_id}/tasks [get]
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
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /module/{module_id}/task/{task_id}/status [get]
// @Tags /tasks module
func GetModuleTaskFiles(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// get task id
	taskID := c.Param("task_id")

	// define queue name
	filePath := "module/" + moduleID + "/task/" + taskID

	// get result of file
	getTaskFile(c, filePath)
}

// GetModuleTaskContext godoc
// @Summary Get the context of a module task
// @Description get task statuses (context)
// @Produce  json
// @Param module_id path string true "Module ID"
// @Param task_id path string true "Task ID"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /module/{module_id}/task/{task_id}/context [get]
// @Tags /tasks module
func GetModuleTaskContext(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// get task id
	taskID := c.Param("task_id")

	// define queue name
	filePath := "module/" + moduleID + "/task/" + taskID

	// get result of file
	getTaskContext(c, filePath)
}

// CreateClusterTask godoc
// @Summary Create and queue a cluster task
// @Description create a cluster task
// @Produce  json
// @Success 200 {object} response.StatusCreated{code=int,message=string,data=models.Task}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /cluster/tasks [post]
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
// @Router /node/{node_id}/tasks [post]
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
// @Router /module/{module_id}/tasks [post]
// @Tags /tasks module
func CreateModuleTask(c *gin.Context) {
	// get param
	moduleID := c.Param("module_id")

	// define queue name
	queueName := "module/" + moduleID + "/tasks"

	// create task
	createTask(c, queueName)
}
