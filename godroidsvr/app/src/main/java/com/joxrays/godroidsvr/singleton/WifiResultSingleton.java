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

package com.joxrays.godroidsvr.singleton;

import com.joxrays.godroidsvr.message.SimpleWifiInfo;

import java.util.List;

public class WifiResultSingleton extends BaseSingleton {

    private List<SimpleWifiInfo> list;

    protected WifiResultSingleton() {
    }

    public static WifiResultSingleton getInstance() {
        return (WifiResultSingleton) BaseSingleton.getInstance(WifiResultSingleton.class);
    }

    public void set(List<SimpleWifiInfo> list) {
        this.list = list;
    }

    public List<SimpleWifiInfo> get() {
        return this.list;
    }
}
