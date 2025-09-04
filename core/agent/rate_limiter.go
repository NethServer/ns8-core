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
	fillingInterval int
	timeUnit string
	capacity int
	effectiveTime time.Duration
}

func (c *RateLimitConf) loadAndValidate(fillingInt string, timeUnit string, capacity string) {
	conv, err := strconv.Atoi(fillingInt)
	if fillingInt == "" || err != nil {
		c.fillingInterval = 300
	} else {
		c.fillingInterval = conv
	}

	ok := c.checkTimeUnitValidity(timeUnit)
	if timeUnit == "" || !ok {
		c.timeUnit = "millis"
	} else {
		c.timeUnit = timeUnit
	}

	convCap, err := strconv.Atoi(capacity)
	if capacity == "" || err != nil {
		c.capacity = 30
	} else {
		c.capacity = convCap
	}

	c.buildEffectiveTime()
}

func (c *RateLimitConf) checkTimeUnitValidity(timeUnit string) bool {
	switch timeUnit {
	case "millis":
		return true
	case "second":
		return true
	case "minute":
		return true
	}

	return false
}

func (c *RateLimitConf) buildEffectiveTime() {
	switch c.timeUnit {
	case "millis":
		c.effectiveTime = time.Duration(c.fillingInterval) * time.Millisecond
	case "second":
		c.effectiveTime = time.Duration(c.fillingInterval) * time.Second
	case "minute":
		c.effectiveTime = time.Duration(c.fillingInterval) * time.Minute
	}
}