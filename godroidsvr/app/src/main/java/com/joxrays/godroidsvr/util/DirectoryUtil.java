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

import android.content.Context;
import android.os.Environment;

import java.io.File;

public final class DirectoryUtil {

    private DirectoryUtil() {

    }

    // /system
    public static File getRootPath() {
        return Environment.getRootDirectory();
    }

    // /data
    public static File getDataPath() {
        return Environment.getDataDirectory();
    }

    // /data/cache
    public static File getDownloadCachePath() {
        return Environment.getDownloadCacheDirectory();
    }

    // /storage/emulated/0
    public static File getExPublicDir() {
        return Environment.getExternalStorageDirectory();
    }

    // /storage/emulated/0/Download
    // /storage/emulated/0/Documents
    // ...
    public static File getExPublicDirTo(String type) {
        return Environment.getExternalStoragePublicDirectory(type);
    }

    // /storage/emulated/0/Android/data/<PackageName>/cache
    public static File getExPrivateCache(Context context) {
        return context.getExternalCacheDir();
    }

    // /storage/emulated/0/Android/data/<PackageName>/files
    public static File getExPrivateFiles(Context context) {
        return context.getExternalFilesDir(null);
    }

    // /storage/emulated/0/Android/data/<PackageName>/files/Download
    public static File getExPrivateFilesTo(Context context, String type) {
        return context.getExternalFilesDir(type);
    }

    public static File getInCachePath(Context context) {
        return context.getCacheDir();
    }

    public static File getInFilePath(Context context) {
        return context.getFilesDir();
    }

    public static File getInDataCachePath(Context context) {
        return context.getDataDir();
    }

    public static File getInCodeCachePath(Context context) {
        return context.getCodeCacheDir();
    }
}
