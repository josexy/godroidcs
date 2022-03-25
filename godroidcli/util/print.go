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

package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	StdOutput = color.Output
)

var Logo = `
   ______      ____             _     ___________ 
  / ____/___  / __ \_________  (_)___/ / ____/ (_)
 / / __/ __ \/ / / / ___/ __ \/ / __  / /   / / / 
/ /_/ / /_/ / /_/ / /  / /_/ / / /_/ / /___/ / /  
\____/\____/_____/_/   \____/_/\__,_/\____/_/_/`

func typeOf(level Level) (typ string) {
	switch level {
	case INFO:
		typ = color.GreenString("[INFO]")
	case WARN:
		typ = color.YellowString("[WARN]")
	case ERROR:
		typ = color.RedString("[ERROR]")
	case FATAL:
		typ = color.HiRedString("[FATAL]")
	}
	return
}

func Print(format string, v ...interface{}) {
	fmt.Fprintf(StdOutput, format, v...)
}

func Printf(level Level, format string, v ...interface{}) {
	fmt.Fprintf(StdOutput, "%s %s\n", typeOf(level), fmt.Sprintf(format, v...))
	if level == FATAL {
		os.Exit(0)
	}
}

func Info(format string, v ...interface{}) {
	Printf(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	Printf(WARN, format, v...)
}

func ErrorBy(err error) {
	Error("%s", err.Error())
}

func Error(format string, v ...interface{}) {
	Printf(ERROR, format, v...)
}

func FatalBy(err error) {
	Fatal("%s", err.Error())
}

func Fatal(format string, v ...interface{}) {
	Printf(FATAL, format, v...)
}
