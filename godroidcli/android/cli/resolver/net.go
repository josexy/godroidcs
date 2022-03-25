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
	"fmt"
	"time"

	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var NetHelpList = []CommandHelpInfo{
	{internal.Info, "display all network interfaces"},
	{internal.Wifi, "display current connected WI-FI information"},
	{internal.ScanWifi, "scan WI-FI"},
	{internal.HasNetwork, "check the current network connection status"},
	{internal.ActiveNetwork, "display currently active network"},
	{internal.PublicNetwork, "get public IP address"},
}

type Network struct {
	*ResolverContext
	resolver pb.NetResolverClient
}

func NewNetwork(conn *grpc.ClientConn) *Network {
	net := &Network{
		resolver: pb.NewNetResolverClient(conn),
	}
	return net
}

func (n *Network) SetContext(ctx *ResolverContext) {
	n.ResolverContext = ctx
}

func (n *Network) GetNetworkInfo() (*pb.NetInterfaceInfoList, error) {
	return n.resolver.GetNetworkInfo(n.ctx, &pb.Empty{})
}

func (n *Network) GetCurrentWifiInfo() (*pb.DetailWifiInfo, error) {
	return n.resolver.GetCurrentWifiInfo(n.ctx, &pb.Empty{})
}

func (n *Network) ScanWifi() (*pb.ScanWifiInfoList, error) {
	var list *pb.ScanWifiInfoList
	var err error
	for {
		list, err = n.resolver.ScanWifiResult(n.ctx, &pb.Empty{})
		if err != nil {
			return nil, err
		}
		if !list.Empty {
			break
		}
		time.Sleep(time.Second)
	}
	return list, err
}

func (n *Network) CheckNetworkConnectivity() (*pb.Boolean, error) {
	return n.resolver.CheckNetworkConnectivity(n.ctx, &pb.Empty{})
}

func (n *Network) GetActivityNetworkDetailInfo() (*pb.DetailActiveNetworkInfoList, error) {
	return n.resolver.GetActiveNetworkInfo(n.ctx, &pb.Empty{})
}

// GetPublicNetworkInfo get current public network from website https://ipinfo.io/json
func (n *Network) GetPublicNetworkInfo() (*pb.PublicNetworkInfo, error) {
	return n.resolver.GetPublicNetworkInfo(n.ctx, &pb.Empty{})
}

func (n *Network) dumpNetworkInfo() {
	var list *pb.NetInterfaceInfoList
	list, n.Error = n.GetNetworkInfo()
	if util.AssertErrorNotNil(n.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Name"),
		util.Yellow("Up"),
		util.Blue("IP"),
		util.Red("BroadcastIP"),
		"MTU",
		"MAC",
		"Prefix",
	})
	for _, info := range list.Values {
		addresses := info.InetAddrs
		for _, addr := range addresses {
			table.AddRow(pt.Row{
				util.Green(info.Name),
				util.Yellow(util.BoolToStr(info.Up)),
				util.Blue(addr.Ip),
				util.Red(addr.BroadcastIp),
				util.Int32ToStr(info.Mtu),
				info.MacAddress,
				util.Int32ToStr(addr.PrefixLength),
			})
		}
	}
	table.Filter(n.Param.Node).Print()
}

func (n *Network) dumpWifiInfo() {
	var wi *pb.DetailWifiInfo
	wi, n.Error = n.GetCurrentWifiInfo()
	if util.AssertErrorNotNil(n.Error) {
		return
	}

	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("SSID", wi.Ssid)
	fn("BSSID", wi.Bssid)
	fn("IP", wi.Ip)
	fn("MAC", wi.Mac)
	fn("Gateway", wi.Gateway)
	fn("Netmask", wi.Netmask)
	fn("Frequency", fmt.Sprintf("%.2fGHz", float32(wi.Frequency)/1000))
	fn("Signal", wi.Signal)
	fn("Dns1", wi.Dns1)
	fn("Dns2", wi.Dns2)
	fn("ServerAddress", wi.ServerAddress)
	fn("NetworkId", util.Int32ToStr(wi.NetworkId))
	fn("LinkSpeed", util.Int32ToStr(wi.LinkSpeed)+"Mbps")
	fn("TxLinkSpeed", util.Int32ToStr(wi.TxSpeed)+"Mbps")
	fn("RxLinkSpeed", util.Int32ToStr(wi.RxSpeed)+"Mbps")
	fn("Status", wi.Status)
	table.Filter(n.Param.Node).Print()
}

func (n *Network) dumpScanWifiList() {
	var list *pb.ScanWifiInfoList
	list, n.Error = n.ScanWifi()
	if util.AssertErrorNotNil(n.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("SSID"),
		util.Yellow("BSSID"),
		util.Blue("Frequency"),
		util.Red("Signal"),
	})
	for _, wi := range list.Values {
		table.AddRow(pt.Row{
			util.Green(wi.Ssid),
			util.Yellow(wi.Bssid),
			util.Blue(fmt.Sprintf("%.2fGHz", float32(wi.Frequency)/1000)),
			util.Red(wi.Signal),
		})
	}
	table.Filter(n.Param.Node).Print()
}

func (n *Network) dumpNetworkConnectivity() {
	var ok *pb.Boolean
	ok, n.Error = n.CheckNetworkConnectivity()
	if util.AssertErrorNotNil(n.Error) {
		return
	}
	if ok.Value {
		util.Info("network connectivity: %s", util.Green("CONNECTED"))
	} else {
		util.Warn("network connectivity: %s", util.Red("LOST"))
	}
}

func (n *Network) dumpAllActiveNetworkInfo() {
	var list *pb.DetailActiveNetworkInfoList
	list, n.Error = n.GetActivityNetworkDetailInfo()
	if util.AssertErrorNotNil(n.Error) {
		return
	}

	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Name"),
		util.Yellow("IP"),
		util.Blue("Type"),
		util.Red("MTU"),
		"Signal",
		"DNS",
		"Network",
		"Proxy",
	})
	for _, info := range list.Values {
		proxy := ""
		if info.Proxy != nil {
			if info.Proxy.Pac != "" {
				proxy = info.Proxy.Pac
			} else if info.Proxy.Host != "" {
				proxy = fmt.Sprintf("%s:%d", info.Proxy.Host, info.Proxy.Port)
			}
		}

		table.AddRow(pt.Row{
			util.Green(info.Name),
			util.Yellow(info.Ip),
			util.Blue(info.Type),
			util.Red(util.Int32ToStr(info.Mtu)),
			info.Signal,
			info.Dns,
			util.BoolToStr(info.HasNetwork),
			proxy,
		})
	}

	table.Filter(n.Param.Node).Print()
}

func (n *Network) dumpPublicNetworkInfo() {
	var pni *pb.PublicNetworkInfo
	pni, n.Error = n.GetPublicNetworkInfo()
	if util.AssertErrorNotNil(n.Error) {
		return
	}

	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("IP", pni.Ip)
	fn("Country", pni.Country)
	fn("Region", pni.Region)
	fn("City", pni.City)
	fn("Location", pni.Location)
	fn("ISP", pni.Isp)
	fn("Timezone", pni.Timezone)
	fn("Hostname", pni.Hostname)

	table.Filter(n.Param.Node).Print()
}

// Run
// > cmd net info
// > cmd net wifi
// > cmd net scan_wifi
// > cmd net connectivity
// > cmd net active_network
// > cmd net public_network
func (n *Network) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.Info:
		n.dumpNetworkInfo()
	case internal.Wifi:
		n.dumpWifiInfo()
	case internal.ScanWifi:
		n.dumpScanWifiList()
	case internal.HasNetwork:
		n.dumpNetworkConnectivity()
	case internal.ActiveNetwork:
		n.dumpAllActiveNetworkInfo()
	case internal.PublicNetwork:
		n.dumpPublicNetworkInfo()
	default:
		return false
	}
	return true
}
