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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	guuid "github.com/google/uuid"
	"github.com/mpvl/unique"
	"github.com/nqd/flat"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/NethServer/ns8-core/core/api-server/audit"
	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/redis"
	"github.com/NethServer/ns8-core/core/api-server/response"
)

func getList(c *gin.Context, queueName string) {
	// init redis connection
	redisConnection := redis.Instance()

	// inspect redis queue: KEYS <queue>
	listArray, errRedis := redisConnection.Keys(ctx, queueName).Result()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting list from redis queue",
			Data:    errRedis.Error(),
		}))
		return
	}

	// define result list
	var list []string

	// loop list array
	for _, l := range listArray {
		// get parts
		parts := strings.Split(l, "/")
		object := parts[1]

		// append object to list
		list = append(list, object)
	}

	// unique the list
	sort.Strings(list)
	unique.Strings(&list)

	// close redis connection
	redisConnection.Close()

	// store to audit
	claims := jwt.ExtractClaims(c)
	parts := strings.Split(queueName, "/")

	auditData := models.Audit{
		ID:        0,
		User:      claims["id"].(string),
		Action:    "list-" + parts[0],
		Data:      queueName,
		Timestamp: time.Now().UTC(),
	}
	audit.Store(auditData)

	// return tasks
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"list": list, "queue": queueName},
	}))
}

func getTasks(c *gin.Context, queueName string) {
	// init redis connection
	redisConnection := redis.Instance()

	// inspect redis tasks queue: KEYS <queue>
	tasksArray, errRedis := redisConnection.Keys(ctx, queueName).Result()

	// handle redis error
	if errRedis != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "error getting tasks from redis queue",
			Data:    errRedis.Error(),
		}))
		return
	}

	// define map
	var taskMap map[string]map[string]interface{}

	// loop tasks array
	for _, t := range tasksArray {
		// extract task id
		parts := strings.Split(t, "/")
		file := parts[len(parts)-1]
		task := strings.Join(parts[:len(parts)-1], "/")

		// check if map exists
		if taskMap == nil {
			taskMap = make(map[string]map[string]interface{})
		}
		// check if chil of map exists
		if taskMap[task] == nil {
			taskMap[task] = make(map[string]interface{})
		}

		// get content of key: GET <task_id>
		taskContent, errRedisContent := redisConnection.Get(ctx, t).Result()

		// handle error
		if errRedisContent != nil {
			c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
				Code:    400,
				Message: "error getting tasks status from redis key",
				Data:    errRedis.Error(),
			}))
			return
		}

		var taskContentData interface{}
		// We expect a JSON-encoded output. If the unmarshal fails, return the output as-is.
		if errOutputDecode := json.Unmarshal([]byte(taskContent), &taskContentData); errOutputDecode != nil {
			taskContentData = taskContent
		}

		// assign map
		taskMap[task][file] = taskContentData
	}

	// close redis connection
	redisConnection.Close()

	// store to audit
	claims := jwt.ExtractClaims(c)
	auditData := models.Audit{
		ID:        0,
		User:      claims["id"].(string),
		Action:    "list-task",
		Data:      queueName,
		Timestamp: time.Now().UTC(),
	}
	audit.Store(auditData)

	// return tasks
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"tasks": taskMap, "queue": queueName},
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

	// store to audit
	claims := jwt.ExtractClaims(c)
	auditData := models.Audit{
		ID:        0,
		User:      claims["id"].(string),
		Action:    "status-task",
		Data:      filePath,
		Timestamp: time.Now().UTC(),
	}
	audit.Store(auditData)

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

	var contextData interface{}
	// We expect a JSON-encoded context. If the unmarshal fails, return the context as-is.
	if errOutputDecode := json.Unmarshal([]byte(contextC), &contextData); errOutputDecode != nil {
		contextData = contextC
	}

	// store to audit
	claims := jwt.ExtractClaims(c)
	auditData := models.Audit{
		ID:        0,
		User:      claims["id"].(string),
		Action:    "context-task",
		Data:      filePath,
		Timestamp: time.Now().UTC(),
	}
	audit.Store(auditData)

	// return file response
	c.JSON(http.StatusOK, structs.Map(response.StatusOK{
		Code:    200,
		Message: "success",
		Data:    gin.H{"context": contextData, "file": filePath},
	}))
	return
}

func createTask(c *gin.Context, queueName string) {
	// bind json body
	var jsonTask models.TaskJSON
	if err := c.ShouldBindBodyWith(&jsonTask, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
			Code:    400,
			Message: "request fields malformed",
			Data:    err.Error(),
		}))
		return
	}

	// init redis connection
	redisConnection := redis.Instance()

	// check redis role
	redisRole := redisConnection.Do(ctx, "role").String()
	if !strings.Contains(redisRole, "master") {
		c.JSON(http.StatusForbidden, structs.Map(response.StatusForbidden{
			Code:    403,
			Message: "task submission is forbidden in a worker node",
			Data:    redis.GetLeaderHostAddress(ctx, redisConnection),
		}))
		return
	}

	// check the agent is connected and alive
	if err := redis.CheckClientIdle(c, redisConnection, strings.TrimSuffix(queueName, "/tasks"), 8); err != nil {
		c.JSON(http.StatusNotFound, structs.Map(response.StatusNotFound{
			Code:    404,
			Message: "client idle check failed",
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
	task.Extra = jsonTask.Extra
	task.Queue = queueName
	task.User = info["id"].(string)
	task.Timestamp = time.Now().UTC()
	task.Parent = jsonTask.Parent

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

	// convert to map and flat it
	var jsonDyn map[string]interface{}
	json.Unmarshal(stringTask, &jsonDyn)
	in, _ := flat.Flatten(jsonDyn, nil)

	// search for sensitve data, in sensitive list
	for k, _ := range in {
		for _, s := range configuration.Config.SensitiveList {
			if strings.Contains(strings.ToLower(k), strings.ToLower(s)) {
				in[k] = "XXX"
			}
		}
	}

	// unflat the map
	out, _ := flat.Unflatten(in, nil)

	// convert to json string
	jsonOut, _ := json.Marshal(out)

	// store to audit
	auditData := models.Audit{
		ID:        0,
		User:      task.User,
		Action:    "create-task",
		Data:      string(jsonOut),
		Timestamp: time.Now().UTC(),
	}
	audit.Store(auditData)

	// return status created
	c.JSON(http.StatusCreated, structs.Map(response.StatusCreated{
		Code:    201,
		Message: "task queued successfully",
		Data:    task,
	}))
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
	queueName := "task/cluster/*"

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
	filePath := "task/cluster/" + taskID

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
	filePath := "task/cluster/" + taskID

	// get result of file
	getTaskContext(c, filePath)
}

// GetNodes godoc
// @Summary Get the list of nodes
// @Description get nodes
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /nodes [get]
// @Tags /tasks node
func GetNodes(c *gin.Context) {
	// define queue name
	queueName := "node/*"

	// get list
	getList(c, queueName)
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
	queueName := "task/node/" + nodeID + "/*"

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
	filePath := "task/node/" + nodeID + "/" + taskID

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
	filePath := "task/node/" + nodeID + "/" + taskID

	// get result of file
	getTaskContext(c, filePath)
}

// GetModules godoc
// @Summary Get the list of modules
// @Description get modules
// @Produce  json
// @Success 200 {object} response.StatusOK{code=int,message=string,data=[]string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /modules [get]
// @Tags /tasks node
func GetModules(c *gin.Context) {
	// define queue name
	queueName := "module/*"

	// get list
	getList(c, queueName)
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
	queueName := "task/module/" + moduleID + "/*"

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
	filePath := "task/module/" + moduleID + "/" + taskID

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
	filePath := "task/module/" + moduleID + "/" + taskID

	// get result of file
	getTaskContext(c, filePath)
}

// CreateClusterTask godoc
// @Summary Create and queue a cluster task
// @Description create a cluster task
// @Produce  json
// @Param payload body models.TaskJSON true "Task Payload"
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
// @Param payload body models.TaskJSON true "Task Payload"
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
// @Param payload body models.TaskJSON true "Task Payload"
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
