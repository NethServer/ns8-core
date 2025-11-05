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
	"sync"
	"sync/atomic"
	"time"
)

type workersLimiter struct {
	wg      sync.WaitGroup
	counter int64
	limit   int64
	delay   time.Duration
}

func NewWorkersLimiter(limit int, delay time.Duration) *workersLimiter {
	w := new(workersLimiter)
	w.counter = 0
	w.limit = int64(limit)
	w.delay = delay
	return w
}

func (w *workersLimiter) Add(delta int) {
	// call wg.Add first so it panics like the standard WaitGroup if delta causes it
	w.wg.Add(delta)
	atomic.AddInt64(&w.counter, int64(delta))
}

func (w *workersLimiter) Done() {
	// call wg.Done first to preserve panic semantics
	w.wg.Done()
	atomic.AddInt64(&w.counter, -1)
}

func (w *workersLimiter) Wait() {
	w.wg.Wait()
}

// ObserveOverload checks if the current number of running workers exceeds
// the configured limit. If the limit is reached, it waits for a short
// delay before re-checking. Returns true if the system is still
// overloaded after the delay, signaling that new work should be rejected.
func (w *workersLimiter) ObserveOverload() bool {
	if atomic.LoadInt64(&w.counter) >= w.limit {
		time.Sleep(w.delay)
		if atomic.LoadInt64(&w.counter) >= w.limit {
			return true
		}
	}
	return false
}
