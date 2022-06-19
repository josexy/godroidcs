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

package com.joxrays.godroidsvr.limiter;

public class TokenBucketLimiter implements Limiter {
    private final long capacity;                      // bucket capacity
    private final long rate;                          // permits-per-second
    private long currentTokenNum;
    private long lastAddTokenTime;

    public TokenBucketLimiter(long capacity, long rate) {
        this.capacity = capacity;
        this.rate = rate;
        this.currentTokenNum = capacity;
        this.lastAddTokenTime = System.currentTimeMillis();
    }

    @Override
    public boolean allow() {
        return acquire(1);
    }

    public boolean acquire(int permits) {
        if (permits > currentTokenNum) {
            long accessTime = System.currentTimeMillis();
            long durationMs = accessTime - lastAddTokenTime;
            long newTokenNum = (long) (durationMs * rate * 1.0 / 1000);
            if (newTokenNum > 0) {
                currentTokenNum = Math.min(currentTokenNum + newTokenNum, capacity);
                this.lastAddTokenTime = accessTime;
            }
            if (permits > currentTokenNum) return false;
        }
        this.currentTokenNum -= permits;
        return true;
    }
}

