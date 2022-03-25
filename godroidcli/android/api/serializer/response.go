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

package serializer

import (
	"github.com/josexy/godroidcli/status"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func BuildResponseWithData(code int, data interface{}) Response {
	return Response{
		Code:    code,
		Message: status.GetMessageFromCode(code),
		Data:    data,
	}
}

func BuildErrorResponse(err error) Response {
	return Response{
		Code:    status.Error,
		Message: status.GetMessageFromCode(status.Error),
		Error:   err.Error(),
	}
}
