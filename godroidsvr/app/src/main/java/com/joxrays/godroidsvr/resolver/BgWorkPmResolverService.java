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
import android.graphics.drawable.Drawable;
import android.os.Build;
import android.util.Pair;

import androidx.annotation.RequiresApi;

import com.joxrays.godroidsvr.message.PackageMetaInfoList;
import com.joxrays.godroidsvr.message.ParamBytes;
import com.joxrays.godroidsvr.message.Status;
import com.joxrays.godroidsvr.observer.DownloadStreamHandler;
import com.joxrays.godroidsvr.observer.UploadStreamObserver;
import com.joxrays.godroidsvr.util.BitmapUtil;
import com.joxrays.godroidsvr.util.PackageInstallUtil;
import com.joxrays.godroidsvr.util.PackageUtil;
import com.joxrays.godroidsvr.message.AppSize;
import com.joxrays.godroidsvr.message.ApplicationInfo;
import com.joxrays.godroidsvr.message.Bytes;
import com.joxrays.godroidsvr.message.PackageInfo;
import com.joxrays.godroidsvr.message.String;
import com.joxrays.godroidsvr.message.Empty;
import com.joxrays.godroidsvr.message.StringList;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;

import io.grpc.stub.StreamObserver;

public class BgWorkPmResolverService extends PmResolverGrpc.PmResolverImplBase {

    private final BgWorkBaseResolverGroup group;

    public BgWorkPmResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getApplicationInfo(String request, StreamObserver<ApplicationInfo> responseObserver) {
        Pair<ApplicationInfo, Exception> pair = PackageUtil.getApplicationInfo(this.group.getContext(), request.getValue());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getPackageInfo(String request, StreamObserver<PackageInfo> responseObserver) {
        Pair<PackageInfo, Exception> pair = PackageUtil.getPackageInfo(this.group.getContext(), request.getValue());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getAllPackageInfo(Empty request, StreamObserver<PackageMetaInfoList> responseObserver) {
        responseObserver.onNext(
                PackageMetaInfoList.newBuilder()
                        .addAllValues(PackageUtil.getAllPackageInfo(this.group.getContext()))
                        .build()
        );
        responseObserver.onCompleted();
    }

    @Override
    public void getAllUserPackageInfo(Empty request, StreamObserver<PackageMetaInfoList> responseObserver) {
        responseObserver.onNext(
                PackageMetaInfoList.newBuilder()
                        .addAllValues(PackageUtil.getAllUserPackageInfo(this.group.getContext()))
                        .build()
        );
        responseObserver.onCompleted();
    }

    @Override
    public void getAllSystemPackageInfo(Empty request, StreamObserver<PackageMetaInfoList> responseObserver) {
        responseObserver.onNext(
                PackageMetaInfoList.newBuilder()
                        .addAllValues(PackageUtil.getAllSystemPackageInfo(this.group.getContext()))
                        .build()
        );
        responseObserver.onCompleted();
    }

    @Override
    public void getApplicationSize(String request, StreamObserver<AppSize> responseObserver) {
        Pair<AppSize, Exception> pair = PackageUtil.getAppSize(this.group.getContext(), request.getValue());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @Override
    public StreamObserver<ParamBytes> installApk(StreamObserver<Status> responseObserver) {
        // install apk
        return new UploadStreamObserver(this.group.getContext(), responseObserver, UploadStreamObserver.Type.Install);
    }

    @RequiresApi(api = Build.VERSION_CODES.S)
    @Override
    public void uninstallApk(String request, StreamObserver<Empty> responseObserver) {
        if (!PackageUtil.checkPackageExist(this.group.getContext(), request.getValue())) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ErrorExceptionUtil.ErrorNotFoundPackage));
        } else {
            PackageInstallUtil.uninstallApk(this.group.getContext(), request.getValue());
            responseObserver.onNext(Empty.newBuilder().build());
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getApk(String request, StreamObserver<String> responseObserver) {
        Pair<ApplicationInfo, Exception> info = PackageUtil.getApplicationInfo(this.group.getContext(), request.getValue());
        if (info.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(info.second));
            return;
        }
        responseObserver.onNext(String.newBuilder().setValue(info.first.getSourceDir()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void getIcon(String request, StreamObserver<Bytes> responseObserver) {
        Pair<Drawable, Exception> pair = PackageUtil.getApplicationIcon(this.group.getContext(), request.getValue());
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
            return;
        }
        if (pair.first == null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ErrorExceptionUtil.ErrorDrawableIsNull));
            return;
        }
        // convert Drawable to InputStream
        Pair<InputStream, Exception> pair2 = BitmapUtil.bitmapToInputStream(
                BitmapUtil.drawableToBitmap(pair.first), Bitmap.CompressFormat.PNG, 100);

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
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ex));
            return;
        }
        responseObserver.onCompleted();
    }

    private void getArray(StreamObserver<StringList> responseObserver, Pair<List<java.lang.String>, Exception> pair) {
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
            return;
        }
        responseObserver.onNext(StringList.newBuilder()
                .addAllValues(pair.first)
                .build());
        responseObserver.onCompleted();
    }

    @Override
    public void getPermissions(String request, StreamObserver<StringList> responseObserver) {
        getArray(responseObserver, PackageUtil.getPackagePermissions(this.group.getContext(), request.getValue()));
    }

    @Override
    public void getActivities(String request, StreamObserver<StringList> responseObserver) {
        getArray(responseObserver, PackageUtil.getPackageActivities(this.group.getContext(), request.getValue()));
    }

    @Override
    public void getServices(String request, StreamObserver<StringList> responseObserver) {
        getArray(responseObserver, PackageUtil.getPackageServices(this.group.getContext(), request.getValue()));
    }

    @Override
    public void getReceivers(String request, StreamObserver<StringList> responseObserver) {
        getArray(responseObserver, PackageUtil.getPackageReceivers(this.group.getContext(), request.getValue()));
    }

    @Override
    public void getProviders(String request, StreamObserver<StringList> responseObserver) {
        getArray(responseObserver, PackageUtil.getPackageProvider(this.group.getContext(), request.getValue()));
    }

    @Override
    public void getSharedLibFiles(String request, StreamObserver<StringList> responseObserver) {
        getArray(responseObserver, PackageUtil.getPackageSharedLibFiles(this.group.getContext(), request.getValue()));
    }
}
