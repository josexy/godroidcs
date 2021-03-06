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
	"testing"
	"time"
)

func TestSimpleLimiter(t *testing.T) {
	limiter := NewSimpleLimiter(10, time.Millisecond*50)
	for i := 0; i < 100; i++ {
		if limiter.Allow() {
			t.Log("ok")
		} else {
			t.Log("fail")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestTokenBucketLimiter(t *testing.T) {
	limiter := NewTokenBucketLimiter(30, 15)
	for i := 0; i < 100; i++ {
		if limiter.Allow() {
			t.Log("ok")
		} else {
			t.Log("fail")
		}
		time.Sleep(time.Millisecond)
	}
}
