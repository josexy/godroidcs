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

package filter

import (
	"fmt"
	"testing"
)

var p = NewCmdParser()

func testCmdParser(s string) {
	fmt.Println("--> " + s)
	if err := p.Parse(s); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
		p.dumpSyntaxTreeTest()
	}
}

func TestCmdParser_GOOD_Parse(t *testing.T) {
	var s string
	s = `cmd | grep hello "hello golang"`
	testCmdParser(s)
	s = `cmd pm all_packages`
	testCmdParser(s)
	s = `cmd device info |grep B|export csv 1.txt |grep r o|export 1.txt csv 2.txt html 3.txt`
	testCmdParser(s)
	s = `cmd pm package com.android.chrome | grep android | print  `
	testCmdParser(s)
	s = `cmd pm	package "com.android.chrome" | grep "android" | export file.txt`
	testCmdParser(s)
	s = `cmd fs list /storage/emulated/0/|grep txt|grep "hello world" |export "hello world" | grep hello world`
	testCmdParser(s)
	s = `cmd pm all_packages |grep android google | export markdown output.md
	|export csv output.csv| export html output.html | export output.txt | print
	|export csv output2.csv html "output.html" output.txt
	|grep "hello world"|grep hello golang`
	testCmdParser(s)
}

func TestCmdParser_BAD_Parse(t *testing.T) {
	var s string
	s = `cmd pm package xxx | export markdown`
	testCmdParser(s)
	s = `cmd pm package xxx | export`
	testCmdParser(s)
	s = `cmd pm package xxx | | grep x`
	testCmdParser(s)
	s = `cmd pm package | out | grep x`
	testCmdParser(s)
	s = `cmd pm package | grep xx yy zz | print |`
	testCmdParser(s)
}
