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
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"

	"github.com/NethServer/ns8-scratchpad/core/api-server/audit"
	"github.com/NethServer/ns8-scratchpad/core/api-server/configuration"
	"github.com/NethServer/ns8-scratchpad/core/api-server/models"
	"github.com/NethServer/ns8-scratchpad/core/api-server/response"
)

// GetAudits godoc
// @Summary Get audit logs for user, action and timestamp
// @Description get audit logs (user, action, data, timestamp)
// @Produce  json
// @Param user query string false "user search by user"
// @Param action query string false "action search by action"
// @Param from query string false "filter result search by from date ISO8601"
// @Param to query string false "filter result search by to date ISO8601"
// @Success 200 {object} response.StatusOK{code=int,message=string,data=string}
// @Header 200 {string} Authorization "Bearer <valid.JWT.token>"
// @Failure 400 {object} response.StatusBadRequest{code=int,message=string,data=object}
// @Router /audit [get]
// @Tags /audit audit
func GetAudits(c *gin.Context) {
	// get query params
	user := c.Query("user")
	action := c.Query("action")
	from := c.Query("from")
	to := c.Query("to")

	// define query
	query := "SELECT * FROM audit WHERE true"

	// check params
	if len(user) > 0 {
		query += " AND user = ?"
	}

	if len(action) > 0 {
		query += " AND action = ?"
	}

	if len(from) > 0 {
		query += " AND timestamp >= ?"
	}

	if len(to) > 0 {
		query += " AND timestamp <= ?"
	}

	// execute query
	results := audit.Query(query, user, action, from, to)

	// store audit
	claims := jwt.ExtractClaims(c)
	auditData := models.Audit{
		User:      claims["id"].(string),
		Action:    "list-audits",
		Data:      "",
		Timestamp: time.Now(),
	}
	audit.Store(auditData)

	// return results
	if len(results) > 0 {
		c.JSON(http.StatusOK, structs.Map(response.StatusOK{
			Code:    200,
			Message: "success",
			Data:    gin.H{"audits": results},
		}))
	} else {
		if len(configuration.Config.AuditFile) == 0 {
			c.JSON(http.StatusBadRequest, structs.Map(response.StatusBadRequest{
				Code:    400,
				Message: "audit is disabled. AUDIT_FILE is not set in the environment",
				Data:    gin.H{"audits": nil},
			}))
		} else {
			c.JSON(http.StatusNotFound, structs.Map(response.StatusNotFound{
				Code:    404,
				Message: "no audits found",
				Data:    gin.H{"audits": nil},
			}))
		}
	}

}
