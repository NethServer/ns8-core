/*
 * Copyright (C) 2023 Nethesis S.r.l.
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

import ()

type SocketAction struct {
	Action  string      `json:"action" structs:"action"`
	Payload interface{} `json:"payload" structs:"payload"`
}

type LogsStartAction struct {
	Id         string `json:"id" structs:"id"`
	Mode       string `json:"mode" structs:"mode"`
	Lines      string `json:"lines" structs:"lines"`
	Filter     string `json:"filter" structs:"filter"`
	From       string `json:"from" structs:"from"`
	To         string `json:"to" structs:"to"`
	Entity     string `json:"entity" structs:"entity"`
	EntityName string `json:"entity_name" structs:"entity_name"`
	TimeZone   string `json:"timezone" structs:"timezone"`
	Instance   string `json:"instance" structs:"instance"`
}

type LogsStopAction struct {
	Id  string `json:"id" structs:"id"`
	Pid string `json:"pid" structs:"pid"`
}
