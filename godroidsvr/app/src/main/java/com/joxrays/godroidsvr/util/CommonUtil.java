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

import android.app.ActivityManager;
import android.app.AppOpsManager;
import android.app.NotificationManager;
import android.app.job.JobScheduler;
import android.app.usage.StorageStatsManager;
import android.app.usage.UsageStatsManager;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.pm.PackageManager;
import android.location.LocationManager;
import android.media.AudioManager;
import android.media.projection.MediaProjectionManager;
import android.net.ConnectivityManager;
import android.net.wifi.WifiManager;
import android.os.storage.StorageManager;
import android.telephony.PhoneNumberUtils;
import android.telephony.TelephonyManager;
import android.util.Pair;
import android.view.WindowManager;

import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.File;
import java.io.InputStreamReader;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Locale;

public final class CommonUtil {

    private CommonUtil() {
    }

    public static String formatPhoneNumber(Context context, String phone) {
        String countryIso = CommonUtil.getTelephonyManager(context).getNetworkCountryIso().toUpperCase();
        String v = PhoneNumberUtils.formatNumber(phone, countryIso);
        return v != null ? v : phone;
    }

    public static Exception executeRedirectFile(File file, String... command) {
        try {
            Process process = new ProcessBuilder(command).redirectOutput(file).start();
            return null;
        } catch (Exception ex) {
            return ex;
        }
    }

    public static Pair<List<String>, Exception> executeLines(String... command) {
        try {
            List<String> list = new ArrayList<>();
            Process process = new ProcessBuilder(command).start();
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    list.add(line);
                }
                return Pair.create(list, null);
            }
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static Pair<String, Exception> executeAll(String... command) {
        try {
            Process process = new ProcessBuilder(command).start();
            StringBuilder builder = new StringBuilder();
            try (BufferedInputStream in = new BufferedInputStream(process.getInputStream())) {
                byte[] buffer = new byte[4096];
                while (true) {
                    int n = in.read(buffer);
                    if (n <= 0) break;
                    builder.append(new String(buffer, 0, n));
                }
                return Pair.create(builder.toString().trim(), null);
            }
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    public static String dateTime(long timestamp) {
        return new SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault()).format(new Date(timestamp));
    }

    public static String now() {
        return dateTime(new Date().getTime());
    }

    public static String nowForFile() {
        return new SimpleDateFormat("yyyy_MM_dd_HH_mm_ss", Locale.getDefault()).format(new Date());
    }

    public static int getSignalLevel(int level) {
        return WifiManager.calculateSignalLevel(level, 5);
    }

    public static String getSignalStrengthFromLevel(int level) {
        return getSignalStrength(getSignalLevel(level));
    }

    public static String getSignalStrength(int signal) {
        switch (signal) {
            case 4:
                return "Excellent";
            case 3:
                return "Good";
            case 2:
                return "Low";
            case 1:
                return "Weak";
            default:
                return "Very Weak";
        }
    }

    public static String format(String format, Object... args) {
        return String.format(Locale.getDefault(), format, args);
    }

    public static String formatLine(String format, Object... args) {
        return CommonUtil.format(format + "\n", args);
    }

    public static PackageManager getPackageManager(Context context) {
        return context.getPackageManager();
    }

    public static UsageStatsManager getUsageStatsManager(Context context) {
        return (UsageStatsManager) context.getSystemService(Context.USAGE_STATS_SERVICE);
    }

    public static StorageStatsManager getStorageStatsManager(Context context) {
        return (StorageStatsManager) context.getSystemService(Context.STORAGE_STATS_SERVICE);
    }

    public static StorageManager getStorageManager(Context context) {
        return (StorageManager) context.getSystemService(Context.STORAGE_SERVICE);
    }

    public static ActivityManager getActivityManager(Context context) {
        return (ActivityManager) context.getSystemService(Context.ACTIVITY_SERVICE);
    }

    public static WifiManager getWifiManager(Context context) {
        return (WifiManager) context.getSystemService(Context.WIFI_SERVICE);
    }

    public static ConnectivityManager getConnectivityManager(Context context) {
        return (ConnectivityManager) context.getSystemService(Context.CONNECTIVITY_SERVICE);
    }

    public static ClipboardManager getClipboardManager(Context context) {
        return (ClipboardManager) context.getSystemService(Context.CLIPBOARD_SERVICE);
    }

    public static TelephonyManager getTelephonyManager(Context context) {
        return (TelephonyManager) context.getSystemService(Context.TELEPHONY_SERVICE);
    }

    public static JobScheduler getJobScheduler(Context context) {
        return (JobScheduler) context.getSystemService(Context.JOB_SCHEDULER_SERVICE);
    }

    public static AudioManager getAudioManager(Context context) {
        return (AudioManager) context.getSystemService(Context.AUDIO_SERVICE);
    }

    public static NotificationManager getNotificationManager(Context context) {
        return (NotificationManager) context.getSystemService(Context.NOTIFICATION_SERVICE);
    }

    public static MediaProjectionManager getMediaProjectionManager(Context context) {
        return (MediaProjectionManager) context.getSystemService(Context.MEDIA_PROJECTION_SERVICE);
    }

    public static LocationManager getLocationManager(Context context) {
        return (LocationManager) context.getSystemService(Context.LOCATION_SERVICE);
    }

    public static AppOpsManager getAppOpsManager(Context context) {
        return (AppOpsManager) context.getSystemService(Context.APP_OPS_SERVICE);
    }

    public static WindowManager getWindowManager(Context context) {
        return (WindowManager) context.getSystemService(Context.WINDOW_SERVICE);
    }
}
