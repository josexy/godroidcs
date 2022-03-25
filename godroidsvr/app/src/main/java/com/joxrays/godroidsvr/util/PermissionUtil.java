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

import android.Manifest;
import android.app.Activity;
import android.app.AppOpsManager;
import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.net.Uri;
import android.os.Build;
import android.os.Environment;
import android.provider.Settings;

import androidx.annotation.RequiresApi;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import java.util.ArrayList;
import java.util.List;

public final class PermissionUtil {

    private PermissionUtil() {

    }

    private static boolean checkPermission(Context context, String permission) {
        return PackageManager.PERMISSION_GRANTED == ContextCompat.checkSelfPermission(context, permission);
    }

    public static boolean checkReadPhoneStatePermission(Context context) {
        return checkPermission(context, Manifest.permission.READ_PHONE_STATE);
    }

    public static boolean checkAccessLocation(Context context) {
        return checkPermission(context, Manifest.permission.ACCESS_FINE_LOCATION);
    }

    public static boolean checkUsageStatsPermission(Context context) {
        boolean is;
        AppOpsManager appOps = CommonUtil.getAppOpsManager(context);
        int mode = appOps.unsafeCheckOp(AppOpsManager.OPSTR_GET_USAGE_STATS, android.os.Process.myUid(), context.getPackageName());
        if (mode == AppOpsManager.MODE_DEFAULT) {
            is = (context.checkCallingOrSelfPermission(Manifest.permission.PACKAGE_USAGE_STATS) == PackageManager.PERMISSION_GRANTED);
        } else {
            is = (mode == AppOpsManager.MODE_ALLOWED);
        }
        return is;
    }

    public static boolean checkStorageManagePermission(Context context) {
        int result1 = ContextCompat.checkSelfPermission(context, Manifest.permission.READ_EXTERNAL_STORAGE);
        int result2 = ContextCompat.checkSelfPermission(context, Manifest.permission.WRITE_EXTERNAL_STORAGE);
        return result1 == PackageManager.PERMISSION_GRANTED && result2 == PackageManager.PERMISSION_GRANTED;
    }

    public static boolean checkInstallPermission(Context context) {
        return context.getPackageManager().canRequestPackageInstalls();
    }

    public static boolean checkWriteSettingsPermission(Context context) {
        return Settings.System.canWrite(context);
    }

    public static boolean checkReadContactPermission(Context context) {
        return checkPermission(context, Manifest.permission.READ_CONTACTS);
    }

    public static boolean checkWriteContactPermission(Context context) {
        return checkPermission(context, Manifest.permission.WRITE_CONTACTS);
    }

    public static boolean checkReadSmsPermission(Context context) {
        return checkPermission(context, Manifest.permission.READ_SMS);
    }

    public static boolean checkSendSmsPermission(Context context) {
        return checkPermission(context, Manifest.permission.SEND_SMS);
    }

    public static boolean checkReadCallLogPermission(Context context) {
        return checkPermission(context, Manifest.permission.READ_CALL_LOG);
    }

    public static boolean checkWriteCallLogPermission(Context context) {
        return checkPermission(context, Manifest.permission.WRITE_CALL_LOG);
    }

    public static boolean checkCallPhonePermission(Context context) {
        return checkPermission(context, Manifest.permission.CALL_PHONE);
    }

    public static void requireUsageStatsPermission(Context context) {
        Intent intent = new Intent(Settings.ACTION_USAGE_ACCESS_SETTINGS);
        context.startActivity(intent);
    }

    @RequiresApi(api = Build.VERSION_CODES.R)
    public static void requireStorageManagePermissionAndroidR(Activity activity) {
        try {
            Intent intent = new Intent(Settings.ACTION_MANAGE_ALL_FILES_ACCESS_PERMISSION);
            intent.setData(Uri.parse("package:" + activity.getPackageName()));
            activity.startActivity(intent);
        } catch (Exception ex) {
            Intent intent = new Intent(Settings.ACTION_MANAGE_ALL_FILES_ACCESS_PERMISSION);
            activity.startActivity(intent);
        }
    }

    public static void requireInstallPermission(Context context) {
        Intent intent = new Intent(Settings.ACTION_MANAGE_UNKNOWN_APP_SOURCES, Uri.parse("package:" + context.getPackageName()));
        context.startActivity(intent);
    }


    public static void requireWriteSettingsPermission(Context context) {
        Intent intent = new Intent(Settings.ACTION_MANAGE_WRITE_SETTINGS);
        intent.setData(Uri.parse("package:" + context.getPackageName()));
        context.startActivity(intent);
    }

    public static void requireNeedPermissions(Activity activity) {
        List<String> permissions = new ArrayList<>();
        if (!checkAccessLocation(activity)) {
            permissions.add(Manifest.permission.ACCESS_FINE_LOCATION);
        }
        if (!checkReadContactPermission(activity)) {
            permissions.add(Manifest.permission.READ_CONTACTS);
        }
        if (!checkWriteContactPermission(activity)) {
            permissions.add(Manifest.permission.WRITE_CONTACTS);
        }
        if (!checkReadSmsPermission(activity)) {
            permissions.add(Manifest.permission.READ_SMS);
        }
        if (!checkSendSmsPermission(activity)) {
            permissions.add(Manifest.permission.SEND_SMS);
        }
        if (!checkCallPhonePermission(activity)) {
            permissions.add(Manifest.permission.CALL_PHONE);
        }
        if (!checkReadCallLogPermission(activity)) {
            permissions.add(Manifest.permission.READ_CALL_LOG);
        }
        if (!checkWriteCallLogPermission(activity)) {
            permissions.add(Manifest.permission.WRITE_CALL_LOG);
        }
        if (!checkReadPhoneStatePermission(activity)) {
            permissions.add(Manifest.permission.READ_PHONE_STATE);
        }
        // < Android 11
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.R) {
            if (!checkStorageManagePermission(activity))
                permissions.add(Manifest.permission.WRITE_EXTERNAL_STORAGE);
        }

        String[] perms = new String[permissions.size()];
        perms = permissions.toArray(perms);
        if (perms.length > 0)
            ActivityCompat.requestPermissions(activity, perms, 0x200);

        if (!checkWriteSettingsPermission(activity))
            requireWriteSettingsPermission(activity);
        if (!checkInstallPermission(activity))
            requireInstallPermission(activity);
        if (!checkUsageStatsPermission(activity))
            requireUsageStatsPermission(activity);

        // >= Android 11
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R && !Environment.isExternalStorageManager())
            requireStorageManagePermissionAndroidR(activity);
    }
}
