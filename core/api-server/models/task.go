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

package models

import (
	"time"
)

type Task struct {
	ID        string      `json:"id" structs:"id"`
	Action    string      `json:"action" structs:"action"`
	Data      interface{} `json:"data" structs:"data"`
	Extra     interface{} `json:"extra" structs:"extra"`
	Queue     string      `json:"queue" structs:"queue"`
	User      string      `json:"user" structs:"user"`
	Timestamp time.Time   `json:"timestamp" structs:"timestamp"`
	Parent    string      `json:"parent" structs:"parent"`
}

type TaskJSON struct {
	ID     string      `json:"id" structs:"id" example:"null"`
	Action string      `json:"action" structs:"action" example:"list-modules"`
	Data   interface{} `json:"data" structs:"data"`
	Extra  interface{} `json:"extra" structs:"extra"`
	Parent string      `json:"parent" structs:"parent"`
}
