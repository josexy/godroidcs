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
	"testing"
)

func TestCalcBytes(t *testing.T) {
	t.Log(CalcFileBytes(512))
	t.Log(CalcFileBytes(1024))
	t.Log(CalcFileBytes(1024 * 512))
	t.Log(CalcFileBytes(1024 * 1024))
	t.Log(CalcFileBytes(1024 * 1024 * 512))
	t.Log(CalcFileBytes(1024 * 1024 * 1024))
	t.Log(CalcFileBytes(1024 * 1024 * 1024 * 512))
	t.Log(CalcFileBytes(1024 * 1024 * 1024 * 1024))
	t.Log(CalcFileBytes(1024 * 1024 * 1024 * 1024 * 512))
	t.Log(CalcFileBytes(1024 * 1024 * 1024 * 1024 * 1024))
}

func TestTrimQuote(t *testing.T) {
	t.Logf("%q\n", Trim(`"hello world"`))
	t.Logf("%q\n", Trim(`hello world`))
	t.Logf("%q\n", Trim(`'hello world'`))
	t.Logf("%q\n", Trim(`"hello world'`))
	t.Logf("%q\n", Trim(`"'hello world'"`))
}

func TestTimeOfHMS(t *testing.T) {
	t.Log(TimeOfHMS(45))
	t.Log(TimeOfHMS(60))
	t.Log(TimeOfHMS(60 * 60))
	t.Log(TimeOfHMS(60 * 60 * 60))
	t.Log(TimeOfHMS(60*60*12 + 60*5 + 125))
}

func TestStringToBytes(t *testing.T) {
	s := "hello world"
	b := StringToBytes(s)
	t.Log(s)
	t.Log(b)
}

func TestBytesToString(t *testing.T) {
	b := []byte("hello world")
	a := BytesToString(b)
	t.Log(b)
	t.Log(a)
}

func TestRawWidthStr(t *testing.T) {
	s := Green("hello %s world", HiBlue("Golang"))
	t.Log(s)
	s = RawWidthStr(s)
	t.Log(s)
}
