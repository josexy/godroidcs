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
	"bufio"
	"fmt"
	"os"

	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var ContactHelpList = []CommandHelpInfo{
	{internal.All, "display all contacts"},
	{internal.Get, "display a contact via id"},
	{internal.Delete, "delete a contact via id"},
	{internal.Add, "add a new contact record"},
}

var contactTypeMap = map[string]string{
	"1": "Home",
	"2": "Work",
	"3": "Mobile",
	"4": "Other",
}

type Contact struct {
	*ResolverContext
	resolver      pb.ContactResolverClient
	idUriCacheMap map[string]string
}

func NewContact(conn *grpc.ClientConn) *Contact {
	return &Contact{
		resolver:      pb.NewContactResolverClient(conn),
		idUriCacheMap: make(map[string]string),
	}
}

func (c *Contact) SetContext(ctx *ResolverContext) {
	c.ResolverContext = ctx
}

func (c *Contact) GetContactInfo(id string) (*pb.ContactInfo, error) {
	return c.resolver.GetContactInfo(c.ctx, &pb.String{Value: id})
}

func (c *Contact) GetAllContactInfo() (*pb.ContactMetaInfoList, error) {
	return c.resolver.GetAllContactInfo(c.ctx, &pb.Empty{})
}

func (c *Contact) DeleteContact(uri string) (err error) {
	_, err = c.resolver.DeleteContact(c.ctx, &pb.String{Value: uri})
	return
}

func (c *Contact) AddContact(name string, phones, emails []pb.StringPair) (err error) {
	var phoneList []*pb.ContactInfo_PhoneInfo
	var emailList []*pb.ContactInfo_EmailInfo

	for i := 0; i < len(phones); i++ {
		phoneList = append(phoneList, &pb.ContactInfo_PhoneInfo{Number: phones[i].First, Type: phones[i].Second})
	}
	for i := 0; i < len(emails); i++ {
		emailList = append(emailList, &pb.ContactInfo_EmailInfo{Email: emails[i].First, Type: emails[i].Second})
	}
	_, err = c.resolver.AddContact(c.ctx,
		&pb.ContactInfo{
			Name:   name,
			Phones: phoneList,
			Emails: emailList,
		},
	)
	return
}

func (c *Contact) dumpAllContactInfo() {
	var list *pb.ContactMetaInfoList
	list, c.Error = c.GetAllContactInfo()
	if util.AssertErrorNotNil(c.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("ID"),
		util.Yellow("Name"),
	})
	// reset cache
	for k := range c.idUriCacheMap {
		delete(c.idUriCacheMap, k)
	}
	for _, info := range list.Values {

		var phone, email, id string
		id = util.Int32ToStr(info.Id)
		c.idUriCacheMap[id] = info.Uri

		table.AddRow(pt.Row{
			util.Green(id),
			util.Yellow(info.Name),
			util.Blue(phone),
			util.Red(email),
		})
	}
	table.Filter(c.Param.Node).Print()
}

func (c *Contact) dumpContactInfo(id string) {
	if _, ok := c.idUriCacheMap[id]; !ok {
		util.Warn("The Contact uri cache is empty, please execute 'cmd contact all' to refresh the cache")
		return
	}
	var ci *pb.ContactInfo
	ci, c.Error = c.GetContactInfo(id)
	if util.AssertErrorNotNil(c.Error) {
		return
	}
	table := pt.NewTable()
	table.AddRow(pt.Row{util.Green("ID"), util.Int32ToStr(ci.Id)})
	table.AddRow(pt.Row{util.Green("Name"), ci.Name})

	if ci.Emails != nil {
		for i := 0; i < len(ci.Emails); i++ {
			name := ""
			if i == 0 {
				name = util.Green("Email")
			}
			table.AddRow(pt.Row{name, fmt.Sprintf("%s (%s)", ci.Emails[i].Email, ci.Emails[i].Type)})
		}
	}

	if ci.Phones != nil {
		for i := 0; i < len(ci.Phones); i++ {
			name := ""
			if i == 0 {
				name = util.Green("Phone")
			}
			table.AddRow(pt.Row{name, fmt.Sprintf("%s (%s)", ci.Phones[i].Number, ci.Phones[i].Type)})
		}
	}

	table.Filter(c.Param.Node).Print()
}

func (c *Contact) dumpDeleteContact(id string) {
	if uri, ok := c.idUriCacheMap[id]; ok {
		c.Error = c.DeleteContact(uri)
		util.AssertErrorNotNil(c.Error)
		delete(c.idUriCacheMap, id)
	} else {
		util.Warn("The Contact uri cache is empty, please execute 'cmd contact all' to refresh the cache")
	}
}

func (c *Contact) dumpAddContact() {
	printTypeList := func() {
		util.Print(`1. Home
2. Work
3. Mobile
4. Other
please select a type: `)
	}

	chooseType := func(index string) string {
		if val, ok := contactTypeMap[index]; ok {
			return val
		}
		return contactTypeMap["4"]
	}

	util.Print("please enter a contact name: ")

	var err error
	var name, phone, email, typIndex string
	var numPhones, numEmails int
	var phones, emails []pb.StringPair

	reader := bufio.NewScanner(os.Stdin)
	readline := func() string {
		reader.Scan()
		return reader.Text()
	}
	name = readline()

	util.Print("please enter the number of %s to add: ", util.Green("phone numbers"))
	numPhones, err = util.StrToInt(readline())
	if util.AssertErrorNotNil(err) {
		return
	}

	for i := 0; i < numPhones; i++ {
		util.Print("please enter phone number: ")
		phone = readline()
		printTypeList()
		typIndex = readline()
		phones = append(phones, pb.StringPair{First: phone, Second: chooseType(typIndex)})
	}

	util.Print("please enter the number of %s to add: ", util.Green("emails"))
	numEmails, err = util.StrToInt(readline())
	if util.AssertErrorNotNil(err) {
		return
	}

	for i := 0; i < numEmails; i++ {
		util.Print("please enter email: ")
		email = readline()
		printTypeList()
		typIndex = readline()
		emails = append(emails, pb.StringPair{First: email, Second: chooseType(typIndex)})
	}

	c.Error = c.AddContact(name, phones, emails)
}

// Run
// > cmd contact all
// > cmd contact get [ID]
// > cmd contact delete [ID]
// > cmd contact add
func (c *Contact) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.All:
		c.dumpAllContactInfo()
	case internal.Get:
		c.dumpContactInfo(param.Args[1])
	case internal.Delete:
		c.dumpDeleteContact(param.Args[1])
	case internal.Add:
		c.dumpAddContact()
	default:
		return false
	}
	return true
}
