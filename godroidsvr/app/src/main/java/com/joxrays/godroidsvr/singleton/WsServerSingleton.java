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

package com.joxrays.godroidsvr.singleton;

import com.joxrays.godroidsvr.base.WsServer;

import java.net.InetSocketAddress;

public class WsServerSingleton extends BaseSingleton {

    private WsServer server;

    protected WsServerSingleton() {

    }

    public static WsServerSingleton getInstance() {
        return (WsServerSingleton) BaseSingleton.getInstance(WsServerSingleton.class);
    }

    public void startServer(int port) {
        server = new WsServer(new InetSocketAddress(port));
        server.run();
    }

    public void stop() {
        if (server != null) {
            try {
                server.stop();
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
    }

    public void sendData(byte[] data) {
        // the websocket server was closed
        if (server == null)
            return;

        // there are no websocket clients connected to server
        if (!server.ready.get())
            return;

        server.broadcast(data);
    }
}
