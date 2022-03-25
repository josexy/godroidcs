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
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/josexy/godroidcli/util"
)

const (
	MdSplit1 = "|"
	MdSplit2 = "-"
)

type MdContext struct {
	*RedirectContext
	table *PrettyTable
	rd    *bufio.Reader
	wr    *bufio.Writer
}

func NewMarkdown(table *PrettyTable) *MdContext {
	if table == nil {
		return nil
	}
	ctx := &MdContext{
		table:           table,
		RedirectContext: &RedirectContext{},
	}
	ctx.BaseRedirect = ctx
	return ctx
}

func (m *MdContext) parseLine(line string) (row []string) {
	list := strings.Split(line, MdSplit1)
	for i := range list {
		s := strings.TrimSpace(list[i])
		if len(s) > 0 {
			row = append(row, s)
		}
	}
	return
}

func (m *MdContext) ReadFrom(reader io.Reader) (int64, error) {
	m.table.mu.Lock()
	defer m.table.mu.Unlock()
	m.rd = bufio.NewReader(reader)
	var h bool
	m.table.Clear()
	var lines int64
	for {
		line, _, err := m.rd.ReadLine()
		if err != nil || line == nil {
			if err == io.EOF {
				break
			}
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		lines++
		r := m.parseLine(string(line))
		if !h {
			h = true
			m.table.SetHeader(r)
			// skip line separator
			_, _, _ = m.rd.ReadLine()
		} else {
			m.table.AddRow(r)
		}
	}
	return lines, nil
}

func (m *MdContext) SaveTo(writer io.Writer) (err error) {
	m.table.mu.Lock()
	defer m.table.mu.Unlock()
	n := len(m.table.rows)

	m.wr = bufio.NewWriter(writer)
	_ = m.wr.WriteByte(MdSplit1[0])
	for i := range m.table.header {
		if _, err = m.wr.WriteString(util.RawWidthStr(m.table.header[i])); err != nil {
			return err
		}
		_ = m.wr.WriteByte(MdSplit1[0])
	}
	if n > 0 {
		_ = m.wr.WriteByte(NewLine)
	}

	// separator line
	_ = m.wr.WriteByte(MdSplit1[0])
	for range m.table.header {
		_, _ = m.wr.WriteString(strings.Repeat(MdSplit2, 4))
		_ = m.wr.WriteByte(MdSplit1[0])
	}
	_ = m.wr.WriteByte(NewLine)
	for i := 0; i < n; i++ {
		_ = m.wr.WriteByte(MdSplit1[0])
		for j := range m.table.rows[i] {
			if _, err = m.wr.WriteString(util.RawWidthStr(m.table.rows[i][j])); err != nil {
				return err
			}
			_ = m.wr.WriteByte(MdSplit1[0])
		}
		if i+1 < n {
			_ = m.wr.WriteByte(NewLine)
		}
	}
	err = m.wr.Flush()
	return err
}
