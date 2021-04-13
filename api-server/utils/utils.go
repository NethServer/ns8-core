/*
 * Copyright (C) 2021 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part ofNethServer project.
 *
 *NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 *NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along withNethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package utils

import (
	"os"
	"strconv"
	"time"
)

func LogError(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
}

func Contains(a string, values []string) bool {
	for _, b := range values {
		if b == a {
			return true
		}
	}
	return false
}

func EpochToHumanDate(epochTime int) string {
	i, err := strconv.ParseInt(strconv.Itoa(epochTime), 10, 64)
	if err != nil {
		return "-"
	}
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02 15:04:05")
}
