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

package ws

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/josexy/godroidcli/android/observer"
	"github.com/josexy/godroidcli/util"
)

type ScreenCaptureClient struct {
	addr         string
	conn         *websocket.Conn
	msgObservers []observer.MessageObserver
	imgChan      chan image.Image
	limiter      Limiter
	closeChan    chan struct{}
}

func NewWsScreenCaptureClient(addr string) *ScreenCaptureClient {
	return &ScreenCaptureClient{
		addr:      addr,
		imgChan:   make(chan image.Image, 1024),
		limiter:   NewSimpleLimiter(30, time.Second),
		closeChan: make(chan struct{}),
	}
}

func (client *ScreenCaptureClient) AddObserver(observer observer.MessageObserver) {
	client.msgObservers = append(client.msgObservers, observer)
}

func (client *ScreenCaptureClient) Start() error {
	u := url.URL{Scheme: "ws", Host: client.addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	client.conn = c
	util.Info("connect ws server [%s] successfully", client.addr)

	go func() {
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				return
			}
			if bytes.Equal(data, []byte("OK")) {
				continue
			}

			// rate limit
			if !client.limiter.Allow() {
				continue
			}

			ready := true
			for _, observer := range client.msgObservers {
				if !observer.IsReady() {
					ready = false
					break
				}
			}
			if ready {
				go client.handleMessage(data)
			}
		}
	}()

	go client.notify()

	return nil
}

func (client *ScreenCaptureClient) handleMessage(data []byte) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}
	client.imgChan <- img
}

func (client *ScreenCaptureClient) notify() {
	for {
		select {
		case img, ok := <-client.imgChan:
			if !ok {
				return
			}
			for _, observer := range client.msgObservers {
				if observer.IsReady() {
					observer.Update(img)
				}
			}
		}
	}
}

func (client *ScreenCaptureClient) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
	close(client.imgChan)
	client.limiter.Done()
}
