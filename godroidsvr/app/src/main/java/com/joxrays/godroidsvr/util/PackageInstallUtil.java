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

import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageInstaller;
import android.os.Build;

import androidx.annotation.RequiresApi;

import com.joxrays.godroidsvr.receiver.PackageInstallerReceiver;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;

public final class PackageInstallUtil {

    private PackageInstallUtil() {

    }

    // install apk from stream
    public static PackageInstaller.Session initPackageInstallerSession(Context context) throws IOException {
        PackageInstaller installer = CommonUtil.getPackageManager(context).getPackageInstaller();
        PackageInstaller.SessionParams params = new PackageInstaller.SessionParams(PackageInstaller.SessionParams.MODE_FULL_INSTALL);
        int sessionId = installer.createSession(params);
        return installer.openSession(sessionId);
    }

    public static OutputStream openSession(PackageInstaller.Session session) throws IOException {
        return session.openWrite("package", 0, -1);
    }

    @RequiresApi(api = Build.VERSION_CODES.S)
    public static void commitPackageInstallerSession(Context context, PackageInstaller.Session session) {
        Intent intent = new Intent(context, PackageInstallerReceiver.class);
        intent.setAction(PackageInstallerReceiver.PACKAGE_ACTION);

        PendingIntent pi = PendingIntent.getBroadcast(context, 0, intent,
                PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_MUTABLE);
        session.commit(pi.getIntentSender());
    }

    private static void addApkToInstallSession(String apkFile, PackageInstaller.Session session) throws
            IOException {
        try (OutputStream out = session.openWrite("package", 0, -1);
             InputStream in = new FileInputStream(apkFile)) {
            byte[] buffer = new byte[8192];
            int n;
            while ((n = in.read(buffer)) > 0) {
                out.write(buffer, 0, n);
            }
            out.flush();
        }
    }

    // install apk from file
    @RequiresApi(api = Build.VERSION_CODES.S)
    public static void installApk(Context context, File apkFile) {
        PackageInstaller.SessionParams params;
        PackageInstaller installer = CommonUtil.getPackageManager(context).getPackageInstaller();
        params = new PackageInstaller.SessionParams(PackageInstaller.SessionParams.MODE_FULL_INSTALL);
        try {
            int sessionId = installer.createSession(params);
            try (PackageInstaller.Session session = installer.openSession(sessionId)) {
                addApkToInstallSession(apkFile.getAbsolutePath(), session);

                Intent intent = new Intent(context, PackageInstallerReceiver.class);
                intent.setAction(PackageInstallerReceiver.PACKAGE_ACTION);

                PendingIntent pi = PendingIntent.getBroadcast(context, 0, intent,
                        PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_MUTABLE);
                session.commit(pi.getIntentSender());
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        } finally {
            FilesUtil.deleteFile(apkFile);
        }
    }

    @RequiresApi(api = Build.VERSION_CODES.S)
    public static void uninstallApk(Context context, String packageName) {
        Intent intent = new Intent(context, PackageInstallerReceiver.class);
        intent.setAction(PackageInstallerReceiver.PACKAGE_ACTION);
        PackageInstaller installer = context.getPackageManager().getPackageInstaller();
        PendingIntent pi = PendingIntent.getBroadcast(context, 0, intent, PendingIntent.FLAG_MUTABLE);
        installer.uninstall(packageName, pi.getIntentSender());
    }
}
