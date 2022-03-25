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
import android.app.ActivityManager;
import android.app.usage.StorageStatsManager;
import android.content.Context;
import android.content.res.Configuration;
import android.location.LocationManager;
import android.os.Build;
import android.os.SystemClock;
import android.os.storage.StorageManager;
import android.os.storage.StorageVolume;
import android.provider.Settings;
import android.telephony.TelephonyManager;
import android.util.Pair;

import com.joxrays.godroidsvr.singleton.GPUInfoSingleton;
import com.joxrays.godroidsvr.message.DeviceInfo;
import com.joxrays.godroidsvr.message.DisplayInfo;
import com.joxrays.godroidsvr.message.GPUInfo;
import com.joxrays.godroidsvr.message.MemoryInfo;
import com.joxrays.godroidsvr.message.StorageSpaceInfo;
import com.joxrays.godroidsvr.message.SystemInfo;

import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import java.util.regex.Pattern;

public final class DeviceUtil {

    private DeviceUtil() {

    }

    private static Configuration getConfiguration(Context context) {
        return context.getResources().getConfiguration();
    }

    public static int getWidth(Context context) {
        return context.getResources().getDisplayMetrics().widthPixels;
    }

    public static int getHeight(Context context) {
        return context.getResources().getDisplayMetrics().heightPixels;
    }

    public static float getDensity(Context context) {
        // 0.75 - ldpi 1.0 - mdpi
        // 1.5 - hdpi 2.0 - xhdpi
        // 3.0 - xxhdpi 4.0 - xxxhdpi
        // 3.0 - xxhdpi 4.0 - xxxhdpi
        return context.getResources().getDisplayMetrics().density;
    }

    public static int getDensityDpi(Context context) {
        return getConfiguration(context).densityDpi;
    }

    public static float getFontScale(Context context) {
        return getConfiguration(context).fontScale;
    }

    public static int getMcc(Context context) {
        return getConfiguration(context).mcc;
    }

    public static int getMnc(Context context) {
        return getConfiguration(context).mnc;
    }

    public static String getOrientation(Context context) {
        return getConfiguration(context).orientation == 1 ? "Portrait" : "Landscape";
    }

    // NOTOUCH=1 FINGER=3
    public static boolean getTouchScreen(Context context) {
        return getConfiguration(context).touchscreen == 3;
    }

    public static float getRefreshRate(Context context) {
        return CommonUtil.getWindowManager(context).getDefaultDisplay().getRefreshRate();
    }

    public static boolean isSupportHDR(Context context) {
        return CommonUtil.getWindowManager(context).getDefaultDisplay().isHdr();
    }

    public static String getHdrCapabilities(Context context) {
        final Map<Integer, String> hdrMap = new ConcurrentHashMap<>();
        hdrMap.put(1, "Dolby Vision");
        hdrMap.put(2, "HDR10");
        hdrMap.put(3, "HLG");
        hdrMap.put(4, "HDR10+");
        int[] types = CommonUtil.getWindowManager(context).getDefaultDisplay().getHdrCapabilities().getSupportedHdrTypes();
        if (types == null) return "";
        List<String> list = new ArrayList<>();
        for (int t : types) {
            if (hdrMap.containsKey(t)) {
                list.add(hdrMap.get(t));
            }
        }
        return String.join(",", list);
    }

    @SuppressLint("HardwareIds")
    public static String getAndroidId(Context context) {
        return Settings.Secure.getString(context.getContentResolver(), Settings.Secure.ANDROID_ID);
    }

    public static String getAbi() {
        if (Build.CPU_ABI != null) {
            return Build.CPU_ABI;
        }
        return "";
    }

    public static String getManufacturer() {
        return Build.MANUFACTURER;
    }

    public static String getProduct() {
        return Build.PRODUCT;
    }

    public static String getBrand() {
        return Build.BRAND;
    }

    public static String getBoard() {
        return Build.BOARD;
    }

    public static String getModel() {
        return Build.MODEL;
    }

    public static String getDeviceName() {
        return Build.DEVICE;
    }

    public static String getFingerPrint() {
        return Build.FINGERPRINT;
    }

    public static String getHardware() {
        return Build.HARDWARE;
    }

    public static String getHost() {
        return Build.HOST;
    }

    public static String getUser() {
        return Build.USER;
    }

    public static String getDisplay() {
        return Build.DISPLAY;
    }

    public static String getReleaseVersion() {
        return Build.VERSION.RELEASE;
    }

    public static int getSDK() {
        return Build.VERSION.SDK_INT;
    }

    public static long getBuildTime() {
        return Build.TIME;
    }

    public static String getDefaultLanguage() {
        return Locale.getDefault().getLanguage();
    }

    public static Pair<StorageSpaceInfo, Exception> getStorageSpaceSize(Context context) {
        StorageStatsManager ssm = CommonUtil.getStorageStatsManager(context);
        StorageManager sm = CommonUtil.getStorageManager(context);
        List<StorageVolume> list = sm.getStorageVolumes();
        StorageSpaceInfo.Builder builder = StorageSpaceInfo.newBuilder();
        for (StorageVolume sv : list) {
            final UUID uuid = sv.getUuid() != null ? UUID.fromString(sv.getUuid()) : StorageManager.UUID_DEFAULT;
            try {
                long total = ssm.getTotalBytes(uuid);
                long free = ssm.getFreeBytes(uuid);
                builder.setFreeSize(free)
                        .setUsedSize(total - free)
                        .setTotalSize(total);
                break;
            } catch (Exception ex) {
                return Pair.create(null, ex);
            }
        }
        return Pair.create(builder.build(), null);
    }

    public static MemoryInfo getMemoryInfo(Context context) {
        ActivityManager.MemoryInfo mi = new ActivityManager.MemoryInfo();
        ActivityManager activityManager = CommonUtil.getActivityManager(context);
        activityManager.getMemoryInfo(mi);
        return MemoryInfo.newBuilder()
                .setTotalMem(mi.totalMem)
                .setUsedMem(mi.totalMem - mi.availMem)
                .setAvailableMem(mi.availMem)
                .setThreshold(mi.threshold)
                .setLowMemory(mi.lowMemory)
                .build();
    }

    public static DeviceInfo getDeviceInfo(Context context) {
        return DeviceInfo.newBuilder().setManufacturer(getManufacturer())
                .setAndroidId(getAndroidId(context))
                .setProduct(getProduct())
                .setBrand(getBrand())
                .setBoard(getBoard())
                .setModel(getModel())
                .setDeviceName(getDeviceName())
                .setFingerprint(getFingerPrint())
                .setHardware(getHardware())
                .setRoot(isRoot())
                .setAdb(isEnableAdb(context))
                .setSimCard(isSimCardReady(context))
                .setDeveloper(isEnableDeveloperMode(context))
                .setAirplane(isEnableAirplaneMode(context))
                .setBluetooth(isEnableBlueTooth(context))
                .setLocation(isEnableWifiLocation(context))
                .setBuildTime(getBuildTime()).build();
    }

    public static SystemInfo getSystemInfo(Context context) {
        return SystemInfo.newBuilder().setHost(getHost())
                .setDisplay(getDisplay())
                .setUser(getUser())
                .setReleaseVersion(getReleaseVersion())
                .setSdk(getSDK())
                .setLanguage(getDefaultLanguage())
                .setAbi(getAbi())
                .setKernelRelease(getKernelReleaseNumber())
                .setKernelVersion(getKernelVersion())
                .setUptime(getSystemUptime())
                .setMcc(getMcc(context))
                .setMnc(getMnc(context))
                .build();
    }

    public static DisplayInfo getDisplayInfo(Context context) {
        int mode = 0;
        Pair<Integer, Exception> pair = ControllerUtil.getScreenBrightnessMode(context);
        if (pair.second == null) {
            mode = pair.first;
        }

        return DisplayInfo.newBuilder()
                .setHeight(getHeight(context))
                .setWidth(getWidth(context))
                .setDensity(getDensity(context))
                .setRefreshRate(getRefreshRate(context))
                .setScreenOffTime(getScreenOffTimeout(context))
                .setScreenBrightness(ControllerUtil.getScreenBrightness(context))
                .setScreenBrightnessMode(mode == 1 ? "Manual" : "Auto")
                .setSupportHdr(isSupportHDR(context))
                .setHdrCapabilities(getHdrCapabilities(context))
                .setFontScale(getFontScale(context))
                .setDensityDpi(getDensityDpi(context))
                .setOrientation(getOrientation(context))
                .setTouchScreen(getTouchScreen(context))
                .build();
    }

    private static int getSettingValue(Context context, String key) {
        try {
            return Settings.Secure.getInt(context.getContentResolver(), key);
        } catch (Exception ignored) {
            return 0;
        }
    }

    public static boolean isRoot() {
        for (String file : new String[]{"/sbin/su", "/system/bin/su", "/system/xbin/su",
                "/data/local/xbin/su", "/data/local/bin/su", "/system/sd/xbin/su",
                "/system/bin/failsafe/su", "/data/local/su"}) {
            if (FilesUtil.exist(new File(file))) {
                return true;
            }
        }
        return false;
    }

    public static boolean isEnableAdb(Context context) {
        return getSettingValue(context, Settings.Global.ADB_ENABLED) > 0;
    }

    public static boolean isEnableAirplaneMode(Context context) {
        return getSettingValue(context, Settings.Global.AIRPLANE_MODE_ON) > 0;
    }

    public static boolean isEnableBlueTooth(Context context) {
        return getSettingValue(context, Settings.Global.BLUETOOTH_ON) > 0;
    }

    public static boolean isEnableDeveloperMode(Context context) {
        return getSettingValue(context, Settings.Global.DEVELOPMENT_SETTINGS_ENABLED) > 0;
    }

    public static boolean isEnableWifiLocation(Context context) {
        return CommonUtil.getLocationManager(context).isProviderEnabled(LocationManager.NETWORK_PROVIDER);
    }

    public static boolean isSimCardReady(Context context) {
        return CommonUtil.getTelephonyManager(context).getSimState() == TelephonyManager.SIM_STATE_READY;
    }

    public static int getScreenOffTimeout(Context context) {
        try {
            return Settings.System.getInt(context.getContentResolver(), Settings.System.SCREEN_OFF_TIMEOUT);
        } catch (Exception ex) {
            return 0;
        }
    }

    public static String getKernelVersion() {
        Pair<String, Exception> pair = CommonUtil.executeAll("uname", "-v");
        return pair.first != null ? pair.first : "";
    }

    public static String getKernelReleaseNumber() {
        Pair<String, Exception> pair = CommonUtil.executeAll("uname", "-r");
        return pair.first != null ? pair.first : "";
    }

    public static Pair<List<Integer>, Exception> getCPUsFrequency() {
        // notice the directory "/sys/devices/system/cpu/" maybe not exist in some Android devices
        File file = new File("/sys/devices/system/cpu/");
        File[] cpuFileList = file.listFiles(pathname -> Pattern.matches("cpu[0-9]+", pathname.getName()));
        if (cpuFileList == null || cpuFileList.length == 0) {
            return null;
        }
        List<Integer> list = new ArrayList<>();
        for (File cpuFile : cpuFileList) {
            // notice, the directory "cpufreq" and "scaling_cur_freq" also maybe not exist in some Android devices
            Path freqPath = Paths.get(cpuFile.getAbsolutePath(), "cpufreq", "scaling_cur_freq");
            Pair<String, Exception> pair = FilesUtil.readText(freqPath.toFile());
            if (pair.second != null) {
                return Pair.create(null, pair.second);
            }
            list.add(Integer.parseInt(pair.first.trim()));
        }
        return Pair.create(list, null);
    }

    public static GPUInfo getGPUInfo() {
        return GPUInfoSingleton.getInstance().getGPUInfo();
    }

    public static long getSystemUptime() {
        return SystemClock.uptimeMillis();
    }
}
