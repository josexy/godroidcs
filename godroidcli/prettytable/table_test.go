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
	"fmt"
	"sync"
	"testing"

	"github.com/fatih/color"
)

func TestNewTable(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			TestPrettyTable_SaveCSVFile(t)
		}()
	}

	wg.Wait()
}

func testLoadFile(filename string, pt *PrettyTable, pipeline RedirectExtra, t *testing.T) {
	pt.BindRedirect(pipeline)
	if n, err := pt.ReadFromFile(filename); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(n)
		fmt.Println(pt)
	}
}

func testSaveFile(filename string, pt *PrettyTable, pipeline RedirectExtra, t *testing.T) {
	pt.SetHeader(Header{color.RedString("username"), "password", color.GreenString("nickname"), "message"})
	pt.AddRow(Row{"root", "12345", "admin", "i am administrator"})
	pt.AddRow(Row{"guest", "12345", color.GreenString("guest"), "hello " + color.BlueString("world")})
	pt.AddRow(Row{"user1", color.BlueString("1234") + "56789" + color.RedString("00"), "user1", "hello world!!!"})
	pt.AddRow(Row{color.YellowString("user2"), "123456789", "user2", "hello!!!"})
	pt.Print()

	pt.BindRedirect(pipeline)
	if err := pt.SaveToFile(filename); err != nil {
		t.Fatal(err)
	}
}

func TestTableAlign(t *testing.T) {
	pt := NewTable()

	pe := NewCSV(pt, true)
	pe.ReadFromFile("test_data/test.csv")

	// left
	pt.SetAlign(Left)
	pt.Print()

	// center
	pt.SetAlign(Center)
	pt.reset()
	pt.Print()

	// right
	pt.SetAlign(Right)
	pt.reset()
	pt.Print()
}

func TestPrettyTable_SaveCSVFile(t *testing.T) {
	pt := NewTable()
	testSaveFile("test_data/test.csv", pt, NewCSV(pt, true), t)
	pt.Grep([]byte("user")).Print()
}

func TestPrettyTable_LoadCSVFile(t *testing.T) {
	pt := NewTable()
	testLoadFile("test_data/test.csv", pt, NewCSV(pt, true), t)
	pt.Grep([]byte("hello")).Grep([]byte("user")).Print()
}

func TestPrettyTable_SaveMarkdownFile(t *testing.T) {
	pt := NewTable()
	testSaveFile("test_data/test.md", pt, NewMarkdown(pt), t)
}

func TestPrettyTable_LoadMarkdownFile(t *testing.T) {
	pt := NewTable()
	testLoadFile("test_data/test.md", pt, NewMarkdown(pt), t)
}

func TestPrettyTable_SaveHtmlFile(t *testing.T) {
	pt := NewTable()
	testSaveFile("test_data/test.html", pt, NewHtml(pt, true, true), t)
}

func TestPrettyTable_LoadHtmlFile(t *testing.T) {
	pt := NewTable()
	testLoadFile("test_data/test.html", pt, NewHtml(pt, true, false), t)
}

func TestPrettyTable_VerticalSaveFile(t *testing.T) {
	pt := NewTable()
	pt.ShowBorder(false)
	pt.SetAlign(Center)
	pt.AddRow(Row{color.RedString(color.GreenString("name")), "mike", "hello!"})
	pt.AddRow(Row{"age", "12", "test..."})
	pt.AddRow(Row{"phone", "23", color.GreenString("hi~")})
	pt.Print()

	pipeline := NewCSV(pt, false)
	if err := pipeline.SaveToFile("test_data/vtest.csv"); err != nil {
		t.Fatal(err)
	}
}

func TestPrettyTable_VerticalLoadFile(t *testing.T) {
	pt := NewTable()
	pipeline := NewCSV(pt, false)
	pt.BindRedirect(pipeline)
	if _, err := pipeline.ReadFromFile("test_data/vtest.csv"); err != nil {
		t.Fatal(err)
	}
	pt.Print()
}

func BenchmarkToString(b *testing.B) {
	pt := NewTable()
	pt.SetHeader(Header{"name", "age"})
	pt.SetAlign(Left)
	for i := 0; i < b.N; i++ {
		pt.Clear()
		pt.AddRow(Row{"hello", "世界"})
		pt.AddRow(Row{"你好", "world"})
		_ = pt.String()
	}
}
