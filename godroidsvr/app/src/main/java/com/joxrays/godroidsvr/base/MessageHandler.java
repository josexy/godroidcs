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

package com.joxrays.godroidsvr.base;

import android.graphics.Bitmap;
import android.os.Handler;
import android.os.Looper;
import android.os.Message;

import androidx.annotation.NonNull;

import com.joxrays.godroidsvr.singleton.WsServerSingleton;
import com.joxrays.godroidsvr.util.BitmapUtil;

import java.io.ByteArrayOutputStream;

public class MessageHandler extends Handler {

    public MessageHandler(@NonNull Looper looper) {
        super(looper);
    }

    @Override
    public void handleMessage(Message msg) {
        super.handleMessage(msg);
        Bitmap bitmap = (Bitmap) msg.obj;
        ByteArrayOutputStream out = new ByteArrayOutputStream(9102);
        // scale bitmap
        bitmap = BitmapUtil.scaleBitmap(bitmap, 0.52f, 0.52f);
        // compress bitmap
        bitmap.compress(Bitmap.CompressFormat.JPEG, 60, out);

        // send screen data to websocket client
        WsServerSingleton.getInstance().sendData(out.toByteArray());
    }
}