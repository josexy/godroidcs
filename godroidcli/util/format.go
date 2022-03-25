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
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

type Level int

const (
	OFF Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var doubleWidth = []*unicode.RangeTable{
	unicode.Han,
	unicode.Hangul,
	unicode.Hiragana,
	unicode.Katakana,
}

func Trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\"'"))
}

func CalcFileBytes(value int64) string {
	return CalcBytes(value, 1024)
}

func CalcCommonBytes(value int64) string {
	return CalcBytes(value, 1000)
}

func CalcBytes(value int64, factor int) (s string) {
	v := float64(value)
	f := float64(factor)
	t := 0
	for v >= f {
		v /= f
		if t >= 3 {
			break
		}
		t++
	}
	buf := bytes.NewBufferString(fmt.Sprintf("%.2f", v))
	switch t {
	case 0:
		buf.WriteString("B")
	case 1:
		buf.WriteString("KB")
	case 2:
		buf.WriteString("MB")
	default:
		buf.WriteString("GB")
	}
	return buf.String()
}

func TimeOf(timestamp int64) string {
	if timestamp <= 0 {
		return "N/A"
	}
	return time.UnixMilli(timestamp).Format("2006-01-02 15:04:05")
}

// TimeOfHMS convert total seconds to Hour:Minute:Second format
func TimeOfHMS(t int64) string {
	hour := t / 60 / 60 % 60
	minute := t / 60 % 60
	second := t % 60
	return fmt.Sprintf("%d:%d:%d", hour, minute, second)
}

func TimeOfNow() string {
	return TimeOf(time.Now().UnixMilli())
}

func TimeOfNowForFile() string {
	return time.UnixMilli(time.Now().UnixMilli()).Format("2006_01_02_15_04_05")
}

func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Int32ToStr(i int32) string {
	return IntToStr(int(i))
}

func FloatToStr(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func Float32ToStr(f float32) string {
	return fmt.Sprintf("%.2f", f)
}

func BoolToStr(b bool) string {
	return strconv.FormatBool(b)
}

// rawWidthStr0 remove escape sequences \033[m
func rawWidthStr0(str string) string {
	buf := bytes.NewBuffer(nil)
	re := regexp.MustCompile("\033\\[.*?m(.*?)\033\\[0m")
	data := []byte(str)

	list := re.Split(str, -1)
	indexes := re.FindAllStringSubmatchIndex(str, -1)
	m, n := len(list), len(indexes)
	var i, j int
	for i < m && j < n {
		buf.WriteString(list[i])
		buf.Write(data[indexes[j][2]:indexes[j][3]])
		i++
		j++
	}
	for i < m {
		buf.WriteString(list[i])
		i++
	}
	for j < n {
		buf.Write(data[indexes[j][2]:indexes[j][3]])
		j++
	}
	return buf.String()
}

func RawWidthStr(str string) string {
	for strings.Contains(str, "\033[") {
		str = rawWidthStr0(str)
	}
	return str
}

func WidthOf(r rune) int {
	// for Chinese punctuation, take up two bytes
	if unicode.IsOneOf(doubleWidth, r) || IsPunctuation(r) {
		return 2
	}
	return 1
}

func StringToBytes(s string) (b []byte) {
	v := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bs := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bs.Data = v.Data
	bs.Len = v.Len
	bs.Cap = v.Len
	return
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Split(value string) (string, int, error) {
	vals := strings.Split(value, ":")
	if len(vals) != 2 {
		return "", 0, fmt.Errorf("invalid format: %q", value)
	}
	if port, err := StrToInt(vals[1]); err != nil {
		return "", 0, err
	} else {
		return vals[0], port, nil
	}
}
