/*
 * Copyright (C) 2023 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package utils

import (
	"sync"
	"time"
)

type MuClock struct {
	mux     sync.Mutex
	curTime int64
}

func (clk *MuClock) Sync() {
	clk.mux.Lock()
	clk.curTime = time.Now().Unix()
	clk.mux.Unlock()
}

func (clk *MuClock) Now() int64 {
	clk.mux.Lock()
	defer clk.mux.Unlock()
	return clk.curTime
}
