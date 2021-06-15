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
	"os"

	"github.com/pkg/errors"

	"github.com/NethServer/ns8-scratchpad/core/api-server/configuration"
	"github.com/NethServer/ns8-scratchpad/core/api-server/utils"

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

func Store(user string, action string, data string, timestamp string) {
	// check audit file is set
	if len(configuration.Config.AuditFile) > 0 {
		// open db
		db, err := sql.Open("sqlite3", configuration.Config.AuditFile)
		if err != nil {
			utils.LogError(errors.Wrap(err, "error in audit file schema open"))
		}
		defer db.Close()

		// begin sqlite connection to insert
		tx, err := db.Begin()
		if err != nil {
			utils.LogError(errors.Wrap(err, "error in audit file schema begin"))
		}

		// define statement
		stmt, err := tx.Prepare("INSERT INTO audit(id, user, action, data, timestamp) VALUES(null, ?, ?, ?, ?)")
		if err != nil {
			utils.LogError(errors.Wrap(err, "error in audit file schema prepare"))
		}
		defer stmt.Close()

		// execute statement
		_, err = stmt.Exec(user, action, data, timestamp)
		if err != nil {
			utils.LogError(errors.Wrap(err, "error in audit file schema execute"))
		}
		tx.Commit()
	}
}
