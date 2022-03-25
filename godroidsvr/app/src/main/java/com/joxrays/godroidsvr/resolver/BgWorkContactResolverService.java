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

import com.joxrays.godroidsvr.message.ContactInfo;
import com.joxrays.godroidsvr.message.ContactMetaInfoList;
import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.String;
import com.joxrays.godroidsvr.util.ContactUtil;

import io.grpc.stub.StreamObserver;

public class BgWorkContactResolverService extends ContactResolverGrpc.ContactResolverImplBase {
    private final BgWorkBaseResolverGroup group;

    public BgWorkContactResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getAllContactInfo(Empty request, StreamObserver<ContactMetaInfoList> responseObserver) {
        responseObserver.onNext(ContactMetaInfoList.newBuilder()
                .addAllValues(ContactUtil.getAllContactMetaInfo(this.group.getContext()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void getContactInfo(String request, StreamObserver<ContactInfo> responseObserver) {
        responseObserver.onNext(ContactInfo
                .newBuilder(ContactUtil.getContactInfo(this.group.getContext(), request.getValue()))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void deleteContact(String request, StreamObserver<Empty> responseObserver) {
        ContactUtil.deleteContact(this.group.getContext(), request.getValue());
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void addContact(ContactInfo request, StreamObserver<Empty> responseObserver) {
        ContactUtil.addContact(this.group.getContext(), request);
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }
}
