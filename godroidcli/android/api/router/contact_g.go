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
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/status"
)

type contactId struct {
	ID string `json:"id"`
}

type contactUri struct {
	Uri string `json:"uri"`
}

type contactAdd struct {
	Name   string `json:"name"`
	Phones []struct {
		Phone string `json:"phone"`
		Type  string `json:"type"`
	} `json:"phones"`
	Emails []struct {
		Email string `json:"email"`
		Type  string `json:"type"`
	} `json:"emails"`
}

func (r *ApiRouter) InitContact() {

	// /api/contact/all
	r.mapHandlers[internal.All] = func(ctx *gin.Context) {
		wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
			status.Success, wrapper.NewRpcCall(r.GetAllContactInfo())))
	}

	// /api/contact/get
	r.mapHandlers[internal.Get] = func(ctx *gin.Context) {
		var value contactId
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCall(r.GetContactInfo(value.ID))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/contact/delete
	r.mapHandlers[internal.Delete] = func(ctx *gin.Context) {
		var value contactUri
		if err := ctx.ShouldBind(&value); err == nil {
			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.DeleteContact(value.Uri))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}

	// /api/contact/add
	/* application/json:
	{
		"name" : "test",
		"phones": [
			{
				"phone": "123",
				"type": "Home"
			}
		],
		"emails": [
			{
				"email": "123@123.com",
				"type": "Mobile"
			}
		]
	}
	*/
	r.mapHandlers[internal.Add] = func(ctx *gin.Context) {
		var value contactAdd
		if err := ctx.ShouldBind(&value); err == nil {
			var phones []pb.StringPair
			var emails []pb.StringPair

			for i := 0; i < len(value.Phones); i++ {
				phones = append(phones, pb.StringPair{First: value.Phones[i].Phone, Second: value.Phones[i].Type})
			}

			for i := 0; i < len(value.Emails); i++ {
				emails = append(emails, pb.StringPair{First: value.Emails[i].Email, Second: value.Emails[i].Type})
			}

			wrapper.ResponseJson(ctx, wrapper.SerializeResponseFromRpcCall(
				status.Success, wrapper.NewRpcCallError(r.AddContact(value.Name, phones, emails))))
		} else {
			wrapper.ResponseError(ctx, err)
		}
	}
}
