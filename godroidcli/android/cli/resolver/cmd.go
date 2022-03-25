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

package resolver

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"

	"github.com/josexy/godroidcli/util"
)

type ExecCmd interface {
	Command()
	CommandReadAll() ([]byte, error)
	CommandReadLines() ([]string, error)
	GetContext() string
	SetArgs(...string)
	SetCmdArgs(string, ...string)
}

type Cmd struct {
	cmd *exec.Cmd
}

func NewCmd() *Cmd {
	return new(Cmd)
}

func (c *Cmd) Command() {
	_ = c.cmd.Run()
}

func (c *Cmd) CommandReadAll() (data []byte, err error) {
	return c.cmd.CombinedOutput()
}

// CommandReadLines split the output into multiple lines
func (c *Cmd) CommandReadLines() (lines []string, err error) {
	r1, _ := c.cmd.StdoutPipe()
	r2, _ := c.cmd.StderrPipe()

	rds := [2]io.ReadCloser{r1, r2}
	i := 0
	defer func() {
		_ = r1.Close()
		_ = r2.Close()
	}()
	err = c.cmd.Start()
	if err != nil {
		return nil, err
	}

	var line []byte
	var list []string
	reader := bufio.NewReader(rds[i])
	for {
		line, _, err = reader.ReadLine()
		if err != nil || line == nil {
			if err == io.EOF {
				i++
				if i > 1 {
					break
				} else {
					reader = bufio.NewReader(rds[i])
				}
			}
		}
		if len(line) == 0 {
			continue
		}

		list = append(list, util.BytesToString(bytes.TrimSpace(line)))
	}
	// wait for all subprocesses to exit
	if err = c.cmd.Wait(); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Cmd) SetArgs(v ...string) {
	c.SetCmdArgs(v[0], v[1:]...)
}

func (c *Cmd) SetCmdArgs(name string, v ...string) {
	c.cmd = exec.Command(name, v...)
}

func (c *Cmd) GetContext() string {
	return ""
}
