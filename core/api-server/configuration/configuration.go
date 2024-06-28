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

package configuration

import (
	"os"
	"strings"
)

type Configuration struct {
	ListenAddress string   `json:"listen_address"`
	RedisAddress  string   `json:"redis_address"`
	RedisUser     string   `json:"redis_user"`
	RedisPassword string   `json:"redis_password"`
	Secret        string   `json:"secret"`
	StaticPath    string   `json:"static_path"`
	AuditFile     string   `json:"audit_file"`
	Issuer        string   `json:"issuer"`
	SensitiveList []string `json:"sensitive_list"`
}

var Config = Configuration{}

func Init() {
	// read configuration from ENV
	if os.Getenv("LISTEN_ADDRESS") != "" {
		Config.ListenAddress = os.Getenv("LISTEN_ADDRESS")
	} else {
		Config.ListenAddress = "127.0.0.1:8080"
	}

	if os.Getenv("REDIS_ADDRESS") != "" {
		Config.RedisAddress = os.Getenv("REDIS_ADDRESS")
	} else {
		Config.RedisAddress = "127.0.0.1:6379"
	}
	if os.Getenv("REDIS_USER") != "" {
		Config.RedisUser = os.Getenv("REDIS_USER")
	} else {
		Config.RedisUser = "default"
	}
	if os.Getenv("REDIS_PASSWORD") != "" {
		Config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	} else {
		Config.RedisPassword = ""
	}

	if os.Getenv("SECRET") != "" {
		Config.Secret = os.Getenv("SECRET")
	} else {
		os.Stderr.WriteString("SECRET variable is empty. ")
		os.Exit(1)
	}

	if os.Getenv("STATIC_PATH") != "" {
		Config.StaticPath = os.Getenv("STATIC_PATH")
	} else {
		Config.StaticPath = "/var/lib/nethserver/cluster/ui"
	}

	if os.Getenv("AUDIT_FILE") != "" {
		Config.AuditFile = os.Getenv("AUDIT_FILE")
	} else {
		Config.AuditFile = ""
	}

	if os.Getenv("ISSUER") != "" {
		Config.Issuer = os.Getenv("ISSUER")
	} else {
		Config.Issuer = "NethServer"
	}

	if os.Getenv("SENSITIVE_LIST") != "" {
		Config.SensitiveList = strings.Split(os.Getenv("SENSITIVE_LIST"), ",")
	} else {
		Config.SensitiveList = []string{"password", "secret", "token", "jwt", "admpass", "adminpass"}
	}
}
