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

package main

import (
	"flag"
)

// this module is used for creating connection and opening a session quickly
// some supported commands are listed as follows:
// connect to android server by serial number
// > godroidcli -device SERIAL_NUMBER:PORT
// connect to android server by TCP/IP
// > godroidcli -address IP:PORT
var (
	device  string
	address string
)

func init() {
	flag.StringVar(&device, "device", "", "android device serial number")
	flag.StringVar(&address, "address", "", "android device ip address and port")
}

func GetDeviceValue() (string, bool) {
	if device == "" {
		return device, false
	}
	return device, true
}

func GetAddressValue() (string, bool) {
	if address == "" {
		return address, false
	}
	return address, true
}

func ParseCommand() {
	flag.Parse()
}
