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
import android.content.pm.PackageInstaller;
import android.os.Bundle;

import com.joxrays.godroidsvr.util.LogUtil;

public class PackageInstallerReceiver extends BroadcastReceiver {
    public static final String PACKAGE_ACTION = "PackageInstallerReceiver.SESSION_API_PACKAGE";

    @Override
    public void onReceive(Context context, Intent intent) {
        Bundle extras = intent.getExtras();
        if (PACKAGE_ACTION.equals(intent.getAction())) {
            int status = extras.getInt(PackageInstaller.EXTRA_STATUS);

            switch (status) {
                case PackageInstaller.STATUS_PENDING_USER_ACTION:
                    LogUtil.d("PackageInstaller.STATUS_PENDING_USER_ACTION");
                    Intent confirmIntent = (Intent) extras.get(Intent.EXTRA_INTENT);
                    confirmIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
                    context.startActivity(confirmIntent);
                    break;
                case PackageInstaller.STATUS_SUCCESS:
                    LogUtil.d("PackageInstaller.STATUS_SUCCESS");
                    break;
                default:
                    String message = extras.getString(PackageInstaller.EXTRA_STATUS_MESSAGE);
                    LogUtil.d("failed to install: " + message);
                    break;
            }
        }
    }
}