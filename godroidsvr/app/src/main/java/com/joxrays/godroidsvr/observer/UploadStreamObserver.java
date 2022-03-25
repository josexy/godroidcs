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

package com.joxrays.godroidsvr.observer;

import android.content.Context;
import android.content.pm.PackageInstaller;
import android.os.Build;

import androidx.annotation.RequiresApi;

import com.joxrays.godroidsvr.message.ParamBytes;
import com.joxrays.godroidsvr.message.Status;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;
import com.joxrays.godroidsvr.util.FilesUtil;
import com.joxrays.godroidsvr.util.LogUtil;
import com.joxrays.godroidsvr.util.PackageInstallUtil;

import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;
import java.nio.ByteBuffer;

import io.grpc.stub.StreamObserver;

public class UploadStreamObserver implements StreamObserver<ParamBytes> {

    // show that the upload file is apk or general file
    // if the type of file is apk, then install it.
    // Otherwise save it simply
    public enum Type {
        Save,
        Install,
    }

    private final Context context;
    private final StreamObserver<Status> responseObserver;
    private OutputStream out;
    private Exception exception;
    private final Type type;
    private String upload_file_path;
    private PackageInstaller.Session session;

    public UploadStreamObserver(Context context, StreamObserver<Status> streamObserver, Type type) {
        this.context = context;
        this.responseObserver = streamObserver;
        this.type = type;
    }

    @Override
    public void onNext(ParamBytes value) {
        // firstly save the upload file path
        if (upload_file_path == null && value.getParam() != null && !value.getParam().getValue().isEmpty()) {
            upload_file_path = value.getParam().getValue();
            LogUtil.d("upload file path: " + upload_file_path);
        }
        ByteBuffer buffer = ByteBuffer.wrap(value.getValue().getValue().toByteArray());

        if (out == null) {
            try {
                // save file
                if (type == Type.Save) {
                    File file = new File(upload_file_path);
                    if (FilesUtil.exist(file)) {
                        if (!FilesUtil.deleteFile(file)) {
                            exception = ErrorExceptionUtil.ErrorDeleteFile;
                            return;
                        }
                    }
                    if (!file.createNewFile()) {
                        exception = ErrorExceptionUtil.ErrorCreateFile;
                        return;
                    }
                    out = new FileOutputStream(file);
                } else {
                    // install apk
                    session = PackageInstallUtil.initPackageInstallerSession(context);
                    out = PackageInstallUtil.openSession(session);
                }
            } catch (Exception ex) {
                exception = ex;
                return;
            }
        }
        try {
            out.write(buffer.array(), buffer.position(), buffer.limit());
        } catch (Exception ex) {
            exception = ex;
        }
    }

    @Override
    public void onError(Throwable t) {
    }

    @RequiresApi(api = Build.VERSION_CODES.S)
    @Override
    public void onCompleted() {
        if (exception != null) {
            exception.printStackTrace();
        }

        responseObserver.onNext(Status.newBuilder()
                .setStatus(exception == null ? Status.CODE.SUCCEED : Status.CODE.FAILED)
                .setMessage(exception != null ? exception.getMessage() : "")
                .build());

        responseObserver.onCompleted();

        try {
            // flush and close output stream
            out.flush();
            out.close();
            if (exception == null && type == Type.Install) {
                PackageInstallUtil.commitPackageInstallerSession(context, session);
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }
}