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

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/josexy/godroidcli/android/internal"
)

const apiPrefix = "/api/"

type ApiRouter struct {
	group *gin.RouterGroup
	*internal.SessionProxy
	mapHandlers map[string]gin.HandlerFunc
}

type ApiGroups struct {
	rg        *gin.RouterGroup
	routerMap map[string]*ApiRouter
}

func NewApiGroups(engine *gin.Engine) *ApiGroups {
	return &ApiGroups{
		rg:        engine.Group(apiPrefix),
		routerMap: make(map[string]*ApiRouter),
	}
}

func (api *ApiGroups) AddGroup(name string, proxy *internal.SessionProxy) *ApiRouter {
	router, ok := api.routerMap[name]
	if !ok {
		router = &ApiRouter{
			SessionProxy: proxy,
			group:        api.rg.Group(name),
			mapHandlers:  make(map[string]gin.HandlerFunc),
		}
		api.routerMap[name] = router
	}
	return router
}

func (api *ApiGroups) initRouter(router *ApiRouter) {
	if router != nil {
		for path, handler := range router.mapHandlers {
			router.group.POST("/"+path, handler)
		}
	}
}

func (api *ApiGroups) InitAll() {
	for _, router := range api.routerMap {
		api.initRouter(router)
	}
}
