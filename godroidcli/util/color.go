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

import "github.com/fatih/color"

const (
	GreenAttr = iota
	YellowAttr
	BlueAttr
	RedAttr
	CyanAttr
	WhiteAttr
	HiGreenAttr
	HiYellowAttr
	HiBlueAttr
	HiRedAttr
	HiCyanAttr
	HiWhiteAttr
)

var ColorMapFunc = map[int]func(string, ...interface{}) string{
	GreenAttr:    Green,
	YellowAttr:   Yellow,
	BlueAttr:     Blue,
	RedAttr:      Red,
	CyanAttr:     Cyan,
	WhiteAttr:    White,
	HiGreenAttr:  HiGreen,
	HiYellowAttr: HiYellow,
	HiBlueAttr:   HiBlue,
	HiRedAttr:    HiRed,
	HiCyanAttr:   HiCyan,
	HiWhiteAttr:  HiWhite,
}

var ColorHighMap = map[color.Attribute]*color.Color{
	color.FgHiGreen:  color.New(color.FgHiGreen, color.Bold),
	color.FgHiYellow: color.New(color.FgHiYellow, color.Bold),
	color.FgHiBlue:   color.New(color.FgHiBlue, color.Bold),
	color.FgHiRed:    color.New(color.FgHiRed, color.Bold),
	color.FgHiCyan:   color.New(color.FgHiCyan, color.Bold),
	color.FgHiWhite:  color.New(color.FgHiWhite, color.Bold),
}

func Green(format string, a ...interface{}) string {
	return color.GreenString(format, a...)
}

func Yellow(format string, a ...interface{}) string {
	return color.YellowString(format, a...)
}

func Blue(format string, a ...interface{}) string {
	return color.BlueString(format, a...)
}

func Red(format string, a ...interface{}) string {
	return color.RedString(format, a...)
}

func Cyan(format string, a ...interface{}) string {
	return color.CyanString(format, a...)
}

func White(format string, a ...interface{}) string {
	return color.WhiteString(format, a...)
}

func HiGreen(format string, a ...interface{}) string {
	return ColorHighMap[color.FgHiGreen].Sprintf(format, a...)
}

func HiYellow(format string, a ...interface{}) string {
	return ColorHighMap[color.FgHiYellow].Sprintf(format, a...)
}

func HiBlue(format string, a ...interface{}) string {
	return ColorHighMap[color.FgHiBlue].Sprintf(format, a...)
}

func HiRed(format string, a ...interface{}) string {
	return ColorHighMap[color.FgHiRed].Sprintf(format, a...)
}

func HiCyan(format string, a ...interface{}) string {
	return ColorHighMap[color.FgHiCyan].Sprintf(format, a...)
}

func HiWhite(format string, a ...interface{}) string {
	return ColorHighMap[color.FgHiWhite].Sprintf(format, a...)
}
