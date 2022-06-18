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
	"fmt"
	"time"

	"github.com/josexy/godroidcli/android/cli/resolver"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	"github.com/josexy/godroidcli/status"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

type Session struct {
	sn        string // the serial number of device
	address   string // rpc server ip address
	port      int    // rpc server port
	status    int    // connection status
	conn      *grpc.ClientConn
	pc        *pingContext
	adb       *AdbCmd
	ctx       context.Context
	cancel    context.CancelFunc
	proxy     *internal.SessionProxy
	resolvers map[string]*resolver.ResolverContext
}

const OpenSessionTimeout = time.Second * 4

func NewSession(ctx context.Context, sn, address string, port int, adb *AdbCmd) (*Session, error) {
	s := &Session{
		sn:      sn,
		address: address,
		port:    port,
		adb:     adb,
		status:  alive,
	}

	s.ctx, s.cancel = context.WithCancel(ctx)

	err := s.openSession(OpenSessionTimeout)
	if err != nil {
		return nil, err
	}
	util.Info("successfully connected to the server [%s]", s)
	// set context value
	s.adb.serialNumber = sn

	// supported resolvers
	s.addResolver(s.ctx, internal.Pm, resolver.NewPackageManager(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Fs, resolver.NewFileSystem(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Di, resolver.NewDevice(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Net, resolver.NewNetwork(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Ctrl, resolver.NewController(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Ms, resolver.NewMediaStore(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Sms, resolver.NewSms(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Contact, resolver.NewContact(s.conn), s.adb)
	s.addResolver(s.ctx, internal.CallLog, resolver.NewCallLog(s.conn), s.adb)
	s.addResolver(s.ctx, internal.Phone, resolver.NewPhone(s.conn), s.adb)
	return s, nil
}

func (s *Session) String() string {
	return fmt.Sprintf("%s@%s:%s", util.Green(s.sn),
		util.Yellow(s.address), util.Yellow(util.IntToStr(s.port)))
}

// CreateSessionProxy internal.SessionProxy acts as an intermediate proxy,
// associating all resolvers for the session with the api server
func (s *Session) CreateSessionProxy() *internal.SessionProxy {
	if s.proxy == nil {
		s.proxy = &internal.SessionProxy{}
		s.proxy.IPackageManager = s.GetResolver(internal.Pm).(*resolver.PackageManager)
		s.proxy.INetwork = s.GetResolver(internal.Net).(*resolver.Network)
		s.proxy.IFileSystem = s.GetResolver(internal.Fs).(*resolver.FileSystem)
		s.proxy.IDevice = s.GetResolver(internal.Di).(*resolver.Device)
		s.proxy.IController = s.GetResolver(internal.Ctrl).(*resolver.Controller)
		s.proxy.IMediaStore = s.GetResolver(internal.Ms).(*resolver.MediaStore)
		s.proxy.ISms = s.GetResolver(internal.Sms).(*resolver.Sms)
		s.proxy.IContact = s.GetResolver(internal.Contact).(*resolver.Contact)
		s.proxy.ICallLog = s.GetResolver(internal.CallLog).(*resolver.CallLog)
		s.proxy.IPhone = s.GetResolver(internal.Phone).(*resolver.Phone)
	}
	return s.proxy
}

// addResolver register resolver
func (s *Session) addResolver(ctx context.Context, name string, r resolver.Resolver, cmd resolver.ExecCmd) {
	if s.resolvers == nil {
		s.resolvers = make(map[string]*resolver.ResolverContext)
	}

	// save websocket server address into context
	ctx = context.WithValue(ctx, resolver.ServerAddrContextKey,
		fmt.Sprintf("%s:%d", s.address, s.port+1))

	// attach current session to all resolvers
	rc := resolver.NewResolverContext(ctx, r, s, cmd)
	r.SetContext(rc)
	s.resolvers[name] = rc
}

func (s *Session) GetResolver(name string) resolver.Resolver {
	if ctx, ok := s.resolvers[name]; ok {
		return ctx.Resolver
	}
	return nil
}

func (s *Session) pong() {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				// cancel ping goroutine
				s.cancel()
				return
			case s.status = <-s.pc.statusChan:
			}
		}
	}()
}

// openSession connect to rpc server within timeout
func (s *Session) openSession(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithNoProxy(),
		grpc.WithBlock(),
	}
	// try to connect to rpc server
	s.conn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%d", s.address, s.port), opts...)
	// connect timeout
	if err != nil {
		err = status.ErrConnTimeout
		return err
	}
	s.pc = NewPingContext(s.conn, s.ctx)
	s.pc.ping()
	s.pong()
	return
}

func (s *Session) CloseSession() (err error) {
	if s.conn == nil {
		err = status.ErrConnNotExist
		return
	}
	err = s.conn.Close()
	if s.cancel != nil {
		s.cancel()
	}
	return
}

func (s *Session) executeCommand(param filter.Param) {
	// recover
	defer util.RecoverIllegalOption()

	name := param.Args[1]
	if m, ok := s.resolvers[name]; ok {
		m.DoProcess(filter.Param{Node: param.Node, Args: param.Args[2:]})
	} else {
		util.Error("[%s] subcommand not supported", name)
	}
}
