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

var codeMessageMap = map[int]string{
	Success:     "success",
	Error:       "error",
	Empty:       "empty message",
	InvalidCode: "invalid response code",
}

func GetMessageFromCode(code int) string {
	if message, ok := codeMessageMap[code]; ok {
		return message
	}
	return codeMessageMap[InvalidCode]
}
