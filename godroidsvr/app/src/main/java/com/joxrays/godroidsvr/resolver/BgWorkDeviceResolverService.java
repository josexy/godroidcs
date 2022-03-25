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

import android.location.LocationManager;
import android.util.Pair;

import com.joxrays.godroidsvr.message.BatteryInfo;
import com.joxrays.godroidsvr.message.DeviceInfo;
import com.joxrays.godroidsvr.message.DisplayInfo;
import com.joxrays.godroidsvr.message.GPUInfo;
import com.joxrays.godroidsvr.message.IntegerList;
import com.joxrays.godroidsvr.message.LocationInfo;
import com.joxrays.godroidsvr.message.MemoryInfo;
import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.StorageSpaceInfo;
import com.joxrays.godroidsvr.message.SystemInfo;
import com.joxrays.godroidsvr.util.BatteryUtil;
import com.joxrays.godroidsvr.util.DeviceUtil;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;
import com.joxrays.godroidsvr.util.LocationUtil;

import java.util.List;

import io.grpc.stub.StreamObserver;

public class BgWorkDeviceResolverService extends DeviceResolverGrpc.DeviceResolverImplBase {

    private final BgWorkBaseResolverGroup group;

    public BgWorkDeviceResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getDeviceInfo(Empty request, StreamObserver<DeviceInfo> responseObserver) {
        responseObserver.onNext(DeviceUtil.getDeviceInfo(this.group.getContext()));
        responseObserver.onCompleted();
    }

    @Override
    public void getMemoryInfo(Empty request, StreamObserver<MemoryInfo> responseObserver) {
        responseObserver.onNext(DeviceUtil.getMemoryInfo(this.group.getContext()));
        responseObserver.onCompleted();
    }

    @Override
    public void getStorageSpaceInfo(Empty request, StreamObserver<StorageSpaceInfo> responseObserver) {
        Pair<StorageSpaceInfo, Exception> pair = DeviceUtil.getStorageSpaceSize(this.group.getContext());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getLocationInfo(Empty request, StreamObserver<LocationInfo> responseObserver) {
        Pair<LocationInfo, Exception> pair = LocationUtil.getAddressInfo(this.group.getContext(), LocationManager.NETWORK_PROVIDER);
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
            return;
        }
        responseObserver.onNext(pair.first);
        responseObserver.onCompleted();
    }

    @Override
    public void getCPUsFrequency(Empty request, StreamObserver<IntegerList> responseObserver) {
        Pair<List<Integer>, Exception> pair = DeviceUtil.getCPUsFrequency();
        if (pair == null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ErrorExceptionUtil.ErrorCannotListDir));
        } else if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(IntegerList.newBuilder().addAllValues(pair.first).build());
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getGPUInfo(Empty request, StreamObserver<GPUInfo> responseObserver) {
        responseObserver.onNext(DeviceUtil.getGPUInfo());
        responseObserver.onCompleted();
    }

    @Override
    public void getSystemInfo(Empty request, StreamObserver<SystemInfo> responseObserver) {
        responseObserver.onNext(DeviceUtil.getSystemInfo(this.group.getContext()));
        responseObserver.onCompleted();
    }

    @Override
    public void getDisplayInfo(Empty request, StreamObserver<DisplayInfo> responseObserver) {
        responseObserver.onNext(DeviceUtil.getDisplayInfo(this.group.getContext()));
        responseObserver.onCompleted();
    }

    @Override
    public void getBatteryInfo(Empty request, StreamObserver<BatteryInfo> responseObserver) {
        responseObserver.onNext(BatteryUtil.getBatteryInfo(this.group.getContext()));
        responseObserver.onCompleted();
    }
}
