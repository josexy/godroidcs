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
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/josexy/godroidcli/android/cli/stream"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	"github.com/josexy/godroidcli/progressbar"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var PmHelpList = []CommandHelpInfo{
	{internal.AllPackages, "display all installed packages"},
	{internal.Package, "display package information"},
	{internal.Application, "display application information"},
	{internal.Install, "install apk file via PackageInstaller"},
	{internal.ForceInstall, "force install apk file via adb"},
	{internal.Uninstall, "uninstall application via PackageInstaller"},
	{internal.ForceUninstall, "force uninstall application via adb"},
	{internal.ClearData, "clear all application data and cache via adb"},
	{internal.AppSize, "display the size of application"},
	{internal.GetApk, "download apk file"},
	{internal.GetIcon, "download application icon"},
	{internal.ForceStop, "force stop the application via adb"},
	{internal.Permissions, "display all permissions of application"},
	{internal.Activities, "display all activities of application"},
	{internal.Services, "display all services of application"},
	{internal.Receivers, "display all receivers of application"},
	{internal.Providers, "display all content providers of application"},
	{internal.SharedLibs, "display all shared libraries of application"},
}

type PackageManager struct {
	*ResolverContext
	resolver pb.PmResolverClient
}

func NewPackageManager(conn *grpc.ClientConn) *PackageManager {
	pm := &PackageManager{
		resolver: pb.NewPmResolverClient(conn),
	}
	return pm
}

func (p *PackageManager) SetContext(ctx *ResolverContext) {
	p.ResolverContext = ctx
}

func (p *PackageManager) GetAllPackageInfo() (*pb.PackageMetaInfoList, error) {
	return p.resolver.GetAllPackageInfo(p.ctx, &pb.Empty{})
}

func (p *PackageManager) GetAllUserPackageInfo() (*pb.PackageMetaInfoList, error) {
	return p.resolver.GetAllUserPackageInfo(p.ctx, &pb.Empty{})
}

func (p *PackageManager) GetAllSystemPackageInfo() (*pb.PackageMetaInfoList, error) {
	return p.resolver.GetAllSystemPackageInfo(p.ctx, &pb.Empty{})
}

func (p *PackageManager) GetPackageInfo(packageName string) (*pb.PackageInfo, error) {
	return p.resolver.GetPackageInfo(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetApplicationInfo(packageName string) (*pb.ApplicationInfo, error) {
	return p.resolver.GetApplicationInfo(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetAppSize(packageName string) (as *pb.AppSize, err error) {
	return p.resolver.GetApplicationSize(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetPermissions(packageName string) (*pb.StringList, error) {
	return p.resolver.GetPermissions(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetActivities(packageName string) (*pb.StringList, error) {
	return p.resolver.GetActivities(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetServices(packageName string) (*pb.StringList, error) {
	return p.resolver.GetServices(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetReceivers(packageName string) (*pb.StringList, error) {
	return p.resolver.GetReceivers(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetProviders(packageName string) (*pb.StringList, error) {
	return p.resolver.GetProviders(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) GetSharedLibs(packageName string) (*pb.StringList, error) {
	return p.resolver.GetSharedLibFiles(p.ctx, &pb.String{Value: packageName})
}

func (p *PackageManager) InstallApk(reader io.Reader) error {
	s, err := p.resolver.InstallApk(p.ctx)
	return stream.HandleUploadStream(s, err, nil, reader)
}

func (p *PackageManager) UninstallApp(packageName string) (err error) {
	_, err = p.resolver.UninstallApk(p.ctx, &pb.String{Value: packageName})
	return
}

func (p *PackageManager) GetApkFile(packageName string, writer io.Writer, fn stream.ProgressCallback) error {
	var res *pb.String
	var err error
	res, err = p.resolver.GetApk(p.ctx, &pb.String{Value: packageName})
	if err != nil {
		return err
	}

	// download apk file by FileSystem resolver
	fs := p.GetResolver("fs").(*FileSystem)
	err = fs.DownloadFile(res.Value, writer, fn)
	return err
}

func (p *PackageManager) GetAppIcon(packageName string, writer io.Writer, fn stream.ProgressCallback) error {
	s, err := p.resolver.GetIcon(p.ctx, &pb.String{Value: packageName})
	return stream.HandleDownloadStream(s, err, writer, fn)
}

func (p *PackageManager) command(args ...string) {
	p.cmd.SetArgs(args...)
	d, err := p.cmd.CommandReadAll()
	if len(d) == 0 {
		return
	}
	msg := string(d)
	msg = strings.TrimSpace(msg)
	if err != nil {
		util.Error(msg)
	} else {
		util.Info(msg)
	}
}

func (p *PackageManager) dumpAllPackageInfo(typ string) {
	var list *pb.PackageMetaInfoList
	switch typ {
	case "user":
		list, p.Error = p.GetAllUserPackageInfo()
	case "system":
		list, p.Error = p.GetAllSystemPackageInfo()
	default:
		list, p.Error = p.GetAllPackageInfo()
	}
	if util.AssertErrorNotNil(p.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("PackageName"),
		util.Yellow("ApplicationName"),
		util.Blue("Version"),
		util.Red("System"),
	})

	for _, pi := range list.Values {
		appName := pi.AppName
		if len(appName) >= 28 {
			appName = appName[:28] + "..."
		}
		version := pi.VersionName
		if len(version) >= 10 {
			version = version[:10] + "..."
		}
		table.AddRow(pt.Row{
			util.Green(pi.PackageName),
			util.Yellow(appName),
			util.Blue(version),
			util.Red(util.BoolToStr(pi.SystemApp)),
		})
	}
	table.Filter(p.Param.Node).Print()
}

func (p *PackageManager) dumpPackageInfo(packageName string) {
	var pi *pb.PackageInfo
	pi, p.Error = p.GetPackageInfo(packageName)
	if util.AssertErrorNotNil(p.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}

	fn("PackageName", pi.PackageName)
	fn("ApplicationName", pi.ApplicationInfo.AppName)
	fn("ProcessName", pi.ApplicationInfo.ProcessName)
	fn("FirstInstallTime", util.TimeOf(pi.FirstInstallTime))
	fn("LastUpdatedTime", util.TimeOf(pi.LastUpdatedTime))
	fn("Version", pi.VersionName)
	fn("Installer", pi.Installer)
	fn("System", util.BoolToStr(pi.ApplicationInfo.SystemApp))
	fn("DataDir", pi.ApplicationInfo.DataDir)
	fn("SourceDir", pi.ApplicationInfo.SourceDir)
	fn("MinSDKVersion", util.Int32ToStr(pi.ApplicationInfo.MinSdkVersion))
	fn("TargetSDKVersion", util.Int32ToStr(pi.ApplicationInfo.TargetSdkVersion))
	table.Filter(p.Param.Node).Print()
}

func (p *PackageManager) dumpApplicationInfo(packageName string) {
	var ai *pb.ApplicationInfo
	ai, p.Error = p.GetApplicationInfo(packageName)
	if util.AssertErrorNotNil(p.Error) {
		return
	}
	table := pt.NewTable()
	fn := func(name, value string) {
		table.AddRow(pt.Row{util.Green(name), value})
	}
	fn("ApplicationName", ai.AppName)
	fn("ProcessName", ai.ProcessName)
	fn("System", util.BoolToStr(ai.SystemApp))
	fn("DataDir", ai.DataDir)
	fn("SourceDir", ai.SourceDir)
	fn("MinSDKVersion", util.Int32ToStr(ai.MinSdkVersion))
	fn("TargetSDKVersion", util.Int32ToStr(ai.TargetSdkVersion))
	table.Filter(p.Param.Node).Print()
}

func (p *PackageManager) dumpAppSize(packageName string) {
	var as *pb.AppSize
	as, p.Error = p.GetAppSize(packageName)
	if util.AssertErrorNotNil(p.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("AppBytes"),
		util.Yellow("CacheBytes"),
		util.Blue("DataBytes"),
		util.Red("TotalBytes"),
	})
	table.AddRow(pt.Row{
		util.Green(util.CalcCommonBytes(as.AppBytes)),
		util.Yellow(util.CalcCommonBytes(as.CacheBytes)),
		util.Blue(util.CalcCommonBytes(as.DataBytes)),
		util.Red(util.CalcCommonBytes(as.TotalBytes)),
	})
	table.Filter(p.Param.Node).Print()
}

func (p *PackageManager) dumpInstall(apkFile string) {
	var fp io.ReadCloser
	fp, p.Error = os.Open(apkFile)
	if util.AssertErrorNotNil(p.Error) {
		return
	}
	p.Error = p.InstallApk(fp)
	fp.Close()
	util.AssertErrorNotNil(p.Error)
}

func (p *PackageManager) dumpUninstall(packageName string) {
	p.Error = p.UninstallApp(packageName)
	util.AssertErrorNotNil(p.Error)
}

func (p *PackageManager) dumpGetApkFile(packageName, fileName string) {
	if fileName == "" || fileName[len(fileName)-1] == '/' {
		fileName = filepath.Join(fileName, packageName+".apk")
	}
	fp, err := os.Create(fileName)
	if util.AssertErrorNotNil(err) {
		return
	}
	defer func() { _ = fp.Close() }()

	pb := progressbar.New(filepath.Base(fileName))
	p.Error = p.GetApkFile(packageName, fp, func(present, total int64) {
		pb.Update(present, total)
	})
	util.AssertErrorNotNil(p.Error)
}

func (p *PackageManager) dumpGetAppIcon(packageName, fileName string) {
	if fileName == "" || fileName[len(fileName)-1] == '/' {
		fileName = filepath.Join(fileName, packageName+".jpg")
	}
	buf := bytes.NewBuffer(nil)

	pb := progressbar.New(filepath.Base(fileName))
	p.Error = p.GetAppIcon(packageName, buf, func(present, total int64) {
		pb.Update(present, total)
	})
	if util.AssertErrorNotNil(p.Error) {
		return
	}
	_ = ioutil.WriteFile(fileName, buf.Bytes(), 0644)
}

func (p *PackageManager) dumpList(list *pb.StringList, err error) {
	if util.AssertErrorNotNil(err) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{"Name"})
	for _, s := range list.Values {
		table.AddRow(pt.Row{s})
	}
	table.Filter(p.Param.Node).Print()
}

func (p *PackageManager) dumpGetPermissions(packageName string) {
	p.dumpList(p.GetPermissions(packageName))
}

func (p *PackageManager) dumpGetActivities(packageName string) {
	p.dumpList(p.GetActivities(packageName))
}

func (p *PackageManager) dumpGetServices(packageName string) {
	p.dumpList(p.GetServices(packageName))
}

func (p *PackageManager) dumpGetReceivers(packageName string) {
	p.dumpList(p.GetReceivers(packageName))
}

func (p *PackageManager) dumpGetProviders(packageName string) {
	p.dumpList(p.GetProviders(packageName))
}

func (p *PackageManager) dumpGetSharedLibs(packageName string) {
	p.dumpList(p.GetSharedLibs(packageName))
}

// Run
// > cmd pm all_packages [user/system]
// > cmd pm package com.android.chrome
// > cmd pm application com.android.chrome
// > cmd pm size com.android.chrome
// > cmd pm apk com.android.chrome ./base.apk
// > cmd pm icon com.android.chrome ./app.jpg
// > cmd pm install ./base.apk
// > cmd pm uninstall com.android.chrome
// > cmd pm clear_data com.android.chrome
// > cmd pm force_install ./base.apk
// > cmd pm force_uninstall com.android.chrome
// > cmd pm force_stop com.android.chrome
// > cmd pm permissions com.android.chrome
// > cmd pm activities com.android.chrome
// > cmd pm services com.android.chrome
// > cmd pm receivers com.android.chrome
// > cmd pm providers com.android.chrome
// > cmd pm sharedlibs com.android.chrome
func (p *PackageManager) Run(param filter.Param) bool {
	var first, second string
	if len(param.Args) >= 2 {
		first = util.Trim(param.Args[1])
	}
	if len(param.Args) >= 3 {
		second = util.Trim(param.Args[2])
	}
	switch param.Args[0] {
	case internal.AllPackages:
		p.dumpAllPackageInfo(first)
	case internal.Package:
		p.dumpPackageInfo(first)
	case internal.Application:
		p.dumpApplicationInfo(first)
	case internal.AppSize:
		p.dumpAppSize(first)
	case internal.Install:
		p.dumpInstall(first)
	case internal.Uninstall:
		p.dumpUninstall(first)
	case internal.GetApk:
		p.dumpGetApkFile(first, second)
	case internal.GetIcon:
		p.dumpGetAppIcon(first, second)
	case internal.Permissions:
		p.dumpGetPermissions(first)
	case internal.Activities:
		p.dumpGetActivities(first)
	case internal.Services:
		p.dumpGetServices(first)
	case internal.Receivers:
		p.dumpGetReceivers(first)
	case internal.Providers:
		p.dumpGetProviders(first)
	case internal.SharedLibs:
		p.dumpGetSharedLibs(first)
	case internal.ClearData:
		p.command("shell", "pm", "clear", first)
	case internal.ForceInstall:
		p.command("install", first)
	case internal.ForceUninstall:
		p.command("uninstall", first)
	case internal.ForceStop:
		p.command("shell", "am", "force-stop", first)
	default:
		return false
	}
	return true
}
