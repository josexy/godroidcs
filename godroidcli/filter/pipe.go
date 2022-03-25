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
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/josexy/godroidcli/filter/bm"
	"github.com/josexy/godroidcli/util"
)

/*
- Pipe		"|"
- Grep		"grep"
- IGrep		"igrep"
- Print		"print"
- Export	"export"
	- csv
	- html
	- markdown
	- "output file"
*/

const (
	Grep     = "grep"
	IGrep    = "igrep"
	Print    = "print"
	Export   = "export"
	Csv      = "csv"
	Markdown = "markdown"
	Html     = "html"
)

var emptyPipeStream = PipeStream{}

type IPipeStream interface {
	Grep([]byte) IPipeStream
	// like Grep, but case insensitive
	IGrep([]byte) IPipeStream
	Export(string, string) IPipeStream
	Print() IPipeStream
	PrintColor(int) IPipeStream
	Filter(*Node) IPipeStream
	Bytes() []byte
}

type (
	streamLine struct {
		// begin <= i <= j <= end
		// position of the line
		begin, end int
		// position of the first matched word in line
		i, j int
	}

	PipeStream struct {
		data []byte
		// all matched lines
		sl []streamLine
	}
)

func CreatePipe(data []byte) *PipeStream {
	return &PipeStream{
		data: data,
	}
}

// PipeOutput output final result with color
func PipeOutput(text []byte, node *Node) {
	CreatePipe(text).Filter(node).(*PipeStream).PrintColor(util.RedAttr)
}

func (sp *PipeStream) Bytes() []byte {
	return sp.data
}

func (sp *PipeStream) Len() int {
	return len(sp.data)
}

// Filter combine grep, igrep, export and print command
func (sp *PipeStream) Filter(node *Node) IPipeStream {
	var pipe IPipeStream = sp
	for {
		if node == nil || !node.IsPipe() {
			break
		}
		if node.HasLeft() {
			list := node.Left.Group()
			// grep, igrep, export, print
			switch list[0] {
			case Grep:
				for _, s := range list[1:] {
					pipe = pipe.Grep(util.StringToBytes(s))
				}
			case IGrep:
				for _, s := range list[1:] {
					pipe = pipe.IGrep(util.StringToBytes(s))
				}
			case Export:
				// export to file directly
				for i := 1; i < len(list); i++ {
					pipe = pipe.Export(list[i], "")
				}
			case Print:
				pipe = pipe.Print()
			default:
				// ignored panic
			}
		}

		if node.HasRight() && node.Right.IsPipe() {
			node = node.Right
		} else {
			break
		}
	}
	return pipe
}

// buildPipeStream construct a new stream pipeline from the current parsed result
func (sp *PipeStream) buildPipeStream() *PipeStream {
	buf := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buf)
	n := len(sp.sl)
	nsl := make([]streamLine, 0, 32)

	lastbegin := 0
	for i := 0; i < n; i++ {
		// write line
		line := sp.Bytes()[sp.sl[i].begin:sp.sl[i].end]
		_, _ = w.Write(line)

		ll := len(line)
		// construct a new stream line which associated with buf.Bytes()
		nsl = append(nsl, streamLine{
			begin: lastbegin,
			end:   lastbegin + ll,
			i:     sp.sl[i].i,
			j:     sp.sl[i].j,
		})
		if i+1 < n {
			_ = w.WriteByte('\n')
		}
		// +1: newline '\n'
		lastbegin += ll + 1
	}
	_ = w.Flush()
	pipe := CreatePipe(buf.Bytes())
	pipe.sl = nsl
	return pipe
}

// Print print without color
func (sp *PipeStream) Print() IPipeStream {
	if len(sp.Bytes()) > 0 {
		_, _ = fmt.Fprintln(util.StdOutput, util.BytesToString(sp.Bytes()))
	}
	return sp
}

// PrintColor print with color
func (sp *PipeStream) PrintColor(color int) IPipeStream {
	n := len(sp.sl)
	if n == 0 {
		return sp.Print()
	}
	for i := 0; i < n; i++ {
		line := sp.Bytes()[sp.sl[i].begin:sp.sl[i].end]
		// replace all matched words in line
		matched := line[sp.sl[i].i:sp.sl[i].j]

		cs := util.ColorMapFunc[color]("%s", util.BytesToString(matched))
		newline := bytes.ReplaceAll(line, matched, util.StringToBytes(cs))
		fmt.Fprintln(util.StdOutput, util.BytesToString(newline))
	}
	return sp
}

// findNewLineIndex find forward and backward index from the current matched index.
// if the backworkd index can be found(prior \n),
// the first value is returned, otherwise returns -1
// if the forward index can be found(next \n),
// the second value is returned, otherwise returns len(data)
func (sp *PipeStream) findNewLineIndex(index int) (int, int) {
	// i: backward index
	// j: forward index
	i, j := index, index
	var doneI, doneJ bool

	for {
		if (doneI && doneJ) || (i < 0 || j >= sp.Len()) {
			break
		}
		// backward index
		for ; i >= 0; i-- {
			if sp.Bytes()[i] == '\n' {
				i++
				doneI = true
				break
			}
		}
		// forward index
		for n := sp.Len(); j < n; j++ {
			if sp.Bytes()[j] == '\n' {
				doneJ = true
				break
			}
		}
	}
	return i, j
}

func (sp *PipeStream) grep0(pattern []byte, caseInsensitive bool) IPipeStream {
	if len(pattern) == 0 {
		return &emptyPipeStream
	}
	indexes := bm.BoyerMooreSearch(sp.Bytes(), pattern, false, caseInsensitive)
	if len(indexes) == 0 {
		return &emptyPipeStream
	}
	// clear
	sp.sl = sp.sl[:0]

	ob, oe, m := -1, -1, len(pattern)
	for _, index := range indexes {
		i, j := index, index+m

		if i >= ob && j <= oe {
			continue
		}

		begin, end := sp.findNewLineIndex(index)
		ob, oe = begin, end
		if begin == -1 {
			begin++
		}

		// the range of new line, including:
		// begin and end index of new line in all text
		// first and last index of the first matched word in new line
		if begin >= 0 && end <= sp.Len() {
			sp.sl = append(sp.sl, streamLine{
				begin: begin,
				end:   end,
				i:     i - begin,
				j:     j - begin,
			})
		}
	}
	return sp.buildPipeStream()
}

// Grep match text with pattern through BM algorithm
func (sp *PipeStream) Grep(pattern []byte) IPipeStream {
	return sp.grep0(pattern, true)
}

// IGrep match text like Grep, but it can ignore case
func (sp *PipeStream) IGrep(pattern []byte) IPipeStream {
	return sp.grep0(pattern, false)
}

// Export export all currently matched lines to a file
// for PipeStream, the second parameter of Export is unused
func (sp *PipeStream) Export(file, unused string) IPipeStream {
	s := util.RawWidthStr(util.BytesToString(sp.Bytes()))
	fp, err := os.Create(file)
	if err != nil {
		return sp
	}
	defer func() { _ = fp.Close() }()
	w := bufio.NewWriter(fp)
	_, _ = w.WriteString(s)
	_ = w.Flush()
	return sp
}
