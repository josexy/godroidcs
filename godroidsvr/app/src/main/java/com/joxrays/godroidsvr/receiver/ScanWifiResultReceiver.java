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

package com.joxrays.godroidsvr.receiver;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.net.wifi.ScanResult;

import com.joxrays.godroidsvr.singleton.WifiResultSingleton;
import com.joxrays.godroidsvr.util.LogUtil;
import com.joxrays.godroidsvr.message.SimpleWifiInfo;
import com.joxrays.godroidsvr.util.CommonUtil;

import java.util.ArrayList;
import java.util.List;

public class ScanWifiResultReceiver extends BroadcastReceiver {
    private final List<SimpleWifiInfo> list = new ArrayList<>();

    private void buildScanWifiInfo(List<ScanResult> results) {
        if (results != null) {
            if (results.size() > 0) {
                list.clear();
            }
            for (ScanResult result : results) {
                LogUtil.d("ScanResult: " + result.SSID + "\t" + result.BSSID);

                list.add(SimpleWifiInfo.newBuilder()
                        .setSsid(result.SSID)
                        .setBssid(result.BSSID)
                        .setFrequency(result.frequency)
                        .setSignal(CommonUtil.getSignalStrengthFromLevel(result.level))
                        .build());
            }
            WifiResultSingleton.getInstance().set(list);
        }
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        List<ScanResult> results = CommonUtil.getWifiManager(context).getScanResults();
        buildScanWifiInfo(results);
    }
}
