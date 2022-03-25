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

import android.net.wifi.WifiManager;
import android.os.Build;
import android.util.Pair;

import androidx.annotation.RequiresApi;

import com.joxrays.godroidsvr.message.NetInterfaceInfoList;
import com.joxrays.godroidsvr.singleton.WifiResultSingleton;
import com.joxrays.godroidsvr.message.Boolean;
import com.joxrays.godroidsvr.message.DetailActiveNetworkInfoList;
import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.PublicNetworkInfo;
import com.joxrays.godroidsvr.message.ScanWifiInfoList;
import com.joxrays.godroidsvr.message.DetailWifiInfo;
import com.joxrays.godroidsvr.util.CommonUtil;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;
import com.joxrays.godroidsvr.util.NetworkUtil;

import io.grpc.stub.StreamObserver;

public class BgWorkNetResolverService extends NetResolverGrpc.NetResolverImplBase {

    private final BgWorkBaseResolverGroup group;

    public BgWorkNetResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    @Override
    public void getCurrentWifiInfo(Empty request, StreamObserver<DetailWifiInfo> responseObserver) {
        Pair<DetailWifiInfo, Exception> pair = NetworkUtil.getWifiInfo(this.group.getContext());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    @Override
    public void getNetworkInfo(Empty request, StreamObserver<NetInterfaceInfoList> responseObserver) {
        Pair<NetInterfaceInfoList, Exception> pair = NetworkUtil.getNetworkInfo(this.group.getContext());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void scanWifiResult(Empty request, StreamObserver<ScanWifiInfoList> responseObserver) {
        WifiManager wm = CommonUtil.getWifiManager(this.group.getContext());
        wm.startScan();

        ScanWifiInfoList.Builder builder = ScanWifiInfoList.newBuilder();
        if (WifiResultSingleton.getInstance().get() == null) {
            builder.setEmpty(true);
        } else {
            builder.setEmpty(false).addAllValues(WifiResultSingleton.getInstance().get());
        }
        responseObserver.onNext(builder.build());
        responseObserver.onCompleted();
    }

    @Override
    public void checkNetworkConnectivity(Empty request, StreamObserver<Boolean> responseObserver) {
        responseObserver.onNext(Boolean.newBuilder()
                .setValue(NetworkUtil.hasActiveNetwork(this.group.getContext()))
                .build());
        responseObserver.onCompleted();
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    @Override
    public void getActiveNetworkInfo(Empty request, StreamObserver<DetailActiveNetworkInfoList> responseObserver) {
        responseObserver.onNext(DetailActiveNetworkInfoList.newBuilder()
                .addAllValues(NetworkUtil.getActivityNetworkDetailInfo(this.group.getContext()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void getPublicNetworkInfo(Empty request, StreamObserver<PublicNetworkInfo> responseObserver) {
        Pair<PublicNetworkInfo, Exception> pair = NetworkUtil.getPublicAddressInfo();
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
            return;
        }
        responseObserver.onNext(pair.first);
        responseObserver.onCompleted();
    }
}
