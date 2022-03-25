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
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var PhoneHelpList = []CommandHelpInfo{
	{internal.Dial, "open the dialer interface"},
	{internal.Call, "start dialing"},
}

type Phone struct {
	*ResolverContext
	resolver pb.PhoneResolverClient
}

func NewPhone(conn *grpc.ClientConn) *Phone {
	return &Phone{
		resolver: pb.NewPhoneResolverClient(conn),
	}
}

func (p *Phone) SetContext(ctx *ResolverContext) {
	p.ResolverContext = ctx
}

func (p *Phone) DialPhone(phone string) (err error) {
	_, err = p.resolver.DialPhone(p.ctx, &pb.String{Value: phone})
	return
}

func (p *Phone) CallPhone(phone string) (err error) {
	_, err = p.resolver.CallPhone(p.ctx, &pb.String{Value: phone})
	return
}

// Run
// > cmd phone dial [PHONE NUMBER]
// > cmd phone call [PHONE NUMBER]
func (p *Phone) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.Dial:
		_ = p.DialPhone(util.Trim(param.Args[1]))
	case internal.Call:
		_ = p.CallPhone(util.Trim(param.Args[1]))
	default:
		return false
	}
	return true
}
