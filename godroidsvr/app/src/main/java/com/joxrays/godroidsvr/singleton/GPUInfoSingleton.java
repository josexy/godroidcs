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

import com.joxrays.godroidsvr.message.GPUInfo;

public class GPUInfoSingleton extends BaseSingleton {

    private final GPUInfo.Builder builder = GPUInfo.newBuilder();

    protected GPUInfoSingleton() {

    }

    public static GPUInfoSingleton getInstance() {
        return (GPUInfoSingleton) BaseSingleton.getInstance(GPUInfoSingleton.class);
    }

    public void setInfo(String vendor, String renderer, String version) {
        builder.setVendor(vendor != null ? vendor : "")
                .setRenderer(renderer != null ? renderer : "")
                .setVersion(version != null ? version : "");
    }

    public GPUInfo getGPUInfo() {
        return builder.build();
    }
}
