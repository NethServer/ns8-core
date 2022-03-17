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
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/NethServer/ns8-core/core/api-server/configuration"
	"github.com/NethServer/ns8-core/core/api-server/models"
	"github.com/NethServer/ns8-core/core/api-server/utils"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Init() {
	// check audit file is set
	if len(configuration.Config.AuditFile) > 0 {
		// check if file exists
		if _, err := os.Stat(configuration.Config.AuditFile); os.IsNotExist(err) {
			// create audit db file and initialize it
			create()
		} else {
			// check if db is valid
			db, err := sql.Open("sqlite3", configuration.Config.AuditFile)
			if err != nil {
				// audit file invalid, create and initialize it
				create()
			}
			defer db.Close()

			// try a query
			query := "SELECT * FROM audit LIMIT 1;"

			_, err = db.Exec(query)
			if err != nil {
				// invalid sql schema, create and initialize it
				create()
			}
		}
	} else {
		utils.LogError(errors.Wrap(errors.New("AUDIT_FILE is not set in the environment."), "[Audit DISABLED]"))
	}
}

func create() {
	// purge existing file
	os.Remove(configuration.Config.AuditFile)

	// create new db
	db, errCreate := sql.Open("sqlite3", configuration.Config.AuditFile)
	if errCreate != nil {
		utils.LogError(errors.Wrap(errCreate, "error in audit db file creation"))
	}
	defer db.Close()

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

	// execute create table
	_, errExecute := db.Exec(query)
	if errExecute != nil {
		utils.LogError(errors.Wrap(errExecute, "error in audit file schema init"))
	}
}

func Store(audit models.Audit) {
	// check audit file is set
	if len(configuration.Config.AuditFile) > 0 {
		// open db
		db, err := sql.Open("sqlite3", configuration.Config.AuditFile)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error in audit file schema open"))
		}
		defer db.Close()

		// begin sqlite connection to insert
		tx, err := db.Begin()
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error in audit file schema begin"))
		}

		// define statement
		stmt, err := tx.Prepare("INSERT INTO audit(id, user, action, data, timestamp) VALUES(null, ?, ?, ?, ?)")
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error in audit file schema prepare"))
		}
		defer stmt.Close()

		// execute statement
		_, err = stmt.Exec(audit.User, audit.Action, audit.Data, audit.Timestamp.Format(time.RFC3339))
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][STORE] error in audit file schema execute"))
		}
		tx.Commit()
	}
}

func QueryArgs(query string, args ...interface{}) []models.Audit {
	// define results
	var results []models.Audit

	// check audit file is set
	if len(configuration.Config.AuditFile) > 0 {
		// open db
		db, err := sql.Open("sqlite3", configuration.Config.AuditFile)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit file schema open"))
		}
		defer db.Close()

		// define query
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

		rows, err := db.Query(query, cleanArgs...)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query execution"))
		}
		defer rows.Close()

		// loop rows
		for rows.Next() {
			var auditRow models.Audit
			var timestamp string
			if err := rows.Scan(&auditRow.ID, &auditRow.User, &auditRow.Action, &auditRow.Data, &timestamp); err != nil {
				utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query row extraction"))
			}

			// parse date
			t, err := time.Parse(time.RFC3339, timestamp)

			if err != nil {
				utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit parse timestamp"))
			}

			// append results
			auditRow.Timestamp = t
			results = append(results, auditRow)
		}

		// check rows error
		errRows := rows.Err()
		if errRows != nil {
			utils.LogError(errors.Wrap(errRows, "[AUDIT][QUERY] error in rows query loop"))
		}
	}

	// return results
	return results

}

func Query(query string) []string {
	// define results
	var results []string

	// check audit file is set
	if len(configuration.Config.AuditFile) > 0 {
		// open db
		db, err := sql.Open("sqlite3", configuration.Config.AuditFile)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit file schema open"))
		}
		defer db.Close()

		// define query
		rows, err := db.Query(query)
		if err != nil {
			utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query execution"))
		}
		defer rows.Close()

		// loop rows
		for rows.Next() {
			var field string
			if err := rows.Scan(&field); err != nil {
				utils.LogError(errors.Wrap(err, "[AUDIT][QUERY] error in audit query row extraction"))
			}

			// append results
			results = append(results, field)
		}

		// check rows error
		errRows := rows.Err()
		if errRows != nil {
			utils.LogError(errors.Wrap(errRows, "[AUDIT][QUERY] error in rows query loop"))
		}
	}

	// return results
	return results

}
