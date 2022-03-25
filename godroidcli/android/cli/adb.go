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

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/josexy/godroidcli/android/cli/resolver"
	"github.com/josexy/godroidcli/util"
)

type (
	// Forward port forwarding rule for device
	Forward struct {
		SerialNumber string
		LocalPort    int
		RemotePort   int
	}

	// Device device status
	Device struct {
		SerialNumber string
		Status       string
	}

	// DeviceForward all port forwarding rule for device
	DeviceForward struct {
		Device
		Forwards []Forward
	}

	AdbCmd struct {
		// current context value
		serialNumber string
		// the absolute path of adb command
		Path string
		*resolver.Cmd
		mu sync.Mutex
	}

	DeviceForwardMap map[string]*DeviceForward
)

func NewAdb() *AdbCmd {
	adb := &AdbCmd{
		Cmd: resolver.NewCmd(),
	}
	var err error
	adb.Path, err = exec.LookPath("adb")
	if err != nil {
		adb.Path = util.GetConfig().AdbPath
		if !util.Exist(adb.Path) {
			util.Fatal("adb command not found!")
			os.Exit(0)
		}
	}
	adb.RefreshDeviceList()
	return adb
}

func (adb *AdbCmd) StartServer() {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	adb.SetArgs("start-server")
	adb.Command()
}

func (adb *AdbCmd) StopServer() {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	adb.SetArgs("kill-server")
	adb.Command()
}

// AddForward adb -s SERIAL forward tcp:6666 tcp:9999
func (adb *AdbCmd) AddForward(sn string, local, remote int) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	adb.SetArgs("-s", sn, "forward", fmt.Sprintf("tcp:%d", local), fmt.Sprintf("tcp:%d", remote))
	adb.Command()
}

// RemoveForward adb -s SERIAL forward --remove tcp:6666
func (adb *AdbCmd) RemoveForward(sn string, local int) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	adb.SetArgs("-s", sn, "forward", "--remove", fmt.Sprintf("tcp:%d", local))
	adb.Command()
}

// ListenAtWlan adb -s SERIAL tcpip 7777
func (adb *AdbCmd) ListenAtWlan(sn string, local int) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	adb.SetArgs("-s", sn, "tcpip", util.IntToStr(local))
	adb.Command()
}

// ConnectAtWlan adb -s SERIAL connect 192.168.1.200:7777
func (adb *AdbCmd) ConnectAtWlan(sn, address string) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	old := adb.serialNumber
	adb.serialNumber = sn
	adb.SetArgs("connect", address)
	adb.Command()
	adb.serialNumber = old // restore
}

// DisconnectAtWlan adb -s SERIAL disconnect 192.168.1.200:7777
func (adb *AdbCmd) DisconnectAtWlan(sn, address string) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	old := adb.serialNumber
	adb.serialNumber = sn
	adb.SetArgs("disconnect", address)
	adb.Command()
	adb.serialNumber = old // restore
}

// StopAtWlan adb -s SERIAL usb
// disconnect and stop
func (adb *AdbCmd) StopAtWlan(sn, address string) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	old := adb.serialNumber
	adb.serialNumber = sn
	adb.SetArgs("disconnect", address)
	adb.Command()
	adb.SetArgs("usb")
	adb.Command()
	adb.serialNumber = old // restore
}

func (adb *AdbCmd) SetArgs(args ...string) {
	adb.internalCommand(args...)
}

func (adb *AdbCmd) GetContext() string {
	return adb.serialNumber
}

func (adb *AdbCmd) internalCommand(args ...string) {
	var v []string
	// execute adb command for device
	// adb -s SERIAL
	if adb.serialNumber != "" {
		v = append(v, "-s")
		v = append(v, adb.serialNumber)
	}
	v = append(v, args...)
	adb.SetCmdArgs(adb.Path, v...)
}

func (adb *AdbCmd) CheckPortIsExist(sn string, port int) (DeviceForwardMap, bool) {
	mp := adb.RefreshForwardList()
	if d, ok := mp[sn]; ok {
		for i := 0; i < len(d.Forwards); i++ {
			if port == d.Forwards[i].LocalPort {
				return mp, true
			}
		}
	}
	return mp, false
}

func (adb *AdbCmd) CheckDeviceIsExist(sn string) (DeviceForwardMap, bool) {
	mp := adb.RefreshDeviceList()
	if _, ok := mp[sn]; ok {
		return mp, true
	}
	return mp, false
}

func (adb *AdbCmd) refreshDeviceList0() DeviceForwardMap {
	mp := make(DeviceForwardMap)
	adb.SetArgs("devices")
	lines, err := adb.CommandReadLines()
	if err != nil {
		util.ErrorBy(err)
		return mp
	}
	var start bool
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "List of devices attached") {
			start = true
			continue
		}
		if !start {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		sn, status := fields[0], fields[1]
		mp[sn] = &DeviceForward{
			Device: Device{
				SerialNumber: sn,
				Status:       status,
			},
		}
	}
	return mp
}

func (adb *AdbCmd) RefreshDeviceList() DeviceForwardMap {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	return adb.refreshDeviceList0()
}

func (adb *AdbCmd) RefreshForwardList() DeviceForwardMap {
	adb.mu.Lock()
	defer adb.mu.Unlock()

	mp := adb.refreshDeviceList0()
	// list all port forwarding rules
	adb.SetArgs("forward", "--list")
	lines, err := adb.CommandReadLines()
	if err != nil {
		util.ErrorBy(err)
		return mp
	}
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// SERIAL tcp:xxxx tcp:yyyy
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}

		sn := fields[0]
		local, err := util.StrToInt(fields[1][4:])
		if err != nil {
			continue
		}
		remote, err := util.StrToInt(fields[2][4:])
		if err != nil {
			continue
		}
		mp[sn].Forwards = append(mp[sn].Forwards,
			Forward{
				SerialNumber: sn,
				LocalPort:    local,
				RemotePort:   remote,
			},
		)
	}
	return mp
}

func (adb *AdbCmd) GetAllForwards() (forwards []Forward) {
	mp := adb.RefreshForwardList()
	for _, d := range mp {
		forwards = append(forwards, d.Forwards...)
	}
	return
}

func (adb *AdbCmd) GetAllDevices() (devices []Device) {
	mp := adb.RefreshDeviceList()
	for _, d := range mp {
		devices = append(devices, d.Device)
	}
	return
}

func (adb *AdbCmd) getAllDeviceSerialNumber(string) (list []string) {
	mp := adb.RefreshDeviceList()
	for _, d := range mp {
		list = append(list, d.Device.SerialNumber)
	}
	return
}
