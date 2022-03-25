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

package wrapper

import (
	"encoding/json"

	"github.com/josexy/godroidcli/android/api/serializer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnmarsalToJson(bs []byte) (val interface{}, err error) {
	err = json.Unmarshal(bs, &val)
	return
}

// MarshalToBytes convert Protobuf message to JSON string
func MarshalToBytes(msg proto.Message) (interface{}, error) {
	data, err := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		return nil, err
	}
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func SerializeResponseFromRpcCall(code int, call *RpcCall) serializer.Response {
	if call.Error != nil {
		// error
		return serializer.BuildErrorResponse(call.Error)
	} else if call.Bytes {
		// bytes stream
	} else if call.Msg != nil {
		// proto message
		object, err := MarshalToBytes(call.Msg)
		if err != nil {
			return serializer.BuildErrorResponse(err)
		}
		return serializer.BuildResponseWithData(code, object)
	}
	// normal data
	return serializer.BuildResponseWithData(code, call.Data)
}
