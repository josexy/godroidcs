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

package com.joxrays.godroidsvr.resolver;

import android.content.Context;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.ConcurrentHashMap;

import io.grpc.BindableService;
import io.grpc.ServerServiceDefinition;

public class BgWorkBaseResolverGroup implements BgWorkContext {
    private final Context context;
    private final Map<String, BindableService> serviceMap;

    public BgWorkBaseResolverGroup(Context context) {
        this.context = context;
        this.serviceMap = new ConcurrentHashMap<>();
    }

    public Context getContext() {
        return context;
    }

    /**
     * register grpc Service
     * @param name
     * @param cls
     * @param <T>
     */
    public <T extends BindableService> void registerService(String name, Class<T> cls) {
        if (name == null || cls == null) return;
        if (serviceMap.containsKey(name)) return;
        try {
            T v = cls.getDeclaredConstructor(BgWorkBaseResolverGroup.class).newInstance(this);
            serviceMap.put(name, v);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public ServerServiceDefinition getService(String name) {
        return Objects.requireNonNull(serviceMap.get(name)).bindService();
    }

    public List<ServerServiceDefinition> getServices() {
        List<ServerServiceDefinition> list = new ArrayList<>();
        serviceMap.forEach((k, v) -> {
            list.add(v.bindService());
        });
        return list;
    }
}
