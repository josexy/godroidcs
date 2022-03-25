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
import android.app.Service;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.hardware.display.DisplayManager;
import android.hardware.display.VirtualDisplay;
import android.media.MediaRecorder;
import android.media.projection.MediaProjection;
import android.os.Environment;
import android.os.IBinder;
import android.view.Surface;

import androidx.annotation.Nullable;
import androidx.core.app.NotificationCompat;

import com.joxrays.godroidsvr.singleton.FinalMediaProjectionSingleton;
import com.joxrays.godroidsvr.util.CommonUtil;
import com.joxrays.godroidsvr.util.DeviceUtil;
import com.joxrays.godroidsvr.util.DirectoryUtil;
import com.joxrays.godroidsvr.util.LogUtil;
import com.joxrays.godroidsvr.R;

import java.io.File;
import java.io.IOException;

public class MediaProjectionService extends Service {

    private MediaRecorder mediaRecorder;
    private MediaProjection projection;
    private VirtualDisplay virtualDisplay;
    private boolean isVideoMode;
    File saveVideoFile;

    private void initMediaRecorder(int width, int height) {
        mediaRecorder = new MediaRecorder();
        // mediaRecorder.setAudioSource(MediaRecorder.AudioSource.MIC);
        mediaRecorder.setVideoSource(MediaRecorder.VideoSource.SURFACE);
        mediaRecorder.setOutputFormat(MediaRecorder.OutputFormat.MPEG_4);
        mediaRecorder.setVideoEncoder(MediaRecorder.VideoEncoder.H264);
        // mediaRecorder.setAudioEncoder(MediaRecorder.AudioEncoder.AAC);
        mediaRecorder.setVideoSize(width, height);
        mediaRecorder.setVideoFrameRate(60);

        saveVideoFile = new File(DirectoryUtil.getExPublicDirTo(Environment.DIRECTORY_DCIM), CommonUtil.nowForFile() + ".mp4");
        LogUtil.d("save media projection video file: " + saveVideoFile);
        mediaRecorder.setOutputFile(saveVideoFile);
        mediaRecorder.setVideoEncodingBitRate(1024 * 1024 * 3);
        try {
            mediaRecorder.prepare();
        } catch (IOException e) {
            e.printStackTrace();
        }
        if (isVideoMode) {
            FinalMediaProjectionSingleton.getInstance().setImageFile(saveVideoFile.toString());
        }
    }

    private void createMediaProjectionNotification() {
        String channelId = getResources().getString(R.string.media_projection_channel_id);
        NotificationChannel nc = new NotificationChannel(channelId, "Channel" + channelId, NotificationManager.IMPORTANCE_DEFAULT);
        CommonUtil.getNotificationManager(this).createNotificationChannel(nc);

        Notification notification = new NotificationCompat.Builder(this, channelId)
                .setContentTitle(getResources().getString(R.string.app_name))
                .setContentText(getResources().getString(R.string.notification_content_text2))
                .setLargeIcon(BitmapFactory.decodeResource(getResources(), R.drawable.android))
                .setSmallIcon(R.drawable.android).build();

        // must call startForeground
        startForeground(Integer.parseInt(channelId), notification);
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        int width = DeviceUtil.getWidth(this);
        int height = DeviceUtil.getHeight(this);
        float density = DeviceUtil.getDensity(this);

        int resultCode = intent.getIntExtra("resultCode", -1);
        Intent data = intent.getParcelableExtra("data");
        Surface surface = intent.getParcelableExtra("surface");
        isVideoMode = surface == null;

        // create media projection notification
        createMediaProjectionNotification();
        LogUtil.d("start media projection");
        projection = CommonUtil.getMediaProjectionManager(this).getMediaProjection(resultCode, data);

        if (projection != null) {
            if (isVideoMode) {
                initMediaRecorder(width, height);
                surface = mediaRecorder.getSurface();
            }

            virtualDisplay = projection.createVirtualDisplay("Display0",
                    width, height, (int) density,
                    DisplayManager.VIRTUAL_DISPLAY_FLAG_PUBLIC,
                    surface, null, null);

            if (isVideoMode)
                mediaRecorder.start();
        }
        return super.onStartCommand(intent, flags, startId);
    }

    @Override
    public void onDestroy() {
        super.onDestroy();
        LogUtil.d("stop media projection");
        if (projection != null) {
            projection.stop();
        }
        if (mediaRecorder != null) {
            mediaRecorder.stop();
        }
        if (virtualDisplay != null) {
            virtualDisplay.release();
        }
    }

    @Nullable
    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }
}
