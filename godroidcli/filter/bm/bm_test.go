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

package bm

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestBoyerMooreSearch(t *testing.T) {
	data, err := ioutil.ReadFile("test_data/words.txt")
	if err != nil {
		t.Fatal(err)
	}
	pattern := "you"
	indexes := BoyerMooreSearch(data, []byte(pattern), false, false)
	t.Log(indexes)
}

func TestSubstring(t *testing.T) {
	data, err := ioutil.ReadFile("test_data/words.txt")
	if err != nil {
		t.Fatal(err)
	}
	pattern := "the"
	subs := Substring(data, []byte(pattern), false)
	for _, d := range subs {
		t.Log(string(d))
	}
}

func TestSimpleString(t *testing.T) {
	str := "hEllOhloXhello HEllo HELLo hellO.hello .. HeLlo"

	pattern := []byte("hello")
	indexes := BoyerMooreSearch([]byte(str), pattern, false, false)
	for _, index := range indexes {
		t.Log(str[index : index+len(pattern)])
	}
}

// BUGS: boyer search
func TestSimple2String(t *testing.T) {
	str := "AAAAAAAAAA"
	indexes := BoyerMooreSearch([]byte(str), []byte("AAA"), false, true)
	fmt.Println(indexes)
}

func testBmIndex(data []byte, pattern string) {
	_ = BoyerMooreSearch(data, []byte(pattern), false, true)
}

func testBmSubstring(data []byte, pattern string) {
	_ = Substring(data, []byte(pattern), true)
}

func BenchmarkBmIndex(b *testing.B) {
	data, err := ioutil.ReadFile("test_data/test.txt")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		testBmIndex(data, "you")
	}
}

func BenchmarkBmSubstring(b *testing.B) {
	data, err := ioutil.ReadFile("test_data/test.txt")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		testBmSubstring(data, "you")
	}
}
