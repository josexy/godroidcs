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

type phoneNumber struct {
	Number string `json:"number"`
}

func (r *ApiRouter) InitPhone() {

	var value phoneNumber

	// /api/phone/dial
	r.mapHandlers[internal.Dial] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.DialPhone(value.Number))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/phone/call
	r.mapHandlers[internal.Call] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.CallPhone(value.Number))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}
}
