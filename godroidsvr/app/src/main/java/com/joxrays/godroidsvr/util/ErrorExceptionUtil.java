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

package com.joxrays.godroidsvr.util;

import java.io.FileNotFoundException;
import java.io.IOException;

import io.grpc.Status;

public final class ErrorExceptionUtil {
    public final static Exception ErrorFileLarge = new Exception("the file larger than 1M");
    public final static Exception ErrorNotFoundPackage = new Exception("not found package name");
    public final static Exception ErrorFileNotFound = new FileNotFoundException("file not found");
    public final static Exception ErrorCreateFile = new FileNotFoundException("create new file failed");
    public final static Exception ErrorDeleteFile = new FileNotFoundException("delete previous file failed");
    public final static Exception ErrorNotFoundWifiInfo = new Exception("not found wifi info");
    public final static Exception ErrorCannotListDir = new Exception("cannot list directory due to permission denied or not exist");
    public final static Exception ErrorCannotGetLocation = new Exception("can not get location");
    public final static Exception ErrorCannotGetInfoContentProvider = new Exception("can not get info by content provider");
    public final static Exception ErrorOperandFailed = new IOException("operand failed");
    public final static Exception ErrorStreamCallbackIsNull = new Exception("stream callback is null");
    public final static Exception ErrorDrawableIsNull = new Exception("drawable is null");
    public final static Exception ErrorExecuteTimeout = new Exception("execute timeout");
    public final static Exception ErrorParseJson = new Exception("parse json failed");

    public static Exception getRpcException(Exception ex) {
        return Status.INTERNAL.withCause(ex).withDescription(ex.getMessage()).asRuntimeException();
    }
}
