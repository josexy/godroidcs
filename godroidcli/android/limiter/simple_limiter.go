// Copyright [2021] [josexy]
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package limiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type SimpleLimiter struct {
	current   int32
	limit     int32
	window    time.Duration
	once      sync.Once
	timer     *time.Ticker
	stopTimer chan struct{}
}

func NewSimpleLimiter(limit int, window time.Duration) *SimpleLimiter {
	return &SimpleLimiter{
		limit:     int32(limit),
		window:    window,
		timer:     time.NewTicker(window),
		stopTimer: make(chan struct{}),
	}
}

func (limiter *SimpleLimiter) Allow() bool {
	limiter.once.Do(func() {
		go func() {
			defer limiter.timer.Stop()
			for {
				select {
				case <-limiter.timer.C:
					atomic.StoreInt32(&limiter.current, 0)
				case <-limiter.stopTimer:
					// interrupt for-loop and stop timer.Ticker completely
					return
				}
			}
		}()
	})
	// discard
	if atomic.LoadInt32(&limiter.current) > limiter.limit {
		return false
	}
	atomic.AddInt32(&limiter.current, 1)
	return true
}

func (limiter *SimpleLimiter) Done() {
	close(limiter.stopTimer)
}
