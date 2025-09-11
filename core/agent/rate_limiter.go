/*
 * Copyright (C) 2025 Nethesis S.r.l.
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
	"time"
	"strconv"
	"strings"
)

const (
	DEFAULT_CAPACITY int = 10
	DEFAULT_TIME_RATE int = 300
)

type TokenBucket struct { 
	initialTokens int
	fillerSleepTime time.Duration
	tokens chan struct{}
}

func NewTokenBucketAlgorithm(fillingInterval time.Duration, capacity int) *TokenBucket {
	t := &TokenBucket{
		fillerSleepTime: fillingInterval,
		initialTokens: capacity,
		tokens: make(chan struct{}, capacity),
	}

	go t.tokenFiller()
	return t
}

func (t *TokenBucket) tokenFiller() {
	t.initialFiller()
	for {
		time.Sleep(t.fillerSleepTime)

		t.tokens <- struct{}{}
	}
}

func (t *TokenBucket) takeToken() {
	<- t.tokens
}

func (t *TokenBucket) initialFiller() {
	for i := 0; i < t.initialTokens; i++ {
		t.tokens <- struct{}{}
	}
}

type RateLimitConf struct {
	timeRate int
	capacity int
	effectiveTime time.Duration
}

func (c *RateLimitConf) loadAndValidate(fillingInterval string, capacity string) {
	before, _, found := strings.Cut(fillingInterval, "ms")
	convBefore, err := strconv.Atoi(before)

	if fillingInterval == "" || !found || err != nil {
		c.timeRate = DEFAULT_TIME_RATE
	} else {
		c.timeRate = convBefore
	}

	convCap, err := strconv.Atoi(capacity)
	if capacity == "" || err != nil {
		c.capacity = DEFAULT_CAPACITY
	} else {
		c.capacity = convCap
	}

	c.effectiveTime = time.Duration(c.timeRate) * time.Millisecond
}