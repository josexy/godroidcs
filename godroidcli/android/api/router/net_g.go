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
	"github.com/josexy/godroidcli/android/api/wrapper"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/status"
)

func (r *ApiRouter) InitNet() {

	// /api/net/info
	r.mapHandlers[internal.Info] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetNetworkInfo())))
	}

	// /api/net/wifi
	r.mapHandlers[internal.Wifi] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetCurrentWifiInfo())))
	}

	// /api/net/connectivity
	r.mapHandlers[internal.HasNetwork] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.CheckNetworkConnectivity())))
	}

	// /api/net/active_network
	r.mapHandlers[internal.ActiveNetwork] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetActivityNetworkDetailInfo())))
	}

	// /api/net/public_network
	r.mapHandlers[internal.PublicNetwork] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetPublicNetworkInfo())))
	}

	// /api/net/scan_wifi
	r.mapHandlers[internal.ScanWifi] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.ScanWifi())))
	}
}
