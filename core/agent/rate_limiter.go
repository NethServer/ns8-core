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
	"os"
	"sync"
)

const (
	DEFAULT_CAPACITY     = 10
	DEFAULT_TIME_RATE    = 300
	ENV_BUCKET_INTERVAL  = "BUCKET_INTERVAL"
	ENV_BUCKET_CAPACITY  = "BUCKET_CAPACITY"
)

type tokenBucket struct { 
	initialTokens int
	fillerSleepTime time.Duration
	tokens chan struct{}
}

var tokenBucketInstance *tokenBucket
var mu = &sync.Mutex{}

func NewTokenBucketAlgorithm() *tokenBucket {
	mu.Lock()
	defer mu.Unlock()

	if tokenBucketInstance != nil {
		return tokenBucketInstance
	}

	interval := os.Getenv(ENV_BUCKET_INTERVAL)
	capacity := os.Getenv(ENV_BUCKET_CAPACITY)
	confManager := &rateLimitConf{}

	confManager.loadAndValidate(interval, capacity)
	
	t := &tokenBucket{
		fillerSleepTime: confManager.effectiveTime,
		initialTokens: confManager.capacity,
		tokens: make(chan struct{}, confManager.capacity),
	}

	tokenBucketInstance = t

	go tokenBucketInstance.tokenFiller()
	return tokenBucketInstance
}

func (t *tokenBucket) tokenFiller() {
	t.initialFiller()
	for {
		time.Sleep(t.fillerSleepTime)

		t.tokens <- struct{}{}
	}
}

func (t *tokenBucket) takeToken() {
	<- t.tokens
}

func (t *tokenBucket) initialFiller() {
	for i := 0; i < t.initialTokens; i++ {
		t.tokens <- struct{}{}
	}
}

type rateLimitConf struct {
	timeRate int
	capacity int
	effectiveTime time.Duration
}

func (c *rateLimitConf) loadAndValidate(fillingInterval, capacity string) {
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