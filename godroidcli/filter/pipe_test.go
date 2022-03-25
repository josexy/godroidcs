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
	"io/ioutil"
	"testing"

	"github.com/josexy/godroidcli/util"
)

func TestStreamPipe_IGrep(t *testing.T) {
	str := `
	hello world
	HelLOWorld
	Hel1o
	WorldhElLoworLd
	`
	pipe := CreatePipe([]byte(str))
	pipe.IGrep([]byte("hello")).PrintColor(util.HiGreenAttr)
}

func TestStreamPipe_Grep0(t *testing.T) {
	str := `

hello world👋 hello wo😄rld  😄 你好世界

😄hello😄😄world hello 世界你好 world

你好👋
世界

世界你好world世界

世界👋你好

`
	pipe := CreatePipe([]byte(str))
	pipe.Grep([]byte("👋")).Print()
	pipe.Grep([]byte("hello")).Grep([]byte("世界你好")).PrintColor(util.GreenAttr)
}

func TestStreamPipe_Grep1(t *testing.T) {
	data, _ := ioutil.ReadFile("bm/test_data/test.txt")
	CreatePipe(data).
		Grep([]byte("copyright")).
		Print().
		PrintColor(util.BlueAttr)
}

func TestStreamPipe_Grep2(t *testing.T) {
	data, _ := ioutil.ReadFile("bm/test_data/test.txt")
	CreatePipe(data).
		Grep([]byte("You")).
		Grep([]byte("copyright")).
		PrintColor(util.RedAttr)
}

func TestStreamPipe_Print(t *testing.T) {
	data := []byte(
		`
#    hello golang
#
#		hello c++
		hello java
		hello linux
!
!		#golang& Golang
		
#		
		go macos
		
#		windows
#		html,css,js
#		test...
		`,
	)
	CreatePipe(data).
		Grep([]byte("golang")).PrintColor(util.GreenAttr). // green
		Print().                                           // without color
		Grep([]byte("hello")).PrintColor(util.RedAttr)     // red
}

func BenchmarkStreamPipe_Grep(b *testing.B) {
	data, _ := ioutil.ReadFile("bm/test_data/test.txt")
	pipe := CreatePipe(data)
	for i := 0; i < b.N; i++ {
		pipe.Grep([]byte("copyright"))
	}
}
