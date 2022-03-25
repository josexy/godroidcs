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

package prettytable

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type BaseRedirect interface {
	ReadFrom(io.Reader) (int64, error)
	SaveTo(io.Writer) error
}

type RedirectExtra interface {
	BaseRedirect
	ReadFromString(string) (int64, error)
	ReadFromBytes([]byte) (int64, error)
	ReadFromFile(string) (int64, error)
	SaveToBytes() ([]byte, error)
	SaveToFile(string) error
}

type RedirectContext struct {
	BaseRedirect
}

func (c *RedirectContext) ReadFromString(s string) (int64, error) {
	return c.ReadFrom(strings.NewReader(s))
}

func (c *RedirectContext) ReadFromBytes(b []byte) (int64, error) {
	return c.ReadFrom(bytes.NewReader(b))
}

func (c *RedirectContext) ReadFromFile(filename string) (int64, error) {
	var fp, err = os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = fp.Close()
	}()
	return c.ReadFrom(fp)
}

func (c *RedirectContext) SaveToBytes() (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	err = c.SaveTo(buf)
	if err != nil {
		return
	}
	data = buf.Bytes()
	return
}

func (c *RedirectContext) SaveToFile(filename string) error {
	var fp, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
	}()
	return c.SaveTo(fp)
}
