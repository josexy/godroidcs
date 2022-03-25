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

const (
	Pm      = "pm"
	Fs      = "fs"
	Di      = "device"
	Net     = "net"
	Ctrl    = "ctrl"
	Ms      = "ms"
	Sms     = "sms"
	Contact = "contact"
	CallLog = "calllog"
	Phone   = "phone"
)

const (
	AllPackages       = "all_packages"
	AllUserPackages   = "all_user_packages"
	AllSystemPackages = "all_system_packages"
	Package           = "package"
	Application       = "application"
	Install           = "install"
	ForceInstall      = "force_install"
	Uninstall         = "uninstall"
	ForceUninstall    = "force_uninstall"
	ClearData         = "clear_data"
	AppSize           = "size"
	GetApk            = "apk"
	GetIcon           = "icon"
	ForceStop         = "force_stop"
	Permissions       = "permissions"
	Activities        = "activities"
	Services          = "services"
	Receivers         = "receivers"
	Providers         = "providers"
	SharedLibs        = "sharedlibs"
)

const (
	MkDir         = "mkdir"
	RmDir         = "rmdir"
	Delete        = "delete"
	Create        = "create"
	Download      = "download"
	Upload        = "upload"
	ForceUpload   = "force_upload"
	ForceDownload = "force_download"
	List          = "list"
	Cd            = "cd"
	Pwd           = "pwd"
	Move          = "move"
	Copy          = "copy"
	Rename        = "rename"
	AppendText    = "append"
	WriteText     = "write"
	ReadText      = "read"
)

const (
	MediaFiles = "mediafiles"
	Thumbnail  = "thumbnail"
)

const (
	GetBaseFileTree = "file_tree"
)

const (
	Info          = "info"
	Wifi          = "wifi"
	ScanWifi      = "scan_wifi"
	HasNetwork    = "connectivity"
	ActiveNetwork = "active_network"
	PublicNetwork = "public_network"
)

const (
	Mem      = "mem"
	Storage  = "storage"
	System   = "system"
	Battery  = "battery"
	Display  = "display"
	Location = "location"
	CPU      = "cpu"
	GPU      = "gpu"
)

const (
	ScreenCord = "screencord"
	Reboot     = "reboot"
	Input      = "input"
	Brightness = "brightness"
	Clipboard  = "clipboard"
	Volume     = "volume"
)

const (
	GetScreenBrightness     = "get_brightness"
	SetScreenBrightness     = "set_brightness"
	GetScreenBrightnessMode = "get_brightness_mode"
	SetScreenBrightnessMode = "set_brightness_mode"
	GetClipboard            = "get_clipboard"
	SetClipboard            = "set_clipboard"
	GetVolume               = "get_volume"
	SetVolume               = "set_volume"
	IncreaseVolume          = "increase_volume"
	DecreaseVolume          = "decrease_volume"
)

const (
	StartScreenCapture = "start_screen_capture"
	StopScreenCapture  = "stop_screen_capture"
	StartScreenRecord  = "start_screen_record"
	StopScreenRecord   = "stop_screen_record"
)

const (
	Text     = "text"
	KeyEvent = "keyevent"
	Tap      = "tap"
	Swipe    = "swipe"
)

const (
	BrightnessMode = "mode"
	Manual         = "manual"
	Auto           = "auto"
)

const (
	Dial = "dial"
	Call = "call"
)

const (
	All  = "all"
	Get  = "get"
	Send = "send"
	Body = "body"
	Add  = "add"
)
