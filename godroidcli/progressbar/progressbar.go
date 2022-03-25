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

package progressbar

import (
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/josexy/godroidcli/util"
)

var waitForTime = 65 * time.Millisecond

// ProgressBar emulate download progressbar
// format: [FILENAME/130.44MB] [>>>           ] (29%, 3.17 MB/s) [00h:00m:15s]
type ProgressBar struct {
	filename    string
	totalSize   string
	delta       int64
	prePresent  int64
	rate        string
	start       time.Time
	lastTime    time.Time
	lastUpdated time.Time
}

func New(filename string) *ProgressBar {
	return &ProgressBar{
		filename:    filename,
		rate:        "0.00 B/s",
		start:       time.Now(),
		lastTime:    time.Now(),
		lastUpdated: time.Now(),
	}
}

func (pb *ProgressBar) Update(present, total int64) {
	if present <= 0 || total <= 0 {
		return
	}
	pb.delta += present - pb.prePresent
	pb.prePresent = present

	// delay
	if present < total && time.Since(pb.lastUpdated).Nanoseconds() < waitForTime.Nanoseconds() {
		return
	}
	// earse current line and move to the beginning of line
	io.WriteString(util.StdOutput, "\033[2K\r")

	if pb.totalSize == "" {
		pb.totalSize = util.CalcFileBytes(total)
	}

	// render a line
	s := fmt.Sprintf("[%s/%s] [%s] (%s, %s) [%s]",
		util.HiGreen(pb.filename),
		util.HiYellow(pb.totalSize),
		pb.calcTransferProgress(present, total),
		pb.calcProgressPercent(present, total),
		util.HiRed(pb.calcTransferRate(present, total)),
		pb.calcTransferTime(),
	)
	io.WriteString(util.StdOutput, s)

	pb.lastUpdated = time.Now()
	// print a new line when progressbar completed
	if pb.isFinish(present, total) {
		pb.finish()
	}
}

func (pb *ProgressBar) isFinish(present, total int64) bool {
	return present == total
}

func (pb *ProgressBar) finish() {
	io.WriteString(util.StdOutput, "\n")
	absPath, err := filepath.Abs(pb.filename)
	if err == nil {
		util.Info(fmt.Sprintf("save to file: %s", util.HiRed(absPath)))
	}
}

func (pb *ProgressBar) calcProgressPercent(present, total int64) string {
	return fmt.Sprintf("%d%%", int(100*(float32(present)/float32(total))))
}

func (pb *ProgressBar) split(size int64) (a float64, b float64, c string) {
	v := float64(size)
	f := 1024.0
	t := 0
	for v >= f {
		v /= f
		if t >= 3 {
			break
		}
		t++
	}

	switch t {
	case 0:
		c = "B"
	case 1:
		c = "KB"
	case 2:
		c = "MB"
	default:
		c = "GB"
	}
	a, b = math.Modf(v)
	return
}

func (pb *ProgressBar) calcTransferRate(present, total int64) string {
	now := time.Now()
	if int(now.Sub(pb.lastTime).Seconds()) >= 1 {
		a, b, c := pb.split(pb.delta)
		pb.lastTime = now
		pb.rate = fmt.Sprintf("%d.%02d %s/s", int(a), int(math.Round(b*100)), c)
		pb.delta = 0
	}
	return pb.rate
}

func (pb *ProgressBar) calcTransferTime() string {
	t := int(time.Since(pb.start).Seconds())
	hour := t / 60 / 60 % 60
	minute := t / 60 % 60
	second := t % 60
	return fmt.Sprintf("%02dh:%02dm:%02ds", hour, minute, second)
}

func (pb *ProgressBar) calcTransferProgress(present, total int64) string {
	width, err := util.GetTerminalWidth()
	if err != nil {
		return strings.Repeat(" ", 20)
	}
	width -= len(pb.filename)
	width -= len(pb.totalSize)
	width -= 50
	num := int((float32(present) / float32(total)) * float32(width))
	if width < 0 {
		return " "
	}
	return fmt.Sprintf("%s%s", util.Green(strings.Repeat("█", num)), strings.Repeat(" ", width-num))
}
