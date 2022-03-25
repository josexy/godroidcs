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
	"bufio"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/josexy/godroidcli/android/cli/stream"
)

func TestProgressBar(t *testing.T) {
	client := http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse("http://127.0.0.1:7890")
			},
		},
	}

	reqUrl := "https://go.dev/dl/go1.17.5.darwin-arm64.tar.gz"

	fileName := reqUrl[strings.LastIndex(reqUrl, "/")+1:]
	log.Println("download filename:", fileName)

	fp, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)

	// display download progress bar
	pb := New(fileName)
	var fn stream.ProgressCallback = func(present, total int64) {
		pb.Update(present, total)
	}

	resp, err := client.Get(reqUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	reader := bufio.NewReader(resp.Body)

	var present int64
	for {
		n, err := reader.Read(buf[:])
		if err != nil || n == 0 {
			if err == io.EOF {
				err = nil
			}
			break
		}
		writer.Write(buf[:n])
		present += int64(n)
		if fn != nil {
			fn(present, resp.ContentLength)
		}
	}
	writer.Flush()
}
