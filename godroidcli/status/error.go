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

package status

import "errors"

var (
	ErrNoSession            = errors.New("there is no active session")
	ErrSessionUnavailable   = errors.New("the current session connection is not available")
	ErrConnNotExist         = errors.New("connection does not exist")
	ErrConnTimeout          = errors.New("connect to server timeout")
	ErrorIllegalOperation   = errors.New("illegal operation, please check whether the input parameters are correct")
	ErrNotFoundOrNotExisted = errors.New("no device or port forwarding rules exist")
	ErrCannotFindForward    = errors.New("no port forwarding rules exist")
	ErrDeviceNotFound       = errors.New("no device rules exist")
	ErrForwardPortInvalid   = errors.New("invalid port forwarding")
	ErrCannotParseHtml      = errors.New("could not parse html table")
	ErrEmptyString          = errors.New("parse string literal is empty")
	ErrCannotParseCmd       = errors.New("could not parse command")
	ErrProvideParams        = errors.New("more parameters need to be provided")
	ErrLoadConfigFailed     = errors.New("failed to load configuration file")
	ErrUnmarshalToJson      = errors.New("could not unmarshal bytes to json object")
	ErrMarshalProtoToJSON   = errors.New("could not marshal proto message to json object")
	ErrPathEmpty            = errors.New("path cannot be empty")
)

var (
	ErrMediaFileType = errors.New("medaifile type invalid")
)
