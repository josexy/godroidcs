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

import com.joxrays.godroidsvr.util.LogUtil;

import org.java_websocket.WebSocket;
import org.java_websocket.handshake.ClientHandshake;
import org.java_websocket.server.WebSocketServer;

import java.net.InetSocketAddress;
import java.util.concurrent.atomic.AtomicBoolean;

public class WsServer extends WebSocketServer {

    public AtomicBoolean ready = new AtomicBoolean(false);

    public WsServer(InetSocketAddress address) {
        super(address);
        setReuseAddr(true);
        setTcpNoDelay(true);
    }

    @Override
    public void onOpen(WebSocket conn, ClientHandshake handshake) {
        // notify websocket client to complete connection
        conn.send("OK");
        LogUtil.d("new connection " + conn.getRemoteSocketAddress());
        ready.set(true);
    }

    @Override
    public void onClose(WebSocket conn, int code, String reason, boolean remote) {
        LogUtil.d("closed " + conn.getRemoteSocketAddress() + " with exit code " + code + " additional info: " + reason);
        ready.set(false);
    }

    @Override
    public void onMessage(WebSocket conn, String message) {
    }

    @Override
    public void onError(WebSocket conn, Exception ex) {
    }

    @Override
    public void onStart() {
    }
}
