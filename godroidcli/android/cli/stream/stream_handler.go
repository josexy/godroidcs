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
	"io"

	pb "github.com/josexy/godroidcli/protobuf"
	"google.golang.org/grpc"
)

type (
	ProgressCallback func(present, total int64)
	WriteCallback    func(bytes []byte) error
	ReadCallback     func() ([]byte, []byte, error)
)

func (c ProgressCallback) Call(present, total int64) {
	c(present, total)
}

func (c WriteCallback) Write(bytes []byte) error {
	return c(bytes)
}

func (c ReadCallback) Read() ([]byte, []byte, error) {
	return c()
}

type Handler struct {
	cs          grpc.ClientStream
	preSendData int64
}

func NewStreamHandler(cs grpc.ClientStream) *Handler {
	return &Handler{cs: cs, preSendData: -1}
}

func convertBytes2Long(bs []byte) (value int64) {
	for i := 0; i < 8; i++ {
		value <<= 8
		value |= int64(bs[i])
	}
	return
}

// HandleGet receive the bytes stream from server (download)
func (h *Handler) HandleGet(callback WriteCallback) error {
	m := new(pb.Bytes)
	var err error
	state := 0
	// format: | mask | [extra parameter] | bytes stream |
	// the extra parameter indicates the total size of the received byte stream
completed:
	for {
		err = h.cs.RecvMsg(m)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		switch state {
		case 0: // mask flag
			if m.Value[0] == 0xAA {
				state = 1
			} else {
				state = 2
			}
		case 1: // extra parameter
			h.preSendData = convertBytes2Long(m.Value)
			state = 2
		case 2: // bytes stream
			err = callback.Write(m.Value)
			if err != nil {
				break completed
			}
		}
	}
	return err
}

// HandlePut send the bytes stream to server (upload)
func (h *Handler) HandlePut(callback ReadCallback) error {
	var err error
	var param []byte
	var value []byte
	m := new(pb.ParamBytes)
	m.Value = new(pb.Bytes)

	for {
		// param indicates the storage file path
		param, value, err = callback.Read()

		if param != nil {
			m.Param = &pb.String{Value: string(param)}
		} else {
			// reset param
			m.Param = nil
		}

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		m.Value.Value = value

		err = h.cs.SendMsg(m)
		if err != nil {
			break
		}
	}
	// client can not send message but it can receive the final result
	err = h.cs.CloseSend()
	return err
}

// GetFinalResult the final result is stored in pb.Status
func (h *Handler) GetFinalResult() (*pb.Status, error) {
	m := new(pb.Status)
	if err := h.cs.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
