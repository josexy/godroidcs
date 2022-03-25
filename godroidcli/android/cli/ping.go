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

package cli

import (
	"context"
	"time"

	pb "github.com/josexy/godroidcli/protobuf"
	"google.golang.org/grpc"
)

const DelayTime = time.Millisecond * 1000

const (
	alive = iota
	unavailable
)

type pingContext struct {
	client     pb.PingTestClient
	ctx        context.Context
	cancel     context.CancelFunc
	statusChan chan int
}

func NewPingContext(conn *grpc.ClientConn, ctx context.Context) *pingContext {
	pc := &pingContext{
		client:     pb.NewPingTestClient(conn),
		statusChan: make(chan int),
	}
	pc.ctx, pc.cancel = context.WithCancel(ctx)
	return pc
}

func (p *pingContext) ping() {
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				p.cancel()
				return
			default:
			}
			// ping rpc server
			_, err := p.client.Ping(p.ctx, &pb.Empty{})
			if err != nil {
				// rpc server is unavailable
				p.statusChan <- unavailable
			} else {
				// rpc server is alive
				p.statusChan <- alive
			}
			// delay
			time.Sleep(DelayTime)
		}
	}()
}
