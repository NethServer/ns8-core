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
	"bufio"
	"fmt"
	"os"
	"os/exec"

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
	fmt.Printf("<5>user=%s action=%s data=%s\n", audit.User, audit.Action, audit.Data)
}

func QueryArgs(logql string, qargs []string) []models.Audit {
	// define results
	var results []models.Audit

	qargs = append(qargs, logql)

	cmd := exec.Command("/usr/local/bin/logcli", qargs...)
	stdoutReader, _ := cmd.StdoutPipe()
	lineReader := bufio.NewReader(stdoutReader)
	fmt.Printf("%v\n", qargs)
	cmd.Start()


	for {
		line, err := lineReader.ReadString(10) // LF
		if line != "" {
			results = append(results, models.Audit{Data: line})
		}
		if err != nil {
			break
		}
	}

	cmd.Wait()

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
