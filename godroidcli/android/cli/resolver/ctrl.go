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
	"strconv"
	"time"

	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var CtrlHelpList = []CommandHelpInfo{
	{internal.ScreenCord, "record the screen and save it in mp4 format"},
	{internal.Reboot, "reboot the device via adb"},
	{internal.Input, "input command via adb"},
	{internal.Brightness, "display or set the brightness of device"},
	{internal.Clipboard, "display or set the clipboard text of device"},
	{internal.Volume, "display or set the volume of device"},
}

type Controller struct {
	*ResolverContext
	resolver pb.ControlResolverClient
}

func NewController(conn *grpc.ClientConn) *Controller {
	ctrl := &Controller{
		resolver: pb.NewControlResolverClient(conn),
	}
	return ctrl
}

func (c *Controller) SetContext(ctx *ResolverContext) {
	c.ResolverContext = ctx
}

// StartScreenCapture android opens the Websocket service and pushes the stream to the client,
// so as to realize the android screen projection
func (c *Controller) StartScreenCapture() (err error) {
	_, err = c.resolver.StartScreenCapture(c.ctx, &pb.Empty{})
	return
}

// StopScreenCapture android stop the Websocket service
func (c *Controller) StopScreenCapture() (err error) {
	_, err = c.resolver.StopScreenCapture(c.ctx, &pb.Empty{})
	return
}

// StartScreenRecord android record screen and save it to a temporary directory
func (c *Controller) StartScreenRecord() (err error) {
	_, err = c.resolver.StartScreenRecord(c.ctx, &pb.Empty{})
	return
}

// StopScreenRecord android stop recording screen
func (c *Controller) StopScreenRecord() (*pb.String, error) {
	return c.resolver.StopScreenRecord(c.ctx, &pb.Empty{})
}

// ScreenRecord record the screen and save it in mp4 format
// note: the video file is saved on Android and downloaded when recording is stopped
func (c *Controller) ScreenRecord(local string) {
	util.Info("start media projection...")
	util.Info("now you may need to go back to the device to grant permissions")
	c.Error = c.StartScreenRecord()
	if util.AssertErrorNotNil(c.Error) {
		return
	}

	util.Info("save to video file when stop recording (Ctrl+C)")

	<-util.MakeInterruptChan()
	util.Warn("stop media projection...")

	var s *pb.String
	s, c.Error = c.StopScreenRecord()
	if util.AssertErrorNotNil(c.Error) {
		return
	}
	if s.Value != "" {
		util.Info("video file path: %s", s.Value)
		// delay for a while and wait for data preparation to complete
		time.Sleep(time.Second * 2)
		// fetch bytes stream
		fs := c.GetResolver("fs").(*FileSystem)
		fs.DownloadGeneralFile(s.Value, local)
		util.Info("delete the video file on Android device")
		fs.DeleteFile(s.Value)
	}
}

// DoText > ctrl input text "hello world"
func (c *Controller) DoText(text string) {
	c.cmd.SetArgs("shell", internal.Input, internal.Text, text)
	c.cmd.Command()
}

// DoKeyEvent > ctrl input keyevent CODE
func (c *Controller) DoKeyEvent(code string) {
	c.cmd.SetArgs("shell", internal.Input, internal.KeyEvent, code)
	c.cmd.Command()
}

// DoSwipe > ctrl input swipe X Y NX NY DELAY
func (c *Controller) DoSwipe(x, y, nx, ny, delay string) {
	if delay != "" {
		c.cmd.SetArgs("shell", internal.Input, internal.Swipe, x, y, nx, ny, delay)
	} else {
		c.cmd.SetArgs("shell", internal.Input, internal.Swipe, x, y, nx, ny)
	}
	c.cmd.Command()
}

// DoTap > ctrl input tap X Y
func (c *Controller) DoTap(x, y string) {
	c.cmd.SetArgs("shell", internal.Input, internal.Tap, x, y)
	c.cmd.Command()
}

func (c *Controller) GetScreenBrightness() (*pb.Integer, error) {
	return c.resolver.GetScreenBrightness(c.ctx, &pb.Empty{})
}

func (c *Controller) SetScreenBrightness(value int) (err error) {
	_, err = c.resolver.SetScreenBrightness(c.ctx, &pb.Integer{Value: int32(value)})
	return
}

func (c *Controller) GetScreenBrightnessMode() (*pb.Integer, error) {
	return c.resolver.GetScreenBrightnessMode(c.ctx, &pb.Empty{})
}

func (c *Controller) SetScreenBrightnessMode(value bool) (err error) {
	_, err = c.resolver.SetScreenBrightnessMode(c.ctx, &pb.Boolean{Value: value})
	return
}

func (c *Controller) GetClipboardText() (*pb.String, error) {
	return c.resolver.GetClipboardText(c.ctx, &pb.Empty{})
}

func (c *Controller) SetClipboardText(text string) (err error) {
	_, err = c.resolver.SetClipboardText(c.ctx, &pb.String{Value: text})
	return
}

func (c *Controller) GetVolume() (*pb.Integer, error) {
	return c.resolver.GetVolume(c.ctx, &pb.Empty{})
}

func (c *Controller) SetVolume(value int) (err error) {
	_, err = c.resolver.SetVolume(c.ctx, &pb.Integer{Value: int32(value)})
	return
}

func (c *Controller) IncreaseVolume() (err error) {
	_, err = c.resolver.IncreaseVolume(c.ctx, &pb.Empty{})
	return
}

func (c *Controller) DecreaseVolume() (err error) {
	_, err = c.resolver.DecreaseVolume(c.ctx, &pb.Empty{})
	return
}

func (c *Controller) Reboot() {
	c.cmd.SetArgs("reboot")
	c.cmd.Command()
}

func (c *Controller) ProcessInputCommand(args ...string) {
	switch args[0] {
	case internal.Text:
		c.DoText(args[1])
	case internal.KeyEvent:
		c.DoKeyEvent(args[1])
	case internal.Swipe:
		if len(args) >= 6 {
			c.DoSwipe(args[1], args[2], args[3], args[4], args[5])
		} else {
			c.DoSwipe(args[1], args[2], args[3], args[4], "")
		}
	case internal.Tap:
		c.DoTap(args[1], args[2])
	}
}

func (c *Controller) dumpBrightness(args ...string) {
	// getter
	if len(args) == 0 {
		var i *pb.Integer
		if i, c.Error = c.GetScreenBrightness(); c.Error == nil {
			util.Info("current brightness value: %d", i.Value)
		} else {
			util.ErrorBy(c.Error)
		}
	} else {
		// getter
		if args[0] == internal.BrightnessMode {
			c.dumpBrightnessMode(args[1:]...)
			return
		}
		// setter
		if v, err := strconv.Atoi(args[0]); err == nil {
			if c.SetScreenBrightness(v); c.Error == nil {
				util.Info("set brightness successfully")
			} else {
				util.ErrorBy(c.Error)
			}
		} else {
			util.ErrorBy(err)
		}
	}
}

func (c *Controller) dumpBrightnessMode(args ...string) {
	if len(args) == 0 {
		var i *pb.Integer
		if i, c.Error = c.GetScreenBrightnessMode(); c.Error == nil {
			var mode = ""
			if c.Error == nil {
				switch i.Value {
				case 1:
					mode = internal.Auto
				case 0:
					mode = internal.Manual
				}
			}
			util.Info("current brightness mode: %s", mode)
		} else {
			util.ErrorBy(c.Error)
		}
	} else {
		if c.Error = c.SetScreenBrightnessMode(args[0] == internal.Auto); c.Error == nil {
			util.Info("set brightness mode successfully")
		} else {
			util.ErrorBy(c.Error)
		}
	}
}

func (c *Controller) dumpClipboard(args ...string) {
	if len(args) == 0 {
		var value *pb.String
		value, c.Error = c.GetClipboardText()
		if util.AssertErrorNotNil(c.Error) {
			return
		}
		filter.PipeOutput(util.StringToBytes(value.Value), c.Param.Node)
	} else {
		c.SetClipboardText(util.Trim(args[0]))
	}
}

func (c *Controller) dumpVolume(args ...string) {
	if len(args) == 0 {
		var i *pb.Integer
		i, c.Error = c.GetVolume()
		if util.AssertErrorNotNil(c.Error) {
			return
		}
		util.Info(fmt.Sprintf("current volume value: %d", i.Value))
	} else {
		value := args[0]
		switch value {
		case "+":
			c.IncreaseVolume()
		case "-":
			c.DecreaseVolume()
		default:
			v, err := strconv.Atoi(value)
			if util.AssertErrorNotNil(err) {
				return
			}
			c.SetVolume(v)
		}
	}
}

// Run
// > ctrl screencord ./example.mp4
// > ctrl input
// > ctrl reboot
// > ctrl brightness [100]
// > ctrl brightness mode [manual/auto]
// > ctrl clipboard ["hello world"]
// > ctrl volume [10]
// > ctrl volume [+/-]
func (c *Controller) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.ScreenCord:
		c.ScreenRecord(util.Trim(param.Args[1]))
	case internal.Input:
		c.ProcessInputCommand(param.Args[1:]...)
	case internal.Reboot:
		c.Reboot()
	case internal.Brightness:
		c.dumpBrightness(param.Args[1:]...)
	case internal.Clipboard:
		c.dumpClipboard(param.Args[1:]...)
	case internal.Volume:
		c.dumpVolume(param.Args[1:]...)
	default:
		return false
	}
	return true
}
