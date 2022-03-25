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
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/josexy/godroidcli/android/api/wrapper"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/status"
	"github.com/josexy/godroidcli/util"
)

const UploadFileId = "upload_file"

type fsPath struct {
	Path string `json:"path"`
}

type fsList struct {
	fsPath
	Mode string `json:"mode"`
}

type fsBaseFileTree struct {
	fsList
	ID string `json:"id"`
}

type fsPathPair struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

type fsPathText struct {
	fsPath
	Text string `json:"text"`
}

func (r *ApiRouter) InitFs() {

	// /api/fs/file_tree
	r.mapHandlers[internal.GetBaseFileTree] = func(ctx *gin.Context) {
		var value fsBaseFileTree
		if err := ctx.ShouldBind(&value); err == nil {
			res, err := r.GetBaseFileTree(value.Path, value.ID, value.Mode)
			var rpcCall *wrapper.RpcCall
			if err != nil {
				rpcCall = wrapper.NewRpcCallError(err)
			} else {
				val, err := wrapper.UnmarsalToJson(util.StringToBytes(res.Value))
				if err != nil {
					rpcCall = wrapper.NewRpcCallError(err)
				} else {
					rpcCall = wrapper.NewRpcCallData(val, err)
				}
			}
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(status.Success, rpcCall))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/list
	r.mapHandlers[internal.List] = func(ctx *gin.Context) {
		var value fsList
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.ListDir(value.Path, value.Mode))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/create
	r.mapHandlers[internal.Create] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.CreateFile(value.Path))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/delete
	r.mapHandlers[internal.Delete] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.DeleteFile(value.Path))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/mkdir
	r.mapHandlers[internal.MkDir] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.MkDir(value.Path))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/rmdir
	r.mapHandlers[internal.RmDir] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.RmDir(value.Path))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/move
	r.mapHandlers[internal.Move] = func(ctx *gin.Context) {
		var value fsPathPair
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.Move(value.Src, value.Dest))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/rename
	r.mapHandlers[internal.Rename] = func(ctx *gin.Context) {
		var value fsPathPair
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.Rename(value.Src, value.Dest))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/copy
	r.mapHandlers[internal.Copy] = func(ctx *gin.Context) {
		var value fsPathPair
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.Copy(value.Src, value.Dest))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/write
	r.mapHandlers[internal.WriteText] = func(ctx *gin.Context) {
		var value fsPathText
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.WriteText(value.Path, value.Text))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/append
	r.mapHandlers[internal.AppendText] = func(ctx *gin.Context) {
		var value fsPathText
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.AppendText(value.Path, value.Text))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/read
	r.mapHandlers[internal.ReadText] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.ReadText(value.Path))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/download
	// download file and redirect download stream to gin.Writer
	r.mapHandlers[internal.Download] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallBytes(r.DownloadFile(value.Path, ctx.Writer, nil))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/fs/upload
	// redirect upload file stream to Android device
	r.mapHandlers[internal.Upload] = func(ctx *gin.Context) {
		var value fsPath
		if err := ctx.ShouldBind(&value); err == nil {
			value.Path = strings.TrimSpace(value.Path)

			if value.Path == "" {
				wrapper.ResponseError(ctx, status.ErrPathEmpty)
				return
			}
			file, err := ctx.FormFile(UploadFileId)
			if err != nil {
				wrapper.ResponseError(ctx, err)
				return
			}
			if value.Path[len(value.Path)-1] == '/' {
				value.Path = filepath.Join(value.Path, file.Filename)
			}

			reader, err := file.Open()
			if err != nil {
				wrapper.ResponseError(ctx, err)
				return
			}
			defer func() { _ = reader.Close() }()
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.UploadFile(reader, value.Path))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}
}
