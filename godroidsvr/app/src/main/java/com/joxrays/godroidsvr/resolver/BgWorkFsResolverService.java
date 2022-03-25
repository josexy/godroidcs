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

import android.util.Pair;

import com.joxrays.godroidsvr.message.Bytes;
import com.joxrays.godroidsvr.message.FileInfoList;
import com.joxrays.godroidsvr.message.ParamBytes;
import com.joxrays.godroidsvr.message.Status;
import com.joxrays.godroidsvr.message.StringPair;
import com.joxrays.godroidsvr.message.StringTuple;
import com.joxrays.godroidsvr.observer.DownloadStreamHandler;
import com.joxrays.godroidsvr.observer.UploadStreamObserver;
import com.joxrays.godroidsvr.message.String;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;
import com.joxrays.godroidsvr.util.FilesUtil;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.nio.file.Path;

import io.grpc.stub.StreamObserver;

public class BgWorkFsResolverService extends FsResolverGrpc.FsResolverImplBase {

    private enum OperandType {
        CreateFile,
        DeleteFile,
        Mkdir,
        Rmdir,
        Move,
        Rename,
        Copy,
        ReadText,
        WriteText,
        AppendText,
    }

    private final BgWorkBaseResolverGroup group;

    public BgWorkFsResolverService(BgWorkBaseResolverGroup group) {
        this.group = group;
    }

    @Override
    public void getBaseFileTree(StringTuple request, StreamObserver<String> responseObserver) {
        if (!FilesUtil.exist(new File(request.getFirst()))) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ErrorExceptionUtil.ErrorFileNotFound));
            return;
        }
        try {
            // the max level can not more than 2
            java.lang.String json = FilesUtil.getBaseTreeFileInfoJson(new File(request.getFirst()),
                    request.getSecond(), request.getThird().equals("all"), 1);
            responseObserver.onNext(String.newBuilder().setValue(json).build());
            responseObserver.onCompleted();
        } catch (Exception ex) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ex));
        }
    }

    @Override
    public void listDir(StringPair request, StreamObserver<FileInfoList> responseObserver) {
        Pair<FileInfoList, Exception> pair = FilesUtil.listDir(request.getFirst(), request.getSecond().equals("all"));
        if (pair.second != null) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(pair.second));
        } else {
            responseObserver.onNext(pair.first);
            responseObserver.onCompleted();
        }
    }

    @Override
    public StreamObserver<ParamBytes> uploadGeneralFile(StreamObserver<Status> responseObserver) {
        return new UploadStreamObserver(this.group.getContext(), responseObserver, UploadStreamObserver.Type.Save);
    }

    @Override
    public void downloadGeneralFile(String request, StreamObserver<Bytes> responseObserver) {
        File file = new File(request.getValue());
        if (!file.exists()) {
            responseObserver.onError(ErrorExceptionUtil.getRpcException(ErrorExceptionUtil.ErrorFileNotFound));
        } else {
            InputStream in;
            try {
                in = new FileInputStream(file);
            } catch (Exception ex) {
                responseObserver.onError(ErrorExceptionUtil.getRpcException(ex));
                return;
            }
            DownloadStreamHandler handler = new DownloadStreamHandler(in);
            handler.sendLongPreData(file.length());
            Exception ex = handler.handle(bytes -> {
                responseObserver.onNext(Bytes.newBuilder().setValue(bytes).build());
            });
            if (ex != null) {
                responseObserver.onError(ErrorExceptionUtil.getRpcException(ex));
                return;
            }
            responseObserver.onCompleted();
        }
    }

    private void handleFileAndDirOp(OperandType type, java.lang.String first, java.lang.String second,
                                    StreamObserver<Status> responseObserver) {
        File file = new File(first);
        boolean status = false;
        Exception exception = null;
        java.lang.String text = "";
        switch (type) {
            case CreateFile:
                Pair<Boolean, Exception> pair1 = FilesUtil.createFile(file);
                status = pair1.first;
                exception = pair1.second;
                break;
            case DeleteFile:
                status = FilesUtil.deleteFile(file);
                break;
            case Mkdir:
                status = FilesUtil.mkDir(file);
                break;
            case Rmdir:
                Pair<Boolean, Long> pair2 = FilesUtil.rmDir(file);
                status = pair2.first;
                break;
            case Move:
                Pair<Path, Exception> pair3 = FilesUtil.move(file, new File(second));
                status = (pair3.second == null);
                exception = pair3.second;
                break;
            case Rename:
                status = FilesUtil.rename(file, new File(second));
                break;
            case Copy:
                pair3 = FilesUtil.copy(file, new File(second));
                status = (pair3.second == null);
                exception = pair3.second;
                break;
            case ReadText:
                Pair<java.lang.String, Exception> pair4 = FilesUtil.readText(file);
                text = pair4.first;
                status = (pair4.second == null);
                exception = pair4.second;
                break;
            case WriteText:
            case AppendText:
                boolean append = type == OperandType.AppendText;
                pair3 = FilesUtil.writeText(file, second, append);
                status = (pair3.second == null);
                exception = pair3.second;
                break;
            default:
                break;
        }

        // operate failed
        if (!status) {
            if (exception == null) {
                exception = ErrorExceptionUtil.ErrorOperandFailed;
            }
            responseObserver.onError(ErrorExceptionUtil.getRpcException(exception));
            return;
        }
        // operate succeed
        responseObserver.onNext(Status.newBuilder()
                .setStatus(Status.CODE.SUCCEED)
                .setMessage(text)
                .build());

        responseObserver.onCompleted();
    }

    @Override
    public void createFile(String request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.CreateFile, request.getValue(), null, responseObserver);
    }

    @Override
    public void deleteFile(String request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.DeleteFile, request.getValue(), null, responseObserver);
    }

    @Override
    public void mkDir(String request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.Mkdir, request.getValue(), null, responseObserver);
    }

    @Override
    public void rmDir(String request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.Rmdir, request.getValue(), null, responseObserver);
    }

    @Override
    public void move(StringPair request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.Move, request.getFirst(), request.getSecond(), responseObserver);
    }

    @Override
    public void rename(StringPair request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.Rename, request.getFirst(), request.getSecond(), responseObserver);

    }

    @Override
    public void copy(StringPair request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.Copy, request.getFirst(), request.getSecond(), responseObserver);
    }

    @Override
    public void readText(String request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.ReadText, request.getValue(), null, responseObserver);
    }

    @Override
    public void writeText(StringPair request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.WriteText, request.getFirst(), request.getSecond(), responseObserver);
    }

    @Override
    public void appendText(StringPair request, StreamObserver<Status> responseObserver) {
        handleFileAndDirOp(OperandType.AppendText, request.getFirst(), request.getSecond(), responseObserver);
    }
}


