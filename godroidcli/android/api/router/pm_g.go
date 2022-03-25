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

type pmName struct {
	Name string `json:"name"`
}

func (r *ApiRouter) InitPm() {
	value := pmName{}

	// /api/pm/all_packages
	r.mapHandlers[internal.AllPackages] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetAllPackageInfo())))
	}

	// /api/pm/all_user_packages
	r.mapHandlers[internal.AllUserPackages] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetAllUserPackageInfo())))
	}

	// /api/pm/all_system_packages
	r.mapHandlers[internal.AllSystemPackages] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetAllSystemPackageInfo())))
	}

	// /api/pm/package
	r.mapHandlers[internal.Package] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetPackageInfo(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/application
	r.mapHandlers[internal.Application] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetApplicationInfo(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/size
	r.mapHandlers[internal.AppSize] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetAppSize(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/permissions
	r.mapHandlers[internal.Permissions] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetPermissions(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/activities
	r.mapHandlers[internal.Activities] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetActivities(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/services
	r.mapHandlers[internal.Services] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetServices(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/receivers
	r.mapHandlers[internal.Receivers] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetReceivers(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/providers
	r.mapHandlers[internal.Providers] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetProviders(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/sharedlibs
	r.mapHandlers[internal.SharedLibs] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetSharedLibs(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/icon
	r.mapHandlers[internal.GetIcon] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(status.Success,
				wrapper.NewRpcCallBytes(r.GetAppIcon(value.Name, ctx.Writer, nil))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/apk
	r.mapHandlers[internal.GetApk] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(status.Success,
				wrapper.NewRpcCallBytes(r.GetApkFile(value.Name, ctx.Writer, nil))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/uninstall
	r.mapHandlers[internal.Uninstall] = func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.UninstallApp(value.Name))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/pm/install
	r.mapHandlers[internal.Install] = func(ctx *gin.Context) {
		file, err := ctx.FormFile(UploadFileId)
		if err != nil {
			wrapper.ResponseError(ctx, err)
			return
		}
		reader, err := file.Open()
		if err != nil {
			wrapper.ResponseError(ctx, err)
			return
		}
		defer func() { _ = reader.Close() }()
		wrapper.SerializeResponseFromRpcCall(status.Success, wrapper.NewRpcCallError(r.InstallApk(reader)))
	}
}
