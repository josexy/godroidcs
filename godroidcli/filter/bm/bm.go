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

package bm

import (
	"bytes"
	"unicode"
)

var bc [256]int

func calcBadChar(pattern []byte) {
	for i := 0; i < 256; i++ {
		bc[i] = -1
	}
	for i := 0; i < len(pattern); i++ {
		bc[pattern[i]] = i
	}
}

func calcGoodSuffix(pattern []byte, suffix []int, prefix []bool) {
	m := len(pattern)
	for i := 0; i < m; i++ {
		suffix[i] = -1
		prefix[i] = false
	}

	for i := 0; i < m-1; i++ {
		j := i
		k := 0
		for j >= 0 && pattern[j] == pattern[m-1-k] {
			j--
			k++
			suffix[k] = j + 1
		}
		if j == -1 {
			prefix[k] = true
		}
	}
}

func moveByGoodSuffix(j, m int, suffix []int, prefix []bool) int {
	k := m - 1 - j
	// match
	if suffix[k] != -1 {
		return j - suffix[k] + 1
	}
	for r := j + 2; r <= m-1; r++ {
		if prefix[m-r] {
			return r
		}
	}
	return m
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func toLower(b byte) byte {
	return byte(unicode.ToLower(rune(b)))
}

// BoyerMooreSearch find indexes of text
// if once flag is true, BoyerMooreSearch will return the first index,
// otherwise will return all indexes
func BoyerMooreSearch(str, pattern []byte, once, caseInsensitive bool) []int {
	if len(str) == 0 || len(pattern) == 0 {
		return nil
	}
	if !caseInsensitive {
		pattern = bytes.ToLower(pattern)
	}
	calcBadChar(pattern)
	n, m := len(str), len(pattern)
	suffix := make([]int, m)
	prefix := make([]bool, m)
	calcGoodSuffix(pattern, suffix, prefix)
	i := 0
	var indexes []int
	for i+m <= n {
		var j int
		var sb, pb byte
		for j = m - 1; j >= 0; j-- {
			sb, pb = str[i+j], pattern[j]
			// case insensitive checking
			if (sb == pb) || (!caseInsensitive && toLower(sb) == pb) {
				continue
			} else {
				break
			}
		}
		// match
		if j < 0 {
			indexes = append(indexes, i)
			if once {
				break
			}
			// skip pattern
			i += m
			continue
		}
		x, y := 0, 0
		if !caseInsensitive {
			sb = toLower(sb)
		}
		x = j - bc[sb]

		if j < m-1 {
			y = moveByGoodSuffix(j, m, suffix, prefix)
		}
		i = i + max(x, y)
	}
	return indexes
}

// Substring find all matching substrings
func Substring(str, pattern []byte, caseInsensitive bool) [][]byte {
	indexes := BoyerMooreSearch(str, pattern, false, caseInsensitive)
	var subs [][]byte
	n, m := len(indexes), len(pattern)
	if n > 0 {
		for _, i := range indexes {
			subs = append(subs, str[i:i+m])
		}
	}
	return subs
}
