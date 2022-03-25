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

package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josexy/godroidcli/android/api/middleware"
	"github.com/josexy/godroidcli/android/api/router"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/status"
	"github.com/josexy/godroidcli/util"
)

type Startup struct {
	ctx    context.Context
	proxy  *internal.SessionProxy
	engine *gin.Engine
}

// New create an api proxy service based on the currently active session
func New(proxy *internal.SessionProxy) *Startup {
	return &Startup{
		ctx:   context.Background(),
		proxy: proxy,
	}
}

func (s *Startup) initApiRouters() {
	groups := router.NewApiGroups(s.engine)
	groups.AddGroup(internal.Pm, s.proxy).InitPm()
	groups.AddGroup(internal.Fs, s.proxy).InitFs()
	groups.AddGroup(internal.Net, s.proxy).InitNet()
	groups.AddGroup(internal.Di, s.proxy).InitDevice()
	groups.AddGroup(internal.Ctrl, s.proxy).InitCtrl()
	groups.AddGroup(internal.Ms, s.proxy).InitMediaStore()
	groups.AddGroup(internal.Sms, s.proxy).InitSms()
	groups.AddGroup(internal.Contact, s.proxy).InitContact()
	groups.AddGroup(internal.CallLog, s.proxy).InitCallLog()
	groups.AddGroup(internal.Phone, s.proxy).InitPhone()
	groups.InitAll()
}

func (s *Startup) prepare() {

	// quiet mode
	gin.DefaultWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)

	s.engine = gin.New()
	s.engine.Use(middleware.Cors(), gin.Recovery())
	s.initApiRouters()
}

// Start start api server and accept an interrupt signal to quit
func (s *Startup) Start() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	config := util.GetConfig()
	if config == nil {
		return status.ErrLoadConfigFailed
	}

	s.prepare()

	server := http.Server{
		Addr:    config.Address,
		Handler: s.engine,
	}
	go func() {
		_ = server.ListenAndServe()
	}()

	util.Info("press %s to quit api server", util.Green("Ctrl+C"))
	util.Info("listen address %s", util.Green(config.Address))

	<-util.MakeInterruptChan()

	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
