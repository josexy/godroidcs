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
	"context"

	"github.com/josexy/godroidcli/filter"
	"github.com/josexy/godroidcli/status"
)

type CommandHelpInfo struct {
	Name  string
	Usage string
}

type Resolver interface {
	SetContext(*ResolverContext)
	Run(filter.Param) bool
}

type AuxResolver interface {
	GetResolver(name string) Resolver
}

type ResolverContext struct {
	Error error
	Param filter.Param
	ctx   context.Context
	cmd   ExecCmd
	Resolver
	AuxResolver
}

func NewResolverContext(ctx context.Context, resolver Resolver, auxResolver AuxResolver, cmd ExecCmd) *ResolverContext {
	return &ResolverContext{
		ctx:         ctx,
		cmd:         cmd,
		Resolver:    resolver,
		AuxResolver: auxResolver,
	}
}

func (ctx *ResolverContext) DoProcess(param filter.Param) {
	ctx.Param = param
	if !ctx.Resolver.Run(param) {
		// throw error
		panic(status.ErrorIllegalOperation)
	}
}
