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
import android.content.Intent;
import android.content.IntentFilter;
import android.os.BatteryManager;

import com.joxrays.godroidsvr.message.BatteryInfo;

public final class BatteryUtil {

    private BatteryUtil() {

    }

    public static Intent getBatteryIntent(Context context) {
        IntentFilter intentFilter = new IntentFilter(Intent.ACTION_BATTERY_CHANGED);
        return context.registerReceiver(null, intentFilter);
    }

    public static int getStatus(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_STATUS, BatteryManager.BATTERY_STATUS_UNKNOWN);
    }

    public static String getStatusString(Intent intent) {
        int value = getStatus(intent);
        String status;
        switch (value) {
            case 2: // BATTERY_STATUS_CHARGING
                status = "Charging";
                break;
            case 3: // BATTERY_STATUS_DISCHARGING
                status = "Discharging";
                break;
            case 4: // BATTERY_STATUS_NOT_CHARGING
                status = "Not charging";
                break;
            case 5: // BATTERY_STATUS_FULL
                status = "Full";
                break;
            default:
                status = "Unknown";
                break;
        }
        return status;
    }

    public static int getHealth(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_HEALTH, BatteryManager.BATTERY_HEALTH_UNKNOWN);
    }

    public static String getHealthString(Intent intent) {
        int value = getHealth(intent);
        String health;
        switch (value) {
            case 2: // BATTERY_HEALTH_GOOD
                health = "Good";
                break;
            case 3: // BATTERY_HEALTH_OVERHEAT
                health = "OverHeat";
                break;
            case 4: // BATTERY_HEALTH_DEAD
                health = "Dead";
                break;
            case 5: // BATTERY_HEALTH_OVER_VOLTAGE
                health = "Over Voltage";
                break;
            case 6: // BATTERY_HEALTH_UNSPECIFIED_FAILURE
                health = "Unspecified Failure";
                break;
            case 7: // BATTERY_HEALTH_COLD
                health = "Cold";
                break;
            default:
                health = "Unknown";
                break;
        }
        return health;
    }

    public static boolean isPresent(Intent intent) {
        return intent.getBooleanExtra(BatteryManager.EXTRA_PRESENT, false);
    }

    public static int getPlugged(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_PLUGGED, 0);
    }

    // current level value of battery
    public static int getLevel(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_LEVEL, 0);
    }

    // the max level value of battery
    public static int getScale(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_SCALE, 0);
    }

    // the voltage of battery
    // for example: 5000mV
    public static int getVoltage(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_VOLTAGE, -1);
    }

    // the temperature of battery
    // for example: 25°C
    public static int getTemperature(Intent intent) {
        return intent.getIntExtra(BatteryManager.EXTRA_TEMPERATURE, -1);
    }

    public static String getTechnology(Intent intent) {
        return intent.getStringExtra(BatteryManager.EXTRA_TECHNOLOGY);
    }

    public static String getPluggedString(Intent intent) {
        int value = getPlugged(intent);
        String plugged;
        switch (value) {
            case BatteryManager.BATTERY_PLUGGED_AC:
                plugged = "AC";
                break;
            case BatteryManager.BATTERY_PLUGGED_USB:
                plugged = "USB";
                break;
            case BatteryManager.BATTERY_PLUGGED_WIRELESS:
                plugged = "Wireless";
                break;
            default:
                plugged = "Unknown";
        }
        return plugged;
    }

    public static BatteryInfo getBatteryInfo(Context context) {
        Intent intent = BatteryUtil.getBatteryIntent(context);
        return BatteryInfo.newBuilder().setHealth(BatteryUtil.getHealthString(intent))
                .setLevel(BatteryUtil.getLevel(intent))
                .setPresent(BatteryUtil.isPresent(intent))
                .setStatus(BatteryUtil.getStatusString(intent))
                .setScale(BatteryUtil.getScale(intent))
                .setVoltage(BatteryUtil.getVoltage(intent))
                .setTemperature(BatteryUtil.getTemperature(intent))
                .setTechnology(BatteryUtil.getTechnology(intent))
                .setPlugged(BatteryUtil.getPluggedString(intent)).build();
    }
}