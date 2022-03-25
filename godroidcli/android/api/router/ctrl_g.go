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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/josexy/godroidcli/android/api/wrapper"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/status"
)

type ctrlValue struct {
	Value string `json:"value"`
}

type ctrlText struct {
	Text string `json:"text"`
}

type ctrlMode struct {
	Mode string `json:"mode"`
}

func (r *ApiRouter) InitCtrl() {

	// /api/ctrl/get_brightness
	r.mapHandlers[internal.GetScreenBrightness] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetScreenBrightness())))
	}

	// /api/ctrl/set_brightness
	r.mapHandlers[internal.SetScreenBrightness] = func(ctx *gin.Context) {
		var value ctrlValue
		if err := ctx.ShouldBind(&value); err == nil {
			val, err := strconv.Atoi(value.Value)
			if err != nil {
				wrapper.ResponseError(ctx, err)
			} else {
				wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
					status.Success, wrapper.NewRpcCallError(r.SetScreenBrightness(val))))
			}
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ctrl/get_clipboard
	r.mapHandlers[internal.GetClipboard] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetClipboardText())))
	}

	// /api/ctrl/set_clipboard
	r.mapHandlers[internal.SetClipboard] = func(ctx *gin.Context) {
		var value ctrlText
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.SetClipboardText(value.Text))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ctrl/get_brightness_mode
	r.mapHandlers[internal.GetScreenBrightnessMode] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetScreenBrightnessMode())))
	}

	// /api/ctrl/set_brightness_mode
	r.mapHandlers[internal.SetScreenBrightnessMode] = func(ctx *gin.Context) {
		var value ctrlMode
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.SetScreenBrightnessMode(value.Mode == "auto"))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ctrl/get_volume
	r.mapHandlers[internal.GetVolume] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetVolume())))
	}

	// /api/ctrl/set_volume
	r.mapHandlers[internal.SetVolume] = func(ctx *gin.Context) {
		var value ctrlValue
		if err := ctx.ShouldBind(&value); err == nil {
			val, err := strconv.Atoi(value.Value)
			if err != nil {
				wrapper.ResponseError(ctx, err)
			} else {
				wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
					status.Success, wrapper.NewRpcCallError(r.SetVolume(val))))
			}
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ctrl/increase_volume
	r.mapHandlers[internal.IncreaseVolume] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCallError(r.IncreaseVolume())))
	}

	// /api/ctrl/decrease_volume
	r.mapHandlers[internal.DecreaseVolume] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCallError(r.DecreaseVolume())))
	}

	// /api/ctrl/start_screen_capture
	// start the websocket service on the device
	r.mapHandlers[internal.StartScreenCapture] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCallError(r.StartScreenCapture())))
	}

	// /api/ctrl/stop_screen_capture
	// stop the websocket service on the device
	r.mapHandlers[internal.StopScreenCapture] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCallError(r.StopScreenCapture())))
	}

	// /api/ctrl/start_screen_record
	r.mapHandlers[internal.StartScreenRecord] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCallError(r.StartScreenRecord())))
	}

	// /api/ctrl/stop_screen_record
	r.mapHandlers[internal.StopScreenRecord] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.StopScreenRecord())))
	}
}
