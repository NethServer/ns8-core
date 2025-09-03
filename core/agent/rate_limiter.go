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
)

const (
	MAX_TOKENS int = 30
	INITIAL_TOKENS int = 10
)

type TokenBucket struct { 
	initialTokens int
	fillerSleepTime time.Duration
	tokens chan struct{}
}

func NewTokenBucketAlgorithm(fillingInterval time.Duration) *TokenBucket {
	t := &TokenBucket{
		fillerSleepTime: fillingInterval,
		initialTokens: INITIAL_TOKENS,
		tokens: make(chan struct{}, MAX_TOKENS),
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