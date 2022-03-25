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

import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.SmsInfoList;
import com.joxrays.godroidsvr.message.String;
import com.joxrays.godroidsvr.message.StringList;
import com.joxrays.godroidsvr.message.StringPair;
import com.joxrays.godroidsvr.util.SmsUtil;

import io.grpc.stub.StreamObserver;

public class BgWorkSmsResolverService extends SmsResolverGrpc.SmsResolverImplBase {
    private final BgWorkBaseResolverGroup group;

    public BgWorkSmsResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getAllBasicSmsInfo(Empty request, StreamObserver<StringList> responseObserver) {
        responseObserver.onNext(StringList.newBuilder()
                .addAllValues(SmsUtil.getAllSmsMetaInfo(this.group.getContext()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void getSmsInfoList(String request, StreamObserver<SmsInfoList> responseObserver) {
        responseObserver.onNext(SmsInfoList.newBuilder()
                .addAllValues(SmsUtil.getSmsInfo(this.group.getContext(), request.getValue()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void sendSms(StringPair request, StreamObserver<Empty> responseObserver) {
        SmsUtil.sendSms(this.group.getContext(), request.getFirst(), request.getSecond());
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }
}
