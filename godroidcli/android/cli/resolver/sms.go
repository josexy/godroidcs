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

package resolver

import (
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var SmsHelpList = []CommandHelpInfo{
	{internal.All, "display all SMS"},
	{internal.Get, "get a SMS by phone number"},
	{internal.Send, "send a new message"},
	{internal.Body, "display SMS body content via id"},
}

type Sms struct {
	*ResolverContext
	resolver       pb.SmsResolverClient
	idBodyCacheMap map[string]string
}

func NewSms(conn *grpc.ClientConn) *Sms {
	return &Sms{
		resolver:       pb.NewSmsResolverClient(conn),
		idBodyCacheMap: make(map[string]string),
	}
}

func (s *Sms) SetContext(ctx *ResolverContext) {
	s.ResolverContext = ctx
}

func (s *Sms) GetSmsInfoList(number string) (*pb.SmsInfoList, error) {
	return s.resolver.GetSmsInfoList(s.ctx, &pb.String{Value: number})
}

func (s *Sms) GetAllBasicSmsInfo() (*pb.StringList, error) {
	return s.resolver.GetAllBasicSmsInfo(s.ctx, &pb.Empty{})
}

func (s *Sms) SendSms(number, message string) (err error) {
	_, err = s.resolver.SendSms(s.ctx, &pb.StringPair{First: number, Second: message})
	return
}

func (s *Sms) dumpAllBasicSmsInfo() {
	var list *pb.StringList
	list, s.Error = s.GetAllBasicSmsInfo()
	if util.AssertErrorNotNil(s.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Number"),
	})
	for _, str := range list.Values {
		table.AddRow(pt.Row{
			util.Green(str),
		})
	}
	table.Filter(s.Param.Node).Print()
}

func (s *Sms) dumpSmsInfoList(number string) {
	var list *pb.SmsInfoList
	list, s.Error = s.GetSmsInfoList(number)
	if util.AssertErrorNotNil(s.Error) {
		return
	}

	// reset message body cache
	for name := range s.idBodyCacheMap {
		delete(s.idBodyCacheMap, name)
	}

	table := pt.NewTable()
	table.SetHeader(pt.Header{
		"ID",
		"Type",
		"Number",
		"Read",
		"Sent Date",
		"Received Date",
	})
	for _, info := range list.Values {
		var typ, id string
		if info.SentReceived {
			typ = util.Green("Sent")
		} else {
			typ = util.Yellow("Received")
		}
		id = util.Int32ToStr(info.Id)
		s.idBodyCacheMap[id] = info.Body
		table.AddRow(pt.Row{
			id,
			typ,
			util.Green(info.Address),
			util.BoolToStr(info.Read),
			util.TimeOf(info.SentDate),
			util.TimeOf(info.ReceivedDate),
		})
	}
	table.Filter(s.Param.Node).Print()
}

func (s *Sms) dumpSendSms(dest, message string) {
	s.Error = s.SendSms(dest, message)
	util.AssertErrorNotNil(s.Error)
}

func (s *Sms) dumpSmsBody(id string) {
	// get message body from cache
	if body, ok := s.idBodyCacheMap[id]; ok {
		filter.PipeOutput(util.StringToBytes(body), s.Param.Node)
	} else {
		util.Warn("The SMS body cache is empty, please execute 'cmd sms get' to refresh the cache")
	}
}

// Run
// > cmd sms all
// > cmd sms get [NUMBER]
// > cmd sms send [NUMBER] [MESSAGE TEXT]
// > cmd sms body [ID]
func (s *Sms) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.All:
		s.dumpAllBasicSmsInfo()
	case internal.Get:
		s.dumpSmsInfoList(util.Trim(param.Args[1]))
	case internal.Send:
		s.dumpSendSms(util.Trim(param.Args[1]), util.Trim(param.Args[2]))
	case internal.Body:
		s.dumpSmsBody(util.Trim(param.Args[1]))
	default:
		return false
	}
	return true
}
