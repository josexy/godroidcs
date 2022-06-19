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

import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

public class SimpleLimiter implements Limiter {
    private final AtomicInteger current = new AtomicInteger(0);
    private final int limit;
    private final long window;
    private final AtomicLong lastUpdateTime;

    public SimpleLimiter(int limit, long windowMs) {
        this.limit = limit;
        this.window = windowMs;
        this.lastUpdateTime = new AtomicLong(System.currentTimeMillis());
    }

    @Override
    public boolean allow() {
        long now = System.currentTimeMillis();
        long last = lastUpdateTime.get();

        // reset
        if ((now - last) >= window) {
            lastUpdateTime.set(now);
            current.set(0);
        }

        if (current.getAndIncrement() > limit) {
            return false;
        }
        return true;
    }
}
