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

package util

import (
	"os"
	"os/exec"
)

var clearFuncMap = map[string]func(){
	"linux":   clear_unix,
	"darwin":  clear_unix,
	"windows": clear_windows,
}

// ClearScreen clear terminal screen
// see: https://stackoverflow.com/a/22896706/8523508
func ClearScreen(os string) {
	if f, ok := clearFuncMap[os]; ok {
		f()
	} else {
		// default
		clearFuncMap["linux"]()
	}
}

// clear_windows execute command and abort: `cmd.exe /c cls`
func clear_windows() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// clear_unix clear terminal screen on unix and unix-like os
// for Unix and Unix-Like os, you also can use `fmt.Fprint(os.Stdout, "\033[2J\033[1;1H")`
func clear_unix() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
