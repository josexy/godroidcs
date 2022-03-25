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

package com.joxrays.godroidsvr.observer;

import com.google.protobuf.ByteString;
import com.joxrays.godroidsvr.util.ErrorExceptionUtil;

import java.io.InputStream;
import java.nio.ByteBuffer;
import java.nio.channels.Channels;
import java.nio.channels.ReadableByteChannel;

public class DownloadStreamHandler {
    private final static int BUFFER_SIZE = 4096;
    private final InputStream in;
    private final ByteBuffer buffer;
    private final ReadableByteChannel channel;

    private ByteString preSendData;
    private boolean preDataSent = false;
    private int sendState = 0;

    public DownloadStreamHandler(InputStream in) {
        this.in = in;
        this.channel = Channels.newChannel(in);
        this.buffer = ByteBuffer.allocate(BUFFER_SIZE);
    }

    /**
     * send extra parameter before sending bytes stream
     *
     * @param n
     */
    public void sendLongPreData(long n) {
        byte[] result = new byte[8];
        for (int i = 7; i >= 0; i--) {
            result[i] = (byte) (n & 0xFF);
            n >>= 8;
        }
        preSendData = ByteString.copyFrom(result);
    }

    private void sendMask(StreamCallback callback) {
        byte mask;
        if (!preDataSent && preSendData != null) {
            mask = (byte) 0xAA;
            sendState = 1;
        } else {
            mask = (byte) 0xBB;
            sendState = 2;
        }
        callback.write(ByteString.copyFrom(new byte[]{mask}));
    }

    public Exception handle(StreamCallback callback) {
        if (callback == null) return ErrorExceptionUtil.ErrorStreamCallbackIsNull;

        // format: | mask | [extra parameter] | bytes stream |
        try {
            completed:
            while (true) {
                switch (sendState) {
                    case 0: // mask byte
                        sendMask(callback);
                        break;
                    case 1: // extra parameter
                        callback.write(preSendData);
                        preDataSent = true;
                        sendState = 2;
                        break;
                    case 2: // main bytes stream
                        int n = channel.read(buffer);
                        if (n <= 0) {
                            break completed;
                        }
                        buffer.flip();
                        callback.write(ByteString.copyFrom(buffer));
                        buffer.clear();
                        break;
                    default:
                        break completed;
                }
            }
            in.close();
            channel.close();
            return null;
        } catch (Exception ex) {
            return ex;
        }
    }
}
