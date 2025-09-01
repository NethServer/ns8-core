/*
 * Copyright (C) 2023 Nethesis S.r.l.
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
	"sync"
)

const (
	MAX_TOKENS uint32 = 1000 
	MIN_TOKEN uint32 = 0
)

type TokenBucket struct { 
	tokens uint32 
	tickerDelay time.Duration
	lock sync.Mutex
	ticker *time.Ticker
}


func NewTokenBucketAlgorithm(fillingInterval time.Duration, baseTokens uint32) *TokenBucket {
	if fillingInterval <= 0 || baseTokens < 0 {
		return nil 
	}

	return &TokenBucket{
		tokens: baseTokens,
		tickerDelay: fillingInterval * time.Second,
	}
}

func (t *TokenBucket) TryTakeToken() bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	if safe := t.checkTokensSafety(); !safe {
		return safe
	}

	t.takeToken()
	return true
}

func (t *TokenBucket) addToken() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if safe := t.checkTokensSafety(); !safe {
		return
	}

	t.tokens += 1
}

func (t *TokenBucket) takeToken() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.tokens -= 1
}

func (t *TokenBucket) TokenFiller() {
	t.ticker = time.NewTicker(t.tickerDelay)
	defer t.ticker.Stop()

	for {
		select {
		case <- t.ticker.C:
			t.addToken()
		}
	}
}

func (t *TokenBucket) checkTokensSafety() bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	switch {
	case t.tokens >= MAX_TOKENS, t.tokens <= MIN_TOKEN:
		return false
	}

	return true
}