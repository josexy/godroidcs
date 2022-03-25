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

import android.graphics.Bitmap;
import android.net.Uri;
import android.util.Pair;

import com.joxrays.godroidsvr.message.Bytes;
import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.MediaStoreInfoList;
import com.joxrays.godroidsvr.message.MediaType;
import com.joxrays.godroidsvr.message.String;
import com.joxrays.godroidsvr.observer.DownloadStreamHandler;
import com.joxrays.godroidsvr.util.BitmapUtil;
import com.joxrays.godroidsvr.util.ContentResolverUtil;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;
import com.joxrays.godroidsvr.util.MediaStoreUtil;

import java.io.IOException;
import java.io.InputStream;

import io.grpc.stub.StreamObserver;

public class BgWorkMediaStoreResolverService extends MediaStoreResolverGrpc.MediaStoreResolverImplBase {

    private final BgWorkBaseResolverGroup group;

    public BgWorkMediaStoreResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getMediaFilesInfo(MediaType request, StreamObserver<MediaStoreInfoList> responseObserver) {
        responseObserver.onNext(MediaStoreInfoList.newBuilder()
                .addAllValues(MediaStoreUtil.getAllMediaFilesInfo(
                        this.group.getContext(), request.getTypeValue()
                ))
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void getMediaFileThumbnail(String request, StreamObserver<Bytes> responseObserver) {
        Pair<Bitmap, Exception> pair = ContentResolverUtil.getThumbnail(
                this.group.getContext(),
                Uri.parse(request.getValue()),
                60,
                60);
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
            return;
        }

        // convert output stream to input stream
        Pair<InputStream, Exception> pair2 = BitmapUtil.bitmapToInputStream(pair.first, Bitmap.CompressFormat.PNG, 100);
        if (pair2.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair2.second));
            return;
        }
        DownloadStreamHandler handler = new DownloadStreamHandler(pair2.first);
        try {
            handler.sendLongPreData(pair2.first.available());
        } catch (IOException ignored) {
        }
        Exception ex = handler.handle(bytes -> {
            responseObserver.onNext(Bytes.newBuilder().setValue(bytes).build());
        });
        if (ex != null) {
            responseObserver.onError(ex);
            return;
        }
        responseObserver.onCompleted();
    }

    @Override
    public void deleteMediaFile(String request, StreamObserver<Empty> responseObserver) {
        MediaStoreUtil.deleteMediaFile(this.group.getContext(), Uri.parse(request.getValue()));
        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void downloadMediaFile(String request, StreamObserver<Bytes> responseObserver) {
        Pair<InputStream, Exception> pair = ContentResolverUtil.openFile(this.group.getContext(), Uri.parse(request.getValue()));
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
            return;
        }
        DownloadStreamHandler handler = new DownloadStreamHandler(pair.first);
        try {
            handler.sendLongPreData(pair.first.available());
        } catch (IOException ignored) {
        }
        Exception ex = handler.handle(bytes -> {
            responseObserver.onNext(Bytes.newBuilder().setValue(bytes).build());
        });
        if (ex != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ex));
            return;
        }
        try {
            pair.first.close();
        } catch (Exception ignored) {
        }
        responseObserver.onCompleted();
    }
}
