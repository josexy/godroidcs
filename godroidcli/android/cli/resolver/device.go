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

	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var DiHelpList = []CommandHelpInfo{
	{internal.Info, "display basic device information"},
	{internal.Mem, "display memory information"},
	{internal.Storage, "display storage space information"},
	{internal.System, "display system information"},
	{internal.Battery, "display battery information"},
	{internal.Display, "display screen information"},
	{internal.Location, "display location information"},
	{internal.GPU, "display GPU information"},
	{internal.CPU, "display CPU frequency information"},
}

type Device struct {
	*ResolverContext
	resolver pb.DeviceResolverClient
}

func NewDevice(conn *grpc.ClientConn) *Device {
	device := &Device{
		resolver: pb.NewDeviceResolverClient(conn),
	}
	return device
}

func (d *Device) SetContext(ctx *ResolverContext) {
	d.ResolverContext = ctx
}

func (d *Device) GetMemoryInfo() (*pb.MemoryInfo, error) {
	return d.resolver.GetMemoryInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetStorageSpaceInfo() (*pb.StorageSpaceInfo, error) {
	return d.resolver.GetStorageSpaceInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetDeviceInfo() (*pb.DeviceInfo, error) {
	return d.resolver.GetDeviceInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetLocationInfo() (*pb.LocationInfo, error) {
	return d.resolver.GetLocationInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetSystemInfo() (*pb.SystemInfo, error) {
	return d.resolver.GetSystemInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetDisplayInfo() (*pb.DisplayInfo, error) {
	return d.resolver.GetDisplayInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetBatteryInfo() (*pb.BatteryInfo, error) {
	return d.resolver.GetBatteryInfo(d.ctx, &pb.Empty{})
}

func (d *Device) GetCPUsFrequency() (*pb.IntegerList, error) {
	return d.resolver.GetCPUsFrequency(d.ctx, &pb.Empty{})
}

func (d *Device) GetGPUInfo() (*pb.GPUInfo, error) {
	return d.resolver.GetGPUInfo(d.ctx, &pb.Empty{})
}

func (d *Device) dumpDeviceInfo() {
	var di *pb.DeviceInfo
	di, d.Error = d.GetDeviceInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("DeviceName", di.DeviceName)
	fn("Manufacturer", di.Manufacturer)
	fn("Product", di.Product)
	fn("Brand", di.Brand)
	fn("Board", di.Board)
	fn("Model", di.Model)
	fn("Hardware", di.Hardware)
	fn("AndroidID", di.AndroidId)
	fn("ROOT", util.BoolToStr(di.Root))
	fn("ADB", util.BoolToStr(di.Adb))
	fn("SIMCard", util.BoolToStr(di.SimCard))
	fn("Developer", util.BoolToStr(di.Developer))
	fn("Airplane", util.BoolToStr(di.Airplane))
	fn("Bluetooth", util.BoolToStr(di.Bluetooth))
	fn("Location", util.BoolToStr(di.Location))
	fn("BuildTime", util.TimeOf(di.BuildTime))
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpSystemInfo() {
	var si *pb.SystemInfo
	si, d.Error = d.GetSystemInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("Host", si.Host)
	fn("User", si.User)
	fn("Display", si.Display)
	fn("Version", si.ReleaseVersion)
	fn("SDK", util.Int32ToStr(si.Sdk))
	fn("Language", si.Language)
	fn("ABI", si.Abi)
	fn("KernelVersion", si.KernelVersion)
	fn("KernelRelease", si.KernelRelease)
	fn("Uptime", util.TimeOfHMS(si.Uptime/1000))
	fn("MCC", util.Int32ToStr(si.Mcc))
	fn("MNC", util.Int32ToStr(si.Mnc))

	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpBatteryInfo() {
	var bi *pb.BatteryInfo
	bi, d.Error = d.GetBatteryInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}

	fn("Level", util.Int32ToStr(bi.Level))
	fn("Max", util.Int32ToStr(bi.Scale))
	fn("Status", bi.Status)
	fn("Health", bi.Health)
	fn("Plugged", bi.Plugged)
	fn("Present", util.BoolToStr(bi.Present))
	fn("Technology", bi.Technology)
	fn("Temperature", fmt.Sprintf("%.1f°C", float32(bi.Temperature)/10.0))
	fn("Voltage", fmt.Sprintf("%dmV", bi.Voltage))
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpDisplayInfo() {
	var di *pb.DisplayInfo
	di, d.Error = d.GetDisplayInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}

	fn("Size", fmt.Sprintf("%dx%d", di.Width, di.Height))
	fn("Density", util.Float32ToStr(di.Density))
	fn("DensityDpi", util.Int32ToStr(di.DensityDpi))
	fn("FontScale", util.Float32ToStr(di.FontScale))
	fn("RefreshRate", util.Float32ToStr(di.RefreshRate)+"Hz")
	fn("Orientation", di.Orientation)
	fn("TouchScreen", util.BoolToStr(di.TouchScreen))
	fn("HDR", util.BoolToStr(di.SupportHdr))
	fn("HDR Capabilities", di.HdrCapabilities)
	fn("ScreenOffTime", util.Int32ToStr(di.ScreenOffTime/1000)+"s")
	fn("ScreenBrightness", util.Int32ToStr(di.ScreenBrightness))
	fn("ScreenBrightnessMode", di.ScreenBrightnessMode)
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpMemoryInfo() {
	var mi *pb.MemoryInfo
	mi, d.Error = d.GetMemoryInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Available"),
		util.Yellow("Used"),
		util.Blue("Total"),
		util.Red("LowMemory"),
		"Threshold",
	})
	table.AddRow(pt.Row{
		util.Green(util.CalcCommonBytes(mi.AvailableMem)),
		util.Yellow(util.CalcCommonBytes(mi.UsedMem)),
		util.Blue(util.CalcCommonBytes(mi.TotalMem)),
		util.Red(util.BoolToStr(mi.LowMemory)),
		util.CalcCommonBytes(mi.Threshold),
	})
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpStorageInfo() {
	var ssi *pb.StorageSpaceInfo
	ssi, d.Error = d.GetStorageSpaceInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}

	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Free"),
		util.Yellow("Used"),
		util.Red("Total"),
	})
	table.AddRow(pt.Row{
		util.Green(util.CalcCommonBytes(ssi.FreeSize)),
		util.Yellow(util.CalcCommonBytes(ssi.UsedSize)),
		util.Red(util.CalcCommonBytes(ssi.TotalSize)),
	})
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpLocationInfo() {
	var li *pb.LocationInfo
	li, d.Error = d.GetLocationInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}

	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("Longitude", util.FloatToStr(li.Longitude))
	fn("Latitude", util.FloatToStr(li.Latitude))
	fn("CountryName", li.CountryName)
	fn("CountryCode", li.CountryCode)
	fn("AdminArea", li.AdminArea)
	fn("Locality", li.Locality)
	fn("SubLocality", li.SubLocality)
	fn("AddressLine", li.AddressLine)
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpCPUFrequency() {
	var list *pb.IntegerList
	list, d.Error = d.GetCPUsFrequency()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{"ID", "Frequency"})
	for i, f := range list.Values {
		table.AddRow(pt.Row{
			"#" + util.IntToStr(i),
			util.Int32ToStr(f/1000) + "MHz",
		})
	}
	table.Filter(d.Param.Node).Print()
}

func (d *Device) dumpGPUInfo() {
	var gi *pb.GPUInfo
	gi, d.Error = d.GetGPUInfo()
	if util.AssertErrorNotNil(d.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("Renderer", gi.Renderer)
	fn("Vendor", gi.Vendor)
	fn("Version", gi.Version)
	table.Filter(d.Param.Node).Print()
}

// Run
// > cmd device info
// > cmd device mem
// > cmd device storage
// > cmd device location
// > cmd device system
// > cmd device display
// > cmd device battery
// > cmd device cpu
// > cmd device gpu
func (d *Device) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.Info:
		d.dumpDeviceInfo()
	case internal.Mem:
		d.dumpMemoryInfo()
	case internal.Storage:
		d.dumpStorageInfo()
	case internal.Location:
		d.dumpLocationInfo()
	case internal.System:
		d.dumpSystemInfo()
	case internal.Battery:
		d.dumpBatteryInfo()
	case internal.Display:
		d.dumpDisplayInfo()
	case internal.CPU:
		d.dumpCPUFrequency()
	case internal.GPU:
		d.dumpGPUInfo()
	default:
		return false
	}
	return true
}
