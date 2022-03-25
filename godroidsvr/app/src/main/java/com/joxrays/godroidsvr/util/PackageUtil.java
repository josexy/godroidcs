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

import android.annotation.SuppressLint;
import android.app.usage.StorageStats;
import android.app.usage.StorageStatsManager;
import android.content.Context;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.content.pm.ProviderInfo;
import android.content.pm.ServiceInfo;
import android.graphics.drawable.Drawable;
import android.os.storage.StorageManager;
import android.os.storage.StorageVolume;
import android.util.Pair;

import com.joxrays.godroidsvr.message.AppSize;
import com.joxrays.godroidsvr.message.ApplicationInfo;
import com.joxrays.godroidsvr.message.PackageInfo;
import com.joxrays.godroidsvr.message.PackageMetaInfo;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.UUID;

public final class PackageUtil {

    private PackageUtil() {

    }

    public static boolean checkPackageExist(Context context, String packageName) {
        PackageManager pm = context.getPackageManager();
        try {
            pm.getPackageInfo(packageName, 0);
            return true;
        } catch (PackageManager.NameNotFoundException ignored) {
        }
        return false;
    }

    public static Pair<ApplicationInfo, Exception> getApplicationInfo(Context context, String packageName) {
        ApplicationInfo nai;
        PackageManager pm = context.getPackageManager();
        try {
            android.content.pm.ApplicationInfo ai = pm.getApplicationInfo(packageName, 0);

            String appName;
            try {
                appName = pm.getApplicationLabel(ai).toString();
            } catch (Exception ex) {
                return Pair.create(null, ex);
            }
            nai = ApplicationInfo.newBuilder()
                    .setSourceDir(ai.publicSourceDir)
                    .setProcessName(ai.processName)
                    .setDataDir(ai.dataDir)
                    .setAppName(appName)
                    .setSystemApp(isSystemApp(ai))
                    .setMinSdkVersion(ai.minSdkVersion)
                    .setTargetSdkVersion(ai.targetSdkVersion)
                    .build();
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }

        return Pair.create(nai, null);
    }

    private static boolean isSystemApp(android.content.pm.ApplicationInfo ai) {
        return (ai.flags & android.content.pm.ApplicationInfo.FLAG_SYSTEM) > 0
                || (ai.flags & android.content.pm.ApplicationInfo.FLAG_UPDATED_SYSTEM_APP) > 0;
    }

    public static Pair<PackageInfo, Exception> getApkPackageInfo(Context context, String apkFileName) {
        PackageManager pm = context.getPackageManager();
        try {
            return buildFromPackageInfo(context, pm.getPackageArchiveInfo(apkFileName, 0));
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    private static Pair<List<String>, Exception> getPackageListInfo(Context context, String packageName, int type) {
        PackageManager pm = context.getPackageManager();
        try {
            android.content.pm.PackageInfo pi = pm.getPackageInfo(packageName, type);
            List<String> list = new ArrayList<>();

            switch (type) {
                case PackageManager.GET_ACTIVITIES:
                    if (pi.activities != null) {
                        for (ActivityInfo ai : pi.activities) {
                            list.add(ai.name);
                        }
                    }
                    break;
                case PackageManager.GET_SERVICES:
                    if (pi.services != null) {
                        for (ServiceInfo si : pi.services) {
                            list.add(si.name);
                        }
                    }
                    break;
                case PackageManager.GET_RECEIVERS:
                    if (pi.receivers != null) {
                        for (ActivityInfo ai : pi.receivers) {
                            list.add(ai.name);
                        }
                    }
                    break;
                case PackageManager.GET_PROVIDERS:
                    if (pi.providers != null) {
                        for (ProviderInfo pri : pi.providers) {
                            list.add(pri.name);
                        }
                    }
                    break;
                case PackageManager.GET_PERMISSIONS:
                    if (pi.requestedPermissions != null) {
                        Collections.addAll(list, pi.requestedPermissions);
                    }
                    break;
                default:
                    if ((type & PackageManager.GET_SHARED_LIBRARY_FILES) != 0) {
                        if (pi.applicationInfo != null && pi.applicationInfo.sharedLibraryFiles != null) {
                            Collections.addAll(list, pi.applicationInfo.sharedLibraryFiles);
                        }
                    }
                    break;
            }

            return Pair.create(list, null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static Pair<List<String>, Exception> getPackageActivities(Context context, String packageName) {
        return getPackageListInfo(context, packageName, PackageManager.GET_ACTIVITIES);
    }

    public static Pair<List<String>, Exception> getPackageServices(Context context, String packageName) {
        return getPackageListInfo(context, packageName, PackageManager.GET_SERVICES);
    }

    public static Pair<List<String>, Exception> getPackageReceivers(Context context, String packageName) {
        return getPackageListInfo(context, packageName, PackageManager.GET_RECEIVERS);
    }

    public static Pair<List<String>, Exception> getPackageProvider(Context context, String packageName) {
        return getPackageListInfo(context, packageName, PackageManager.GET_PROVIDERS);
    }

    public static Pair<List<String>, Exception> getPackagePermissions(Context context, String packageName) {
        return getPackageListInfo(context, packageName, PackageManager.GET_PERMISSIONS);
    }

    public static Pair<List<String>, Exception> getPackageSharedLibFiles(Context context, String packageName) {
        return getPackageListInfo(context, packageName, PackageManager.GET_SHARED_LIBRARY_FILES);
    }

    public static Pair<PackageInfo, Exception> buildFromPackageInfo(Context context, android.content.pm.PackageInfo packageInfo) {
        Pair<ApplicationInfo, Exception> pair = getApplicationInfo(context, packageInfo.packageName);
        if (pair.second != null) {
            return Pair.create(null, pair.second);
        }
        String installer = context.getPackageManager().getInstallerPackageName(packageInfo.packageName);
        return Pair.create(
                PackageInfo.newBuilder()
                        .setPackageName(packageInfo.packageName)
                        .setVersionName(packageInfo.versionName)
                        .setFirstInstallTime(packageInfo.firstInstallTime)
                        .setLastUpdatedTime(packageInfo.lastUpdateTime)
                        .setApplicationInfo(pair.first)
                        .setInstaller(installer != null ? installer : "")
                        .build(),
                null
        );
    }

    public static Pair<PackageInfo, Exception> getPackageInfo(Context context, String packageName) {
        PackageManager pm = context.getPackageManager();
        try {
            return buildFromPackageInfo(context, pm.getPackageInfo(packageName, 0));
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static Pair<Drawable, Exception> getApplicationIcon(Context context, String packageName) {
        PackageManager pm = context.getPackageManager();
        try {
            android.content.pm.ApplicationInfo ai = pm.getApplicationInfo(packageName, 0);
            Drawable icon = ai.loadIcon(pm);
            return Pair.create(icon, null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static List<PackageMetaInfo> getAllPackageInfo(Context context) {
        return getAllPackageInfoByType(context, 0);
    }

    public static List<PackageMetaInfo> getAllUserPackageInfo(Context context) {
        return getAllPackageInfoByType(context, 1);
    }

    public static List<PackageMetaInfo> getAllSystemPackageInfo(Context context) {
        return getAllPackageInfoByType(context, 2);
    }

    public static List<PackageMetaInfo> getAllPackageInfoByType(Context context, int type) {
        List<PackageMetaInfo> list = new ArrayList<>();
        PackageManager pm = context.getPackageManager();

        @SuppressLint("QueryPermissionsNeeded")
        List<android.content.pm.PackageInfo> pis = pm.getInstalledPackages(PackageManager.GET_ACTIVITIES);

        for (android.content.pm.PackageInfo pi : pis) {
            // user
            if (type == 1) {
                if (isSystemApp(pi.applicationInfo)) continue;
            } else if (type == 2) {
                // system
                if (!isSystemApp(pi.applicationInfo)) continue;
            }
            PackageMetaInfo info = PackageMetaInfo.newBuilder()
                    .setPackageName(pi.packageName)
                    .setVersionName(pi.versionName)
                    .setAppName(pm.getApplicationLabel(pi.applicationInfo).toString())
                    .setSystemApp(isSystemApp(pi.applicationInfo))
                    .build();
            list.add(info);
        }
        return list;
    }

    public static Pair<Integer, Exception> getPackageUID(Context context, String packageName) {
        try {
            android.content.pm.PackageManager pm = context.getPackageManager();
            android.content.pm.ApplicationInfo ai = pm.getApplicationInfo(packageName, PackageManager.GET_META_DATA);
            return Pair.create(ai.uid, null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    // the higher Android devices maybe need "android.permission.PACKAGE_USAGE_STATS" permission
    public static Pair<AppSize, Exception> getAppSize(Context context, String packageName) {
        AppSize.Builder builder = AppSize.newBuilder();
        StorageStatsManager ssm = CommonUtil.getStorageStatsManager(context);
        StorageManager sm = CommonUtil.getStorageManager(context);
        List<StorageVolume> list = sm.getStorageVolumes();

        for (StorageVolume sv : list) {
            final UUID uuid = sv.getUuid() != null ? UUID.fromString(sv.getUuid()) : StorageManager.UUID_DEFAULT;
            try {
                Pair<Integer, Exception> pair = getPackageUID(context, packageName);
                if (pair.second != null) {
                    return Pair.create(null, pair.second);
                }
                StorageStats ss = ssm.queryStatsForUid(uuid, pair.first);
                long app = ss.getAppBytes();
                long cache = ss.getCacheBytes();
                long data = ss.getDataBytes();
                builder.setAppBytes(app)
                        .setCacheBytes(cache)
                        .setDataBytes(data)
                        .setTotalBytes(app + cache + data);
            } catch (Exception ex) {
                return Pair.create(null, ex);
            }
        }
        return Pair.create(builder.build(), null);
    }
}
