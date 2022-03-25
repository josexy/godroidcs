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
import android.media.AudioManager;
import android.provider.Settings;
import android.util.Pair;

public final class ControllerUtil {

    private ControllerUtil() {

    }

    public static int getScreenBrightness(Context context) {
        return Settings.System.getInt(context.getContentResolver(),
                Settings.System.SCREEN_BRIGHTNESS, 125);
    }

    public static Exception setScreenBrightness(Context context, int value) {
        if (value < 0) value = 0;
        else if (value > 255) value = 255;
        try {
            Settings.System.putInt(context.getContentResolver(), Settings.System.SCREEN_BRIGHTNESS, value);
            return null;
        } catch (Exception ex) {
            return ex;
        }
    }

    // get screen brightness mode: (manual/auto)
    public static Pair<Integer, Exception> getScreenBrightnessMode(Context context) {
        try {
            int value = Settings.System.getInt(context.getContentResolver(),
                    Settings.System.SCREEN_BRIGHTNESS_MODE);
            return Pair.create(value, null);
        } catch (Exception ex) {
            return Pair.create(null, ex);
        }
    }

    // set screen brightness mode: (0:manual 1:auto)
    public static void setScreenBrightnessMode(Context context, int mode) {
        int value;
        if (mode == 0) value = Settings.System.SCREEN_BRIGHTNESS_MODE_MANUAL;
        else value = Settings.System.SCREEN_BRIGHTNESS_MODE_AUTOMATIC;
        Settings.System.putInt(context.getContentResolver(), Settings.System.SCREEN_BRIGHTNESS_MODE, value);
    }

    public static int getAudioMusicMaxValue(Context context) {
        return CommonUtil.getAudioManager(context).getStreamMaxVolume(AudioManager.STREAM_MUSIC);
    }

    public static int getAudioMusicValue(Context context) {
        return CommonUtil.getAudioManager(context).getStreamVolume(AudioManager.STREAM_MUSIC);
    }

    public static void setAudioMusicValue(Context context, int value, boolean ui) {
        final int max = getAudioMusicMaxValue(context);
        CommonUtil.getAudioManager(context).setStreamVolume(
                AudioManager.STREAM_MUSIC, value % max, ui ? AudioManager.FLAG_SHOW_UI : 0);
    }

    public static void increaseAudioMusic(Context context, boolean ui) {
        CommonUtil.getAudioManager(context).adjustStreamVolume(
                AudioManager.STREAM_MUSIC, AudioManager.ADJUST_RAISE, ui ? AudioManager.FLAG_SHOW_UI : 0);
    }

    public static void decreaseAudioMusic(Context context, boolean ui) {
        CommonUtil.getAudioManager(context).adjustStreamVolume(
                AudioManager.STREAM_MUSIC, AudioManager.ADJUST_LOWER, ui ? AudioManager.FLAG_SHOW_UI : 0);
    }
}
