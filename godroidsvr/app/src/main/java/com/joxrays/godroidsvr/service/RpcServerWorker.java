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

package com.joxrays.godroidsvr.service;

import android.content.Context;

import androidx.annotation.NonNull;
import androidx.work.Worker;
import androidx.work.WorkerParameters;

import com.joxrays.godroidsvr.MainActivity;
import com.joxrays.godroidsvr.resolver.BgWorkBaseResolverGroup;
import com.joxrays.godroidsvr.resolver.BgWorkCallLogResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkContactResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkControlResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkDeviceResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkFsResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkMediaStoreResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkNetResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkPhoneResolverService;
import com.joxrays.godroidsvr.resolver.BgWorkPingTestService;
import com.joxrays.godroidsvr.resolver.BgWorkSmsResolverService;
import com.joxrays.godroidsvr.util.LogUtil;
import com.joxrays.godroidsvr.resolver.BgWorkPmResolverService;

import java.util.concurrent.TimeUnit;

import io.grpc.Server;
import io.grpc.netty.shaded.io.grpc.netty.NettyServerBuilder;
import io.grpc.netty.shaded.io.netty.channel.ChannelOption;

public class RpcServerWorker extends Worker {
    public final static String TAG = "RpcServerWorker";
    private Server server;
    private final int rpcPort;
    private boolean isRunning;

    public RpcServerWorker(@NonNull Context context, @NonNull WorkerParameters workerParams) {
        super(context, workerParams);
        rpcPort = workerParams.getInputData().getInt("PORT", MainActivity.RPC_SERVER_PORT);
        isRunning = false;
    }

    public Exception startServer(int port) {
        // register all rpc service
        BgWorkBaseResolverGroup workBaseResolverGroup = new BgWorkBaseResolverGroup(this.getApplicationContext());
        workBaseResolverGroup.registerService("ping", BgWorkPingTestService.class);
        workBaseResolverGroup.registerService("pm", BgWorkPmResolverService.class);
        workBaseResolverGroup.registerService("fs", BgWorkFsResolverService.class);
        workBaseResolverGroup.registerService("net", BgWorkNetResolverService.class);
        workBaseResolverGroup.registerService("device", BgWorkDeviceResolverService.class);
        workBaseResolverGroup.registerService("ctrl", BgWorkControlResolverService.class);
        workBaseResolverGroup.registerService("ms", BgWorkMediaStoreResolverService.class);
        workBaseResolverGroup.registerService("sms", BgWorkSmsResolverService.class);
        workBaseResolverGroup.registerService("contact", BgWorkContactResolverService.class);
        workBaseResolverGroup.registerService("calllog", BgWorkCallLogResolverService.class);
        workBaseResolverGroup.registerService("phone", BgWorkPhoneResolverService.class);

        LogUtil.d("start server on: " + port);

        server = NettyServerBuilder.forPort(port)
                .addServices(workBaseResolverGroup.getServices())
                .withChildOption(ChannelOption.SO_REUSEADDR, true)
                .build();
        try {
            server.start();
            server.awaitTermination();
            return null;
        } catch (Exception ex) {
            return ex;
        }
    }

    public void stopServer() {
        if (server != null) {
            try {
                LogUtil.d("stop server");
                server.shutdown().awaitTermination(5, TimeUnit.SECONDS);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
    }

    @Override
    public void onStopped() {
        super.onStopped();
        stopServer();
        isRunning = false;
    }

    @NonNull
    @Override
    public Result doWork() {
        if (isRunning) {
            return Result.failure();
        }
        isRunning = true;
        Exception ex = startServer(rpcPort);
        // start rpc server failed
        if (ex != null) {
            ex.printStackTrace();
            return Result.failure();
        }
        // if start rpc server successfully, the server will be blocked.
        // so this method "doWork" don't return any value
        return Result.success();
    }
}
