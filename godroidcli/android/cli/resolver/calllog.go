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

var CallLogHelpList = []CommandHelpInfo{
	{internal.All, "display all call logs"},
	{internal.Get, "display one call log via phone number"},
	{internal.Delete, "delete one call log via phone number"},
}

type CallLog struct {
	*ResolverContext
	resolver pb.CallLogResolverClient
}

func NewCallLog(conn *grpc.ClientConn) *CallLog {
	return &CallLog{
		resolver: pb.NewCallLogResolverClient(conn),
	}
}

func (c *CallLog) SetContext(ctx *ResolverContext) {
	c.ResolverContext = ctx
}

func (c *CallLog) GetCallLogInfo(number string) (*pb.CallLogInfoList, error) {
	return c.resolver.GetCallLogInfo(c.ctx, &pb.String{Value: number})
}

func (c *CallLog) GetAllCallLogInfo() (*pb.CallLogMetaInfoList, error) {
	return c.resolver.GetAllCallLogInfo(c.ctx, &pb.Empty{})
}

func (c *CallLog) DeleteCallLog(number string) (err error) {
	_, err = c.resolver.DeleteCallLog(c.ctx, &pb.String{Value: number})
	return
}

func (c *CallLog) dumpAllCallLogInfo() {
	var list *pb.CallLogMetaInfoList
	list, c.Error = c.GetAllCallLogInfo()
	if util.AssertErrorNotNil(c.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Number"),
		util.Yellow("Date"),
	})
	for _, info := range list.Values {
		table.AddRow(
			pt.Row{
				util.Green(info.Number),
				util.Yellow(util.TimeOf(info.Date)),
			})
	}
	table.Filter(c.Param.Node).Print()
}

func (c *CallLog) dumpCallLogInfo(number string) {
	var cl *pb.CallLogInfoList
	cl, c.Error = c.GetCallLogInfo(number)
	if util.AssertErrorNotNil(c.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		"ID",
		"Phone Number",
		"Duration",
		"Date",
		"Type",
	})
	for _, info := range cl.Values {
		typ := info.Type
		switch typ {
		case "Missed", "Rejected", "Blocked":
			typ = util.Red(typ)
		case "Incoming", "Outgoing":
			typ = util.Green(typ)
		}
		table.AddRow(pt.Row{
			util.Int32ToStr(info.Id),
			info.Number,
			util.TimeOfHMS(int64(info.Duration)),
			util.TimeOf(info.Date),
			typ,
		})
	}
	table.Filter(c.Param.Node).Print()
}

func (c *CallLog) dumpDeleteCallLog(number string) {
	c.Error = c.DeleteCallLog(number)
	util.AssertErrorNotNil(c.Error)
}

// Run
// > cmd calllog all
// > cmd calllog get [NUMBER]
// > cmd calllog delete [NUMBER]
func (c *CallLog) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.All:
		c.dumpAllCallLogInfo()
	case internal.Get:
		c.dumpCallLogInfo(util.Trim(param.Args[1]))
	case internal.Delete:
		c.dumpDeleteCallLog(util.Trim(param.Args[1]))
	default:
		return false
	}
	return true
}
