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
	"html"
	"io"
	"strings"
	"unicode"

	"github.com/josexy/godroidcli/status"
	"github.com/josexy/godroidcli/util"
)

type tokenType int

const (
	table tokenType = iota
	row
	header
	cell
	text
)

const (
	stateStartTag = iota
	stateEndTag
	stateStartEle
	stateStartText
)

type token struct {
	name    []byte
	typ     tokenType
	closure bool
}

type node struct {
	tok   token
	child []*node
}

type htmlTableAst struct {
	tokens   []token
	rd       *bufio.Reader
	root     *node
	i        int
	n        int
	state    int
	unescape bool
}

func newHtmlTableAst(data []byte) *htmlTableAst {
	return &htmlTableAst{
		rd:     bufio.NewReader(bytes.NewReader(data)),
		tokens: make([]token, 0, 32),
		n:      len(data),
	}
}

func (t *htmlTableAst) parse() bool {
	for t.i < t.n {
		if tok, ok := t.tryGetToken(); ok {
			t.tokens = append(t.tokens, tok)
		}
	}
	if len(t.tokens) == 0 {
		return false
	}
	t.i = 1
	t.n = len(t.tokens) - 1

	t.root = &node{tok: t.tokens[0]}
	if t.root.tok.typ != table || t.root.tok.typ != t.tokens[t.n].typ || !t.tokens[t.n].closure {
		return false
	}

	return t.generate(t.root)
}

func (t *htmlTableAst) skipByte() {
	_ = t.readByte()
}

func (t *htmlTableAst) readByte() (b byte) {
	b, _ = t.rd.ReadByte()
	t.i++
	return
}

func (t *htmlTableAst) unreadByte() {
	_ = t.rd.UnreadByte()
	t.i--
}

func (t *htmlTableAst) peekByte() (b byte, err error) {
	var bs []byte
	bs, err = t.rd.Peek(1)
	if err != nil {
		return
	}
	return bs[0], err
}

func (t *htmlTableAst) checkTypeOfToken(tok *token) {
	if !bytes.HasPrefix(tok.name, []byte{'t'}) {
		tok.typ = text
		return
	}
	switch strings.ToLower(string(tok.name)) {
	case "table":
		tok.typ = table
	case "tr":
		tok.typ = row
	case "th":
		tok.typ = header
	case "td":
		tok.typ = cell
	default:
		tok.typ = text
	}
}

func (t *htmlTableAst) tryGetToken() (tok token, ok bool) {
	for t.i < t.n {
		b := t.readByte()

		if unicode.IsSpace(rune(b)) {
			continue
		}

		switch t.state {
		case stateStartTag:
			switch b {
			case '<':
				pk, err := t.peekByte()
				if err == nil {
					switch pk {
					case '/': // </xxx>
						t.skipByte() // skip '/'
						tok.closure = true
						t.state = stateEndTag
					default: // <xxx>
						t.state = stateStartEle
					}
				}
			default: // text
				t.state = stateStartText
				tok.name = append(tok.name, b)
			}
		case stateEndTag: // </xxx> <yyyy></yyyy>
			if b == '>' {
				t.checkTypeOfToken(&tok)
				t.state = stateStartTag
				ok = true
				return
			}
			tok.name = append(tok.name, b)
		case stateStartEle: // <xxx>
			if b == '>' {
				// 1. <xxx> </xxx>
				// 2. <xxx> <yyy>
				// 3. <xxx> text
				t.checkTypeOfToken(&tok)
				t.state = stateStartTag
				ok = true
				return
			}
			tok.name = append(tok.name, b)
		case stateStartText: // <xxx> text </xxx>
			if b == '<' {
				t.checkTypeOfToken(&tok)
				t.unreadByte()
				t.state = stateStartTag
				ok = true
				return
			}
			tok.name = append(tok.name, b)
		default:
			return
		}
	}
	return
}

/*
generate a simple syntax tree
<table>
	<tr>
		<th> header </th>
	</tr>
	<tr>
		<td> text </td>
	</tr>
</table>

syntax tree:
	  <table>
	  /      \
   <tr>      <tr>
   /  \ 	  /  \
  <th> <th>  <td> <td>
  /  		 /
 [header] 	[text]
*/
func (t *htmlTableAst) generate(n *node) bool {
	if n.tok.closure && n.tok.typ == text {
		return false
	}
	if n.tok.typ == text || n.tok.closure {
		return true
	}

	for t.i < t.n {
		tok := t.tokens[t.i]
		t.i++

		nc := &node{
			tok:   tok,
			child: make([]*node, 0, 2),
		}
		if !t.generate(nc) {
			return false
		}
		if tok.closure {
			return tok.typ == n.tok.typ
		} else {
			n.child = append(n.child, nc)
		}
	}
	return true
}

func (t *htmlTableAst) get(typ tokenType, i int) (line []string) {
	defer func() {
		if recover() != nil {
			line = nil
		}
	}()
	r := t.root.child[i].child
	for j := range r {
		if r[j].tok.typ != typ {
			return nil
		}
		uns := string(r[j].child[0].tok.name)
		if t.unescape {
			uns = html.UnescapeString(uns)
		}
		line = append(line, uns)
	}
	return
}

type HtmlContext struct {
	*RedirectContext
	table       *PrettyTable
	ast         *htmlTableAst
	wr          *bufio.Writer
	unescape    bool
	firstHeader bool
}

func NewHtml(table *PrettyTable, unescape, firstHeader bool) *HtmlContext {
	if table == nil {
		return nil
	}
	ctx := &HtmlContext{
		table:           table,
		unescape:        unescape,
		firstHeader:     firstHeader,
		RedirectContext: &RedirectContext{},
	}
	ctx.BaseRedirect = ctx
	return ctx
}

func (h *HtmlContext) ReadFrom(reader io.Reader) (int64, error) {
	h.table.mu.Lock()
	defer h.table.mu.Unlock()

	data, err := io.ReadAll(reader)
	if err != nil {
		return -1, err
	}
	h.ast = newHtmlTableAst(data)
	h.ast.unescape = h.unescape
	if !h.ast.parse() {
		return -1, status.ErrCannotParseHtml
	}
	h.table.Clear()
	if hdr := h.ast.get(header, 0); h != nil {
		if h.firstHeader {
			h.table.SetHeader(hdr)
		} else {
			h.table.AddRow(hdr)
		}
	}
	n := len(h.ast.root.child)
	for i := 1; i < n; i++ {
		if r := h.ast.get(cell, i); r != nil {
			h.table.AddRow(r)
		}
	}
	return int64(n), nil
}

func (h *HtmlContext) SaveTo(writer io.Writer) error {
	h.table.mu.Lock()
	defer h.table.mu.Unlock()

	h.wr = bufio.NewWriter(writer)
	_, _ = h.wr.WriteString("<table>")

	_, _ = h.wr.WriteString("<tr>")
	for i := range h.table.header {
		_, _ = h.wr.WriteString("<th>")
		_, _ = h.wr.WriteString(util.RawWidthStr(h.table.header[i]))
		_, _ = h.wr.WriteString("</th>")
	}
	_, _ = h.wr.WriteString("</tr>")

	for i := range h.table.rows {
		_, _ = h.wr.WriteString("<tr>")
		for j := range h.table.rows[i] {
			_, _ = h.wr.WriteString("<td>")
			_, _ = h.wr.WriteString(util.RawWidthStr(h.table.rows[i][j]))
			_, _ = h.wr.WriteString("</td>")
		}
		_, _ = h.wr.WriteString("</tr>")
	}
	_, _ = h.wr.WriteString("</table>")

	_ = h.wr.Flush()
	return nil
}
