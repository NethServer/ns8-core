/*
 * Copyright (C) 2022 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 */

package main

import (
	"context"
	"log"

	"github.com/NethServer/ns8-core/core/agent/models"
	"github.com/go-redis/redis/v8"
)

func listenEventsAsync(ctx context.Context, rdb *redis.Client, complete chan int) {
	pubsub := rdb.PSubscribe(ctx, "*/event/*")
	defer func() {
		pubsub.Close()
		complete <- 1
	}()

	csyn := make(chan int, 1)

	go func() {
		for msg := range pubsub.Channel() {
			log.Print(msg.Channel)
		}
		csyn <- 1
	}()

	select {
	case <-csyn:
		return
	case <-ctx.Done():
		return
	}

}

func runEvent(actionCtx context.Context, task *models.Event) {

}
