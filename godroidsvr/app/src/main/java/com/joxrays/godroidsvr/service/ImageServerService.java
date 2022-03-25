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

import android.annotation.SuppressLint;
import android.app.Service;
import android.content.Intent;
import android.graphics.PixelFormat;
import android.media.Image;
import android.media.ImageReader;
import android.os.Handler;
import android.os.HandlerThread;
import android.os.IBinder;
import android.os.Message;

import androidx.annotation.Nullable;

import com.joxrays.godroidsvr.MainActivity;
import com.joxrays.godroidsvr.base.MessageHandler;
import com.joxrays.godroidsvr.singleton.WsServerSingleton;
import com.joxrays.godroidsvr.util.BitmapUtil;

class WrapperHandler {
    HandlerThread handlerThread;
    Handler handler;

    public WrapperHandler(HandlerThread handlerThread, Handler handler) {
        this.handlerThread = handlerThread;
        this.handler = handler;
    }
}

public class ImageServerService extends Service implements ImageReader.OnImageAvailableListener {
    private static final String TAG = "ImageServerService";
    private ImageReader reader;
    private WrapperHandler wrapperHandler;

    @Override
    public void onImageAvailable(ImageReader reader) {
        Image image = reader.acquireLatestImage();
        if (image == null)
            return;

        Message message = Message.obtain(wrapperHandler.handler, 1, BitmapUtil.convertImageToBitmap(image));
        message.sendToTarget();

        image.close();
    }

    @Override
    public void onDestroy() {
        super.onDestroy();

        stopService(new Intent(this, MediaProjectionService.class));
        reader.setOnImageAvailableListener(null, null);
        wrapperHandler.handlerThread.quitSafely();

        WsServerSingleton.getInstance().stop();
    }

    @Nullable
    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    @SuppressLint("WrongConstant")
    @Override
    public int onStartCommand(@Nullable Intent intent, int flags, int startId) {
        reader = ImageReader.newInstance(
                this.getResources().getDisplayMetrics().widthPixels,
                this.getResources().getDisplayMetrics().heightPixels,
                PixelFormat.RGBA_8888, 5);

        reader.setOnImageAvailableListener(this, null);

        assert intent != null;

        int port = intent.getIntExtra("port", MainActivity.WS_SERVER_PORT);
        Intent x = new Intent(this, MediaProjectionService.class);
        x.putExtra("resultCode", intent.getIntExtra("resultCode", -1));
        x.putExtra("data", (Intent) intent.getParcelableExtra("data"));
        x.putExtra("surface", reader.getSurface());
        startForegroundService(x);

        // the sub-thread is used for sending screen image data to websocket client
        HandlerThread handlerThread = new HandlerThread(TAG);
        handlerThread.start();
        Handler myHandler = new MessageHandler(handlerThread.getLooper());
        wrapperHandler = new WrapperHandler(handlerThread, myHandler);

        // background websocket server
        new Thread(() -> WsServerSingleton.getInstance().startServer(port)).start();

        return super.onStartCommand(intent, flags, startId);
    }
}
