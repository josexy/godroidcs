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
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"unicode/utf8"

	"github.com/josexy/godroidcli/filter"
	"github.com/josexy/godroidcli/filter/bm"
	"github.com/josexy/godroidcli/util"
)

const (
	Left AlignType = iota
	Right
	Center
)

const (
	Space   = ' '
	NewLine = '\n'
)

var alignFunc = alignFuncMap{
	Left:   leftAlign,
	Center: centerAlign,
	Right:  rightAlign,
}

type (
	AlignType    int
	Header       []string
	Row          []string
	Rows         []Row
	alignFuncMap map[AlignType]func(string, int) string
)

type BorderStyle struct {
	H rune
	V rune

	LeftTop  rune
	Top      rune
	RightTop rune

	Left  rune
	In    rune
	Right rune

	LeftBottom  rune
	Bottom      rune
	RightBottom rune

	showBorder bool
}

type PrettyTable struct {
	header Header
	rows   Rows
	align  AlignType
	buf    []byte
	style  BorderStyle
	writer *bytes.Buffer
	mu     sync.Mutex
	RedirectExtra
}

func NewTable() *PrettyTable {
	if alignFunc == nil {
		alignFunc = make(alignFuncMap)
	}
	pt := &PrettyTable{
		align: Left,
		rows:  make([]Row, 0, 64),
		buf:   make([]byte, 0, 4096),
		style: BorderStyle{
			H: util.H, V: util.V,
			LeftTop: util.PLUS, Top: util.PLUS, RightTop: util.PLUS,
			Left: util.PLUS, In: util.PLUS, Right: util.PLUS,
			LeftBottom: util.PLUS, Bottom: util.PLUS, RightBottom: util.PLUS,
			showBorder: true,
		},
	}

	pt.writer = bytes.NewBuffer(pt.buf)
	return pt
}

func match(str string, pattern []byte, caseInsensitive bool) bool {
	return len(bm.BoyerMooreSearch(util.StringToBytes(str), pattern, true, caseInsensitive)) > 0
}

func (p *PrettyTable) grep0(pattern []byte, caseInsensitive bool) filter.IPipeStream {
	tb := NewTable()
	tb.SetHeader(p.header)
	tb.align = p.align
	tb.style = p.style

	if len(pattern) == 0 {
		return nil
	}

	var mat bool
	for i := 0; i < len(p.rows); i++ {
		mat = false
		for j := 0; j < len(p.rows[i]); j++ {
			// use BM algorithm
			if match(p.rows[i][j], pattern, caseInsensitive) {
				mat = true
				break
			}
		}
		if mat {
			tb.AddRow(p.rows[i])
		}
	}
	return tb
}

// Grep filter some rows of table by pattern and return PipeStream
// for example:
// | grep hello			# filter rows by pattern "hello"
// | grep "hello world"	# filter rows by pattern "hello world"
// | grep hello world	# filter rows by pattern "hello" at first and then "world"
func (p *PrettyTable) Grep(pattern []byte) filter.IPipeStream {
	return p.grep0(pattern, true)
}

func (p *PrettyTable) IGrep(pattern []byte) filter.IPipeStream {
	return p.grep0(pattern, false)
}

// Export write table to file by given format type
// if the format is given but no output file is provided,
// the output file name is the same as the format.
// supported format: csv, markdown, html
// for csv: use CsvContext
// for markdown: use MdContext
// for html: use HtmlContext
// for example:
// ... | export 0.txt 1.txt csv 2.csv html 3.html 4.txt csv | ...
// - common file: 0.txt, 1.txt, 4.txt
// - csv file: 2.csv, csv
// - html file: 3.html
func (p *PrettyTable) Export(t1 string, t2 string) filter.IPipeStream {
	var redirect RedirectExtra
	switch t1 {
	case filter.Csv: // csv
		redirect = NewCSV(p, true)
	case filter.Markdown: // markdown
		if len(p.header) == 0 {
			util.Warn("table has no header, can not export markdown format")
			break
		}
		redirect = NewMarkdown(p)
	case filter.Html: // html
		redirect = NewHtml(p, true, true)
	default:
	}

	if redirect == nil {
		_ = ioutil.WriteFile(t1, util.StringToBytes(util.RawWidthStr(p.String())), 0644)
	} else {
		_ = redirect.SaveToFile(t2)
	}
	return p
}

// Print write table to standard output
func (p *PrettyTable) Print() filter.IPipeStream {
	_, _ = fmt.Fprintln(util.StdOutput, p.String())
	return p
}

// PrintColor ignored
func (p *PrettyTable) PrintColor(int) filter.IPipeStream {
	return p.Print()
}

// Filter handle Grep, Export and Print
func (p *PrettyTable) Filter(node *filter.Node) filter.IPipeStream {
	var tb filter.IPipeStream = p
	for {
		if node == nil || !node.IsPipe() {
			break
		}
		if node.HasLeft() {
			list := node.Left.Group()

			// grep, export, print
			switch list[0] {
			case filter.Grep:
				for _, s := range list[1:] {
					tb = tb.Grep(util.StringToBytes(s))
				}
			case filter.IGrep:
				for _, s := range list[1:] {
					tb = tb.IGrep(util.StringToBytes(s))
				}
			case filter.Export:
				// | export 1.txt 2.txt html 3.txt
				// | export csv 2.txt html 3.txt
				for i := 1; i < len(list); {
					var t2 string
					t1 := list[i]
					// type t1
					if t1 == filter.Csv || t1 == filter.Html || t1 == filter.Markdown {
						// output file t2
						if i+1 < len(list) {
							t2 = list[i+1]
							i++
						} else {
							// t1 = t2
							t2 = t1
						}
					}
					i++
					// the output file name is t2
					tb = tb.Export(t1, t2)
				}
			case filter.Print:
				tb = tb.Print()
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
	return tb
}

// BindRedirect bind redirect interface,
// so user can call the methods of RedirectExtra interface directly
func (p *PrettyTable) BindRedirect(redirect RedirectExtra) {
	if redirect == nil {
		panic(errors.New("bind redirect interface failed"))
	}
	p.RedirectExtra = redirect
}

func (p *PrettyTable) ShowBorder(show bool) {
	p.style.showBorder = show
}

func (p *PrettyTable) SetBorderStyle(style BorderStyle) {
	p.style.showBorder = style.showBorder
	if style.LeftTop != 0 {
		p.style.LeftTop = style.LeftTop
	}
	if style.Top != 0 {
		p.style.Top = style.Top
	}
	if style.RightTop != 0 {
		p.style.RightTop = style.RightTop
	}
	if style.Left != 0 {
		p.style.Left = style.Left
	}
	if style.In != 0 {
		p.style.In = style.In
	}
	if style.Right != 0 {
		p.style.Right = style.Right
	}
	if style.LeftBottom != 0 {
		p.style.LeftBottom = style.LeftBottom
	}
	if style.Bottom != 0 {
		p.style.Bottom = style.Bottom
	}
	if style.RightBottom != 0 {
		p.style.RightBottom = style.RightBottom
	}
}

func (p *PrettyTable) SetAlign(align AlignType) {
	p.align = align
}

func (p *PrettyTable) SetHeader(header Header) {
	p.header = header
}

func (p *PrettyTable) AddRow(row Row) {
	p.rows = append(p.rows, row)
}

func (p *PrettyTable) AddRows(rows Rows) {
	p.rows = append(p.rows, rows...)
}

func (p *PrettyTable) GetRow(index int) (row Row) {
	if index >= 0 && index < len(p.rows) {
		row = p.rows[index]
	}
	return
}

func (p *PrettyTable) wideOfSlice(slice []string) (wides []int) {
	for i := range slice {
		wides = append(wides, widthOfStr(slice[i]))
	}
	return
}

func (p *PrettyTable) calcMaxWide(s1, s2 []int) {
	n1, n2 := len(s1), len(s2)
	if n2 < n1 {
		n1 = n2
	}
	for i := 0; i < n1; i++ {
		if s1[i] < s2[i] {
			s1[i] = s2[i]
		}
	}
}

func (p *PrettyTable) maxWideOfTable() (max []int) {
	max = p.wideOfSlice(p.header)
	for i := 0; i < len(p.rows); i++ {
		rowWides := p.wideOfSlice(p.rows[i])
		if max == nil {
			max = rowWides
		} else {
			p.calcMaxWide(max, rowWides)
		}
	}
	return
}

func (p *PrettyTable) paddingItem(str string, max int) string {
	return alignFunc[p.align](str, max)
}

func (p *PrettyTable) drawLineBorder(l, m, r rune, wides []int, newline bool) []byte {
	n := len(wides)
	line := bytes.NewBuffer(nil)
	wfn := func(l0 rune, wide int) {
		line.WriteRune(l0)
		for i := 0; i < wide+2; i++ {
			line.WriteRune(p.style.H)
		}
	}
	for i := 0; i < n; i++ {
		var l0 rune
		if i == 0 {
			l0 = l
		} else {
			l0 = m
		}
		wfn(l0, wides[i])
	}
	line.WriteRune(r)
	if newline {
		line.WriteByte(NewLine)
	}
	return line.Bytes()
}

func (p *PrettyTable) drawLineItems(items []string, wides []int) {
	n := len(items)
	m := len(p.header)
	if m == 0 {
		m = len(wides)
	}

	for i := 0; i < m; i++ {
		if p.style.showBorder {
			p.writer.WriteRune(p.style.V)
		}
		p.writer.WriteByte(Space)
		var item string
		if i+1 > n {
			item = ""
		} else {
			item = items[i]
		}
		p.writer.WriteString(p.paddingItem(item, wides[i]))
		p.writer.WriteByte(Space)
	}
	if p.style.showBorder {
		p.writer.WriteRune(p.style.V)
	}
	p.writer.WriteByte(NewLine)
}

func (p *PrettyTable) drawHeader(wides []int) {
	if p.style.showBorder {
		p.writer.Write(p.drawLineBorder(p.style.LeftTop, p.style.Top, p.style.RightTop, wides, true))
	}
	p.drawLineItems(p.header, wides)
}

func (p *PrettyTable) drawRows(wides []int) {
	n := len(p.rows)
	for i := 0; i < n; i++ {
		p.drawLineItems(p.rows[i], wides)
	}
	if p.style.showBorder {
		p.writer.Write(p.drawLineBorder(p.style.LeftBottom, p.style.Bottom, p.style.RightBottom, wides, false))
	}
}

func (p *PrettyTable) Update() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.reset()
	p.draw()
}

func (p *PrettyTable) draw() bool {
	maxWides := p.maxWideOfTable()
	if maxWides == nil {
		return false
	}
	n1, n2 := len(p.header), len(p.rows)

	if n1 > 0 {
		p.drawHeader(maxWides)
	}
	if p.style.showBorder && (n1 > 0 || n2 > 0) {
		p.writer.Write(p.drawLineBorder(p.style.Left, p.style.In, p.style.Right, maxWides, true))
	}
	if n2 > 0 {
		p.drawRows(maxWides)
	}
	return true
}

func (p *PrettyTable) reset() {
	p.writer.Reset()
}

// Clear clear header and all rows
func (p *PrettyTable) Clear() {
	p.reset()
	p.header = p.header[:0]
	p.rows = p.rows[:0]
}

func (p *PrettyTable) String() string {
	data := p.Bytes()
	if data != nil {
		return util.BytesToString(data)
	}
	return ""
}

func (p *PrettyTable) Bytes() []byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.writer.Len() > 0 {
		return p.writer.Bytes()
	}
	if p.draw() {
		return p.writer.Bytes()
	}
	return nil
}

func widthOfStr(str string) (num int) {
	str = util.RawWidthStr(str)
	for _, r := range str {
		num += util.WidthOf(r)
	}
	return
}

func leftAlign(str string, max int) string {
	n := utf8.RuneCountInString(str)
	wide := widthOfStr(str) - n
	return fmt.Sprintf("%-*s", max-wide, str)
}

func rightAlign(str string, max int) string {
	n := utf8.RuneCountInString(str)
	wide := widthOfStr(str) - n
	return fmt.Sprintf("%*s", max-wide, str)
}

func centerAlign(str string, max int) string {
	n := utf8.RuneCountInString(str)
	wide := widthOfStr(str)
	padding := (max-wide)/2 + n
	pr := max - padding - (wide - n)
	return fmt.Sprintf("%*s%*s", padding, str, pr, "")
}
