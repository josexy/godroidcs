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
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	AdbPath     string `json:"adb_path"`
	HistoryFile string `json:"history_file"`
	Address     string `json:"address"`
}

var (
	config *Config
	once   sync.Once
)

func GetConfig() *Config {
	if config == nil {
		// singleton
		once.Do(func() {
			config = new(Config)
			fp, err := os.Open("config.json")
			if err != nil {
				ErrorBy(err)
			}
			defer func() { _ = fp.Close() }()
			decoder := json.NewDecoder(bufio.NewReader(fp))

			if err = decoder.Decode(config); err != nil {
				ErrorBy(err)
			}
		})
	}
	return config
}
