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

const (
	PC1  = 0x3002
	PC2  = 0xFF1F
	PC3  = 0xFF01
	PC4  = 0xFF0C
	PC5  = 0x3001
	PC6  = 0xFF1B
	PC7  = 0xFF1A
	PC8  = 0x300C
	PC9  = 0x300D
	PC10 = 0x2018
	PC11 = 0x2019
	PC12 = 0xFF08
	PC13 = 0xFF09
	PC14 = 0x3014
	PC15 = 0x3015
	PC16 = 0x3010
	PC17 = 0x3011
	PC18 = 0x2014
	PC19 = 0x2026
	PC20 = 0x2013
	PC21 = 0xFF0E
	PC22 = 0x300A
	PC23 = 0x300B
	PC24 = 0x3008
	PC25 = 0x3009
)

const (
	PLUS  rune = '+'
	H     rune = '-'
	V     rune = '|'
	DOT   rune = '.'
	EQUAL rune = '='

	SIGN1 rune = '@'
	SIGN2 rune = '#'
	SIGN3 rune = '$'
	SIGN4 rune = '%'
	SIGN5 rune = '\\'
	SIGN6 rune = '*'
	SIGN7 rune = '/'

	SIGN8 rune = '─'
	SIGN9 rune = '│'

	SIGN10 rune = '┌'
	SIGN11 rune = '┬'
	SIGN12 rune = '┐'

	SIGN13 rune = '├'
	SIGN14 rune = '┼'
	SIGN15 rune = '┤'

	SIGN16 rune = '└'
	SIGN17 rune = '┴'
	SIGN18 rune = '┘'
)

func IsPunctuation(r rune) bool {
	return r == PC1 || r == PC2 ||
		r == PC3 || r == PC4 || r == PC5 || r == PC6 ||
		r == PC7 || r == PC8 || r == PC9 || r == PC10 ||
		r == PC11 || r == PC12 || r == PC13 || r == PC14 ||
		r == PC15 || r == PC16 || r == PC17 || r == PC18 ||
		r == PC19 || r == PC20 || r == PC21 || r == PC22 ||
		r == PC23 || r == PC24 || r == PC25
}
