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

package com.joxrays.godroidsvr.service;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.os.IBinder;

import androidx.annotation.Nullable;
import androidx.core.app.NotificationCompat;

import com.joxrays.godroidsvr.util.CommonUtil;
import com.joxrays.godroidsvr.MainActivity;
import com.joxrays.godroidsvr.R;

public class MainNotificationService extends Service {

    private void createMainNotification() {
        String channelId = getResources().getString(R.string.main_channel_id);

        Intent main = new Intent(this, MainActivity.class);
        PendingIntent pi = PendingIntent.getActivity(this,
                0x1001, main,
                PendingIntent.FLAG_IMMUTABLE | PendingIntent.FLAG_UPDATE_CURRENT);

        Notification notification = new NotificationCompat.Builder(this, channelId)
                .setContentTitle(getResources().getString(R.string.app_name))
                .setContentText(getResources().getString(R.string.notification_content_text1))
                .setContentIntent(pi)
                .setLargeIcon(BitmapFactory.decodeResource(getResources(), R.drawable.android))
                .setSmallIcon(R.drawable.android).build();

        NotificationChannel nc = new NotificationChannel(channelId, "Channel" + channelId, NotificationManager.IMPORTANCE_DEFAULT);
        CommonUtil.getNotificationManager(this).createNotificationChannel(nc);
        // must call startForeground
        startForeground(Integer.parseInt(channelId), notification);
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        createMainNotification();
        return super.onStartCommand(intent, flags, startId);
    }

    @Nullable
    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }
}

