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
	"bytes"
	"container/list"
	"fmt"
	"strings"
	"unicode"

	"github.com/josexy/godroidcli/status"
	"github.com/josexy/godroidcli/util"
)

/*
cmd fs list . | grep "xx" ｜ export csv yy
	 cmd
	/   \
   fs	 |
  / 	/  \
 list  grep export
 /    /		 /
 .   "xx"   csv
			/
		  yy

cmd pm size "xx"|grep "yy"|export data|grep "zz"|export o.txt|print
		cmd
		/  \
	  pm	|
	 /	   /  \
   size grep    |
  /		/	 / 	  \
 "xx"  "yy" export grep
			/	   /  \
		 data	 "zz"  |
					 /   \
				 export	  |
				  /		 /
				o.txt   print
*/

type tokenType byte

const (
	tokenString tokenType = iota + 1
	tokenPipe
)

type token struct {
	name []byte
	typ  tokenType
}

func NewToken(name []byte, typ tokenType) *token {
	return &token{name: name, typ: typ}
}

func (t token) String() string {
	return string(t.name)
}

type Param struct {
	Node *Node
	Args []string
}

var EmptyParam Param

type Node struct {
	tok   *token
	Left  *Node
	Right *Node
}

func (node *Node) IsPipe() bool {
	return node.tok.typ == tokenPipe
}

func (node *Node) String() string {
	return "[" + node.tok.String() + "]"
}

func (node *Node) HasLeft() bool {
	return node.Left != nil
}

func (node *Node) HasRight() bool {
	return node.Right != nil
}

// Group find all string by left child nodes
func (node *Node) Group() (list []string) {
	n := node
	for n != nil {
		list = append(list, n.tok.String())
		n = n.Left
	}
	return
}

const (
	startTagState = iota + 1
	startStringState
)

type CmdParser struct {
	reader    *bytes.Reader
	i         int
	n         int
	state     int
	root      *Node
	tokes     []token
	peekTokes *list.List
}

func NewCmdParser() *CmdParser {
	return &CmdParser{peekTokes: list.New()}
}

func (cmd *CmdParser) Root() *Node {
	return cmd.root
}

func (cmd *CmdParser) readByte() byte {
	b, _ := cmd.reader.ReadByte()
	cmd.i++
	return b
}

func (cmd *CmdParser) unread() {
	_ = cmd.reader.UnreadByte()
	cmd.i--
}

func isSpace(b byte) bool {
	return unicode.IsSpace(rune(b))
}

func (cmd *CmdParser) getGetToken() (*token, error) {
	buf := bytes.NewBuffer(nil)
	var wrap bool
	for cmd.i <= cmd.n {
		c := cmd.readByte()
		if c == 0 {
			if buf.Len() > 0 {
				return NewToken(buf.Bytes(), tokenString), nil
			}
			break
		}

		switch cmd.state {
		case startTagState:
			switch {
			case isSpace(c):
				break
			case c == '|':
				cmd.state = startTagState
				buf.WriteByte('|')
				return NewToken(buf.Bytes(), tokenPipe), nil
			default:
				if c == '\'' || c == '"' {
					wrap = true
				} else {
					cmd.unread()
				}
				cmd.state = startStringState
			}
		case startStringState:
			if c == '\'' || c == '"' {
				cmd.state = startTagState
				return NewToken(buf.Bytes(), tokenString), nil
			}
			if !wrap {
				if isSpace(c) || c == '|' {
					if c == '|' {
						cmd.unread()
					}
					cmd.state = startTagState
					return NewToken(buf.Bytes(), tokenString), nil
				}
			}
			buf.WriteByte(c)
		default:
			return nil, nil
		}
	}
	return nil, nil
}

func (cmd *CmdParser) tryGetToken() (*token, error) {
	if cmd.peekTokes.Len() > 0 {
		e := cmd.peekTokes.Front()
		cmd.peekTokes.Remove(e)
		return e.Value.(*token), nil
	} else {
		tok, err := cmd.getGetToken()
		return tok, err
	}
}

func (cmd *CmdParser) pushPeekToken(tok *token) {
	if tok != nil {
		cmd.peekTokes.PushFront(tok)
	}
}

func (cmd *CmdParser) parseCmd() *Node {
	left := cmd.parseString()
	if left == nil {
		panic(status.ErrCannotParseCmd)
	}
	tok, _ := cmd.tryGetToken()
	if tok != nil && tok.typ == tokenPipe {
		cmd.pushPeekToken(tok)
		left.Right = cmd.parsePipe()
	}
	return left
}

func (cmd *CmdParser) parsePipe() *Node {
	var parent *Node
	for {
		tok, _ := cmd.tryGetToken()
		if tok == nil || tok.typ != tokenPipe {
			break
		}
		parent = &Node{tok: tok}
		parent.Left = cmd.parseOp()
		if parent.Left == nil {
			panic(status.ErrProvideParams)
		}
		parent.Right = cmd.parsePipe()
	}
	return parent
}

func (cmd *CmdParser) parseString() *Node {
	var parent *Node
	for {
		tok, _ := cmd.tryGetToken()
		if tok == nil || tok.typ == tokenPipe {
			cmd.pushPeekToken(tok)
			break
		}
		parent = &Node{tok: tok}
		parent.Left = cmd.parseString()
	}
	return parent
}

func (cmd *CmdParser) parseOp() *Node {
	tok, _ := cmd.tryGetToken()
	if tok == nil {
		return nil
	}

	name := string(tok.name)
	left := &Node{tok: tok}
	switch name {
	case Export:
		left.Left = cmd.parseType()
	case IGrep, Grep, Print:
		left.Left = cmd.parseString()
	default:
		panic(fmt.Errorf("can not parse identifier: [%s]", name))
	}
	if name != Print && left.Left == nil {
		panic(status.ErrProvideParams)
	}
	return left
}

func (cmd *CmdParser) parseType() *Node {
	tok, _ := cmd.tryGetToken()
	if tok == nil {
		return nil
	}
	if tok.typ == tokenPipe {
		cmd.pushPeekToken(tok)
		return nil
	}
	name := string(tok.name)
	left := &Node{tok: tok}
	left.Left = cmd.parseString()
	switch name {
	case Csv, Markdown, Html:
		if left.Left == nil {
			panic(status.ErrProvideParams)
		}
	}
	return left
}

func (cmd *CmdParser) Parse(line string) (err error) {
	if len(strings.TrimSpace(line)) == 0 {
		return status.ErrEmptyString
	}
	if cmd.reader == nil {
		cmd.reader = bytes.NewReader([]byte(line))
	} else {
		cmd.reader.Reset([]byte(line))
	}
	cmd.i = 0
	cmd.n = cmd.reader.Len()
	cmd.root = nil
	cmd.state = startTagState
	cmd.tokes = cmd.tokes[:0]
	cmd.peekTokes.Init()

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	cmd.root = cmd.parseCmd()
	// cmd.dumpSyntaxTreeTest()
	return
}

/**
cmd pm package com.android.chrome | grep android | print
┌── [cmd]
├── [pm]
│   └── [package]
│       └── [com.android.chrome]
└── [|]
    ├── [grep]
    │   └── [android]
    └── [|]
        └── [print]
*/
func (cmd *CmdParser) dumpTest(node *Node, v *[]bool, inx, size int) {
	if node == nil {
		return
	}
	if inx == size {
		fmt.Print("┌── ")
	} else if inx+1 == size {
		fmt.Print("└── ")
		(*v)[len(*v)-1] = false
	} else {
		fmt.Print("├── ")
	}
	fmt.Println(util.Green(node.String()))
	n := 0
	if node.HasLeft() {
		n++
	}
	if node.HasRight() {
		n++
	}
	for j := 0; j < n; j++ {
		for i := 0; i < len(*v); i++ {
			if (*v)[i] {
				fmt.Print("│")
			} else {
				fmt.Print(" ")
			}
			fmt.Print(strings.Repeat(" ", 3))
		}
		var nc *Node
		if j == 0 {
			nc = node.Left
			if nc == nil {
				nc = node.Right
			}
		} else if j == 1 {
			nc = node.Right
		}
		*v = append(*v, n > 1)
		cmd.dumpTest(nc, v, j, n)
		*v = (*v)[:len(*v)-1]
	}
}

func (cmd *CmdParser) dumpSyntaxTreeTest() {
	if cmd.root == nil {
		return
	}
	var v []bool
	cmd.dumpTest(cmd.root, &v, 0, 0)
}
