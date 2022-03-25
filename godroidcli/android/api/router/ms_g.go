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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josexy/godroidcli/android/api/wrapper"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/status"
)

type msType struct {
	Type string `json:"type"`
}

type msUri struct {
	Uri string `json:"uri"`
}

func (r *ApiRouter) InitMediaStore() {

	// /api/ms/mediafiles
	r.mapHandlers[internal.MediaFiles] = func(ctx *gin.Context) {
		var value msType
		if err := ctx.ShouldBind(&value); err == nil {
			if index, ok := protobuf.MediaType_Type_value[strings.ToUpper(value.Type)]; ok {
				wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
					status.Success, wrapper.NewRpcCall(r.GetMediaFilesInfo(protobuf.MediaType_Type(index)))))
			} else {
				wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
					status.Success, wrapper.NewRpcCallError(status.ErrMediaFileType)))
			}
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ms/delete
	r.mapHandlers[internal.Delete] = func(ctx *gin.Context) {
		var value msUri
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.DeleteMediaFile(value.Uri))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ms/thumbnail
	r.mapHandlers[internal.Thumbnail] = func(ctx *gin.Context) {
		var value msUri
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.GetMediaFileThumbnail(value.Uri, ctx.Writer, nil))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/ms/download
	r.mapHandlers[internal.Download] = func(ctx *gin.Context) {
		var value msUri
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.DownloadMediaFile(value.Uri, ctx.Writer, nil))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}
}
