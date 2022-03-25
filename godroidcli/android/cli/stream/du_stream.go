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

package stream

import (
	"errors"
	"io"

	pb "github.com/josexy/godroidcli/protobuf"
	"google.golang.org/grpc"
)

// HandleDownloadStream redirect download bytes stream to io.Writer
// you also can convert the download bytes stream to os.File or gin.ResponseWriter
func HandleDownloadStream(cs grpc.ClientStream, err error, writer io.Writer, fn ProgressCallback) error {
	if err != nil {
		return err
	}
	handler := NewStreamHandler(cs)
	var present int64
	err = handler.HandleGet(func(b []byte) error {
		if fn != nil && handler.preSendData != -1 {
			present += int64(len(b))
			fn(present, handler.preSendData)
		}
		_, err := writer.Write(b)
		return err
	})
	return err
}

// HandleUploadStream redirect upload bytes stream to io.Reader
// you also can convert the os.File or gin.ResponseWriter to upload bytes stream
func HandleUploadStream(cs grpc.ClientStream, err error, param []byte, reader io.Reader) error {
	if err != nil {
		return err
	}
	buf := make([]byte, 4096)
	first := true
	// param: indicates whether to save the upload bytes stream to a local file
	if param == nil {
		first = false
	}
	handler := NewStreamHandler(cs)
	err = handler.HandlePut(func() ([]byte, []byte, error) {
		// if the parameter is not equal to nil, it is sent first: {filepath, content}
		// send normally: {content}
		n, err := reader.Read(buf[:])
		if n <= 0 {
			return nil, nil, err
		}
		if first {
			// first send the save file path
			first = false
			return param, buf[:n], err
		}
		return nil, buf[:n], err
	})
	if err != nil {
		return err
	}
	// get the final result from the server
	status, err := handler.GetFinalResult()
	if err != nil {
		return err
	} else if status.Status == pb.Status_FAILED {
		err = errors.New(status.Message)
	}
	return err
}
