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

package internal

import (
	"io"

	"github.com/josexy/godroidcli/android/cli/stream"
	pb "github.com/josexy/godroidcli/protobuf"
)

type IPackageManager interface {
	GetAllPackageInfo() (*pb.PackageMetaInfoList, error)
	GetAllUserPackageInfo() (*pb.PackageMetaInfoList, error)
	GetAllSystemPackageInfo() (*pb.PackageMetaInfoList, error)
	GetPackageInfo(string) (*pb.PackageInfo, error)
	GetApplicationInfo(string) (*pb.ApplicationInfo, error)
	GetPermissions(string) (*pb.StringList, error)
	GetActivities(string) (*pb.StringList, error)
	GetServices(string) (*pb.StringList, error)
	GetReceivers(string) (*pb.StringList, error)
	GetProviders(string) (*pb.StringList, error)
	GetSharedLibs(string) (*pb.StringList, error)
	GetAppSize(string) (*pb.AppSize, error)
	GetAppIcon(string, io.Writer, stream.ProgressCallback) error
	GetApkFile(string, io.Writer, stream.ProgressCallback) error
	InstallApk(io.Reader) error
	UninstallApp(string) error
}

type INetwork interface {
	GetActivityNetworkDetailInfo() (*pb.DetailActiveNetworkInfoList, error)
	GetPublicNetworkInfo() (*pb.PublicNetworkInfo, error)
	GetNetworkInfo() (*pb.NetInterfaceInfoList, error)
	GetCurrentWifiInfo() (*pb.DetailWifiInfo, error)
	CheckNetworkConnectivity() (*pb.Boolean, error)
	ScanWifi() (*pb.ScanWifiInfoList, error)
}

type IFileSystem interface {
	GetBaseFileTree(string, string, string) (*pb.String, error)
	ListDir(string, string) (*pb.FileInfoList, error)
	CreateFile(string) error
	DeleteFile(string) error
	MkDir(string) error
	RmDir(string) error
	Move(string, string) error
	Rename(string, string) error
	Copy(string, string) error
	WriteText(string, string) error
	AppendText(string, string) error
	ReadText(string) (*pb.Status, error)
	UploadFile(io.Reader, string) error
	DownloadFile(string, io.Writer, stream.ProgressCallback) error
}

type IDevice interface {
	GetMemoryInfo() (*pb.MemoryInfo, error)
	GetStorageSpaceInfo() (*pb.StorageSpaceInfo, error)
	GetDeviceInfo() (*pb.DeviceInfo, error)
	GetLocationInfo() (*pb.LocationInfo, error)
	GetSystemInfo() (*pb.SystemInfo, error)
	GetDisplayInfo() (*pb.DisplayInfo, error)
	GetBatteryInfo() (*pb.BatteryInfo, error)
	GetCPUsFrequency() (*pb.IntegerList, error)
	GetGPUInfo() (*pb.GPUInfo, error)
}

type IController interface {
	GetScreenBrightness() (*pb.Integer, error)
	SetScreenBrightness(int) error
	GetScreenBrightnessMode() (*pb.Integer, error)
	SetScreenBrightnessMode(bool) error
	GetClipboardText() (*pb.String, error)
	SetClipboardText(string) error
	GetVolume() (*pb.Integer, error)
	SetVolume(int) error
	IncreaseVolume() error
	DecreaseVolume() error
	StartScreenCapture() error
	StopScreenCapture() error
	StartScreenRecord() error
	StopScreenRecord() (*pb.String, error)
}

type IMediaStore interface {
	GetMediaFilesInfo(pb.MediaType_Type) (*pb.MediaStoreInfoList, error)
	GetMediaFileThumbnail(string, io.Writer, stream.ProgressCallback) error
	DeleteMediaFile(string) error
	DownloadMediaFile(string, io.Writer, stream.ProgressCallback) error
}

type ISms interface {
	GetAllBasicSmsInfo() (*pb.StringList, error)
	GetSmsInfoList(string) (*pb.SmsInfoList, error)
	SendSms(string, string) error
}

type IContact interface {
	GetAllContactInfo() (*pb.ContactMetaInfoList, error)
	GetContactInfo(string) (*pb.ContactInfo, error)
	DeleteContact(string) error
	AddContact(string, []pb.StringPair, []pb.StringPair) error
}

type ICallLog interface {
	GetAllCallLogInfo() (*pb.CallLogMetaInfoList, error)
	GetCallLogInfo(string) (*pb.CallLogInfoList, error)
	DeleteCallLog(string) error
}

type IPhone interface {
	DialPhone(string) error
	CallPhone(string) error
}
