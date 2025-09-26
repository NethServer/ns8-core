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
 * along with NethServer. If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/utils"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type (
	failureModel struct {
		faultySchema bool
		faultyOpen   bool
		faultyConfig bool
	}

	dbUtils struct {
		conn         *sql.DB
		openFail     error
		faultyStatus failureModel
		tableExist   func(err error) bool
		isFaulty     func(f failureModel) (bool, string)
	}
)

var db *dbUtils

func Init() {
	db = &dbUtils{
		faultyStatus: failureModel{},
		tableExist: func(err error) bool {
			return strings.Contains(err.Error(), "already exists")
		},
		isFaulty: func(s failureModel) (bool, string) {
			if s.faultySchema {
				return true, "issues in schema creation"
			}
			if s.faultyOpen { 
				return true, "issues in opening the database file"
			}
			if s.faultyConfig {
				return true, "issues in database configuration"
			}
			
			return false, ""
		},
	}

	if len(configuration.Config.AuditFile) == 0 {
		utils.LogError(errors.Wrap(errors.New("AUDIT_FILE is not set in the environment."), "[AUDIT][INIT]"))
		db.faultyStatus.faultyConfig = true
		return
	}

	createDB()
}

func createDB() {
	db.conn, db.openFail = sql.Open("sqlite3", configuration.Config.AuditFile)
	if db.openFail != nil {
		utils.LogError(errors.Wrap(db.openFail, "[AUDIT][CREATE] error in audit db file creation."))
		db.faultyStatus.faultyOpen = true
		db.conn.Close()
		return
	}

	// define audit schema
	query := `
		CREATE TABLE audit (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user TEXT NOT NULL,
			action TEXT NOT NULL,
			data TEXT NOT NULL,
			timestamp TEXT NOT NULL
		);
	`
	_, errExecute := db.conn.Exec(query)
	if errExecute != nil {
		if db.tableExist(errExecute) {
			return // table already exists, no problem
		}

		utils.LogError(errors.Wrap(errExecute, "[AUDIT][CREATE] error in audit file schema init."))
		db.faultyStatus.faultySchema = true
	}
}

func Store(audit models.Audit) {
	if ok, out := db.isFaulty(db.faultyStatus); ok {
		utils.LogError(errors.Wrap(errors.New("Connection dropped due to " + out), "[AUDIT][STORE]"))
		return
	}

	tx, err := db.conn.Begin()
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error occurred while initiating the transaction, rollback enforced"))
		tx.Rollback()
		return
	}

	stmt, err := tx.Prepare("INSERT INTO audit(id, user, action, data, timestamp) VALUES(null, ?, ?, ?, ?)")
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error occurred while preparing the transaction, rollback enforced"))
		tx.Rollback()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(audit.User, audit.Action, audit.Data, audit.Timestamp.Format(time.RFC3339))
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error occurred when executing the transaction, rollback enforced"))
		tx.Rollback()
		return
	}
	tx.Commit()
}

func QueryArgs(query string, args ...interface{}) []models.Audit {
	if ok, out := db.isFaulty(db.faultyStatus); ok {
		utils.LogError(errors.Wrap(errors.New("Connection dropped due to " + out ), "[AUDIT][QUERY]"))
		return nil
	}

	var results []models.Audit

	cleanArgs := make([]interface{}, 0, len(args))
	for _, item := range args {
		if item != "" {
			if strings.Contains(fmt.Sprintf("%v", item), ",") {
				parts := strings.Split(fmt.Sprintf("%v", item), ",")

				for _, element := range parts {
					cleanArgs = append(cleanArgs, element)
				}
			} else {
				cleanArgs = append(cleanArgs, item)
			}
		}
	}

	rows, err := db.conn.Query(query, cleanArgs...)
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query execution"))
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var auditRow models.Audit
		var timestamp string
		if err := rows.Scan(&auditRow.ID, &auditRow.User, &auditRow.Action, &auditRow.Data, &timestamp); err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query row extraction"))
		}

		t, err := time.Parse(time.RFC3339, timestamp)

		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit parse timestamp"))
		}

		auditRow.Timestamp = t
		results = append(results, auditRow)
	}

	errRows := rows.Err()
	if errRows != nil {
		utils.LogError(errors.Wrap(errRows, "[AUDIT][QUERY] error in rows query loop"))
	}

	return results

}

func Query(query string) []string {
	if ok, out := db.isFaulty(db.faultyStatus); ok {
		utils.LogError(errors.Wrap(errors.New("Connection dropped due to " + out), "[AUDIT][QUERY]"))
		return nil
	}

	rows, err := db.conn.Query(query)
	if err != nil {
		utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query execution"))
		return nil
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var field string
		if err := rows.Scan(&field); err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query row extraction"))
		}

		results = append(results, field)
	}

	errRows := rows.Err()
	if errRows != nil {
		utils.LogError(errors.Wrap(errRows, "[AUDIT][QUERY] error in rows query loop"))
	}

	return results

}
