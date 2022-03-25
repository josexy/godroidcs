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
	"encoding/csv"
	"io"

	"github.com/josexy/godroidcli/util"
)

type CsvContext struct {
	*RedirectContext
	table       *PrettyTable
	firstHeader bool
}

func NewCSV(table *PrettyTable, firstHeader bool) *CsvContext {
	if table == nil {
		return nil
	}
	ctx := &CsvContext{
		table:           table,
		firstHeader:     firstHeader,
		RedirectContext: &RedirectContext{},
	}
	ctx.BaseRedirect = ctx
	return ctx
}

func (c *CsvContext) ReadFrom(r io.Reader) (int64, error) {
	c.table.mu.Lock()
	defer c.table.mu.Unlock()

	var err error
	reader := csv.NewReader(r)
	var record []string
	var h bool
	c.table.Clear()
	var lines int64
	for {
		record, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			continue
		}
		if record != nil {
			if !h && c.firstHeader {
				c.table.SetHeader(record)
				h = true
			} else {
				c.table.AddRow(record)
			}
			lines++
		}
	}
	return lines, nil
}

func (c *CsvContext) SaveTo(w io.Writer) (err error) {
	c.table.mu.Lock()
	defer c.table.mu.Unlock()

	writer := csv.NewWriter(w)
	var line []string
	for i := range c.table.header {
		line = append(line, util.RawWidthStr(c.table.header[i]))
	}
	if len(line) > 0 {
		if err = writer.Write(line); err != nil {
			return err
		}
	}
	line = line[:0]

	for i := range c.table.rows {
		for j := range c.table.rows[i] {
			line = append(line, util.RawWidthStr(c.table.rows[i][j]))
		}
		if err = writer.Write(line); err != nil {
			return err
		}
		line = line[:0]
	}
	writer.Flush()
	return nil
}
