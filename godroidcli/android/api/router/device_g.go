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

func (r *ApiRouter) InitDevice() {

	// /api/device/info
	r.mapHandlers[internal.Info] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetDeviceInfo())))
	}

	// /api/device/system
	r.mapHandlers[internal.System] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetSystemInfo())))
	}

	// /api/device/battery
	r.mapHandlers[internal.Battery] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetBatteryInfo())))
	}

	// /api/device/location
	r.mapHandlers[internal.Location] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetLocationInfo())))
	}

	// /api/device/display
	r.mapHandlers[internal.Display] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetDisplayInfo())))
	}

	// /api/device/mem
	r.mapHandlers[internal.Mem] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetMemoryInfo())))
	}

	// /api/device/storage
	r.mapHandlers[internal.Storage] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetStorageSpaceInfo())))
	}

	// /api/device/gpu
	r.mapHandlers[internal.GPU] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetGPUInfo())))
	}

	// /api/device/cpu
	r.mapHandlers[internal.CPU] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetCPUsFrequency())))
	}
}
