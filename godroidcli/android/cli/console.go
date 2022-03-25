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
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/josexy/godroidcli/android/api"
	"github.com/josexy/godroidcli/android/cli/resolver"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	"github.com/josexy/godroidcli/status"
	"github.com/josexy/godroidcli/util"
)

const (
	AppVersion     = "1.0.0"
	ExecuteTimeout = time.Second * 4
)

const (
	CliLocal     = "!"
	CliHelp      = "help"
	CliExit      = "exit"
	CliVersion   = "version"
	CliLogo      = "logo"
	CliTime      = "time"
	CliList      = "list"
	CliConnect   = "connect"
	CliKill      = "kill"
	CliCheckout  = "checkout"
	CliStart     = "start"
	CliStop      = "stop"
	CliRestart   = "restart"
	CliAdd       = "add"
	CliRm        = "rm"
	CliCmd       = "cmd"
	CliClear     = "clear"
	CliWlan      = "wlan"
	CliDashboard = "dashboard"
)

const (
	Listen     = "listen"
	Connect    = CliConnect
	Disconnect = "disconnect"
	Stop       = "stop"
	Devices    = "devices"
	Forwards   = "forwards"
	Sessions   = "sessions"
)

type (
	CommandInfo struct {
		Usage string
		Func  func(filter.Param)
		root  readline.PrefixCompleterInterface
	}

	ci = CommandInfo

	// command prompt
	LineInfo struct {
		User string
		Host string
		Line string
	}

	Console struct {
		completer *readline.PrefixCompleter
		instance  *readline.Instance
		sessMap   map[string]*Session
		curSess   *Session
		li        *LineInfo
		ctxP      context.Context
		cancel    context.CancelFunc
		wg        sync.WaitGroup
		mu        sync.Mutex
		adb       *AdbCmd
		parser    *filter.CmdParser
		*resolver.Cmd
	}
)

var (
	CommandMap            map[string]ci
	CmdSubCommandHelpInfo map[string][]resolver.CommandHelpInfo
	CmdCommandHelpInfo    []resolver.CommandHelpInfo
	WlanCommandHelpInfo   []resolver.CommandHelpInfo
	ListCommandHelpInfo   []resolver.CommandHelpInfo
)

func NewLineInfo() *LineInfo {
	ctx := new(LineInfo)
	ctx.User = util.GetUser()
	ctx.Host = util.GetHostName()
	return ctx
}

func NewConsole() *Console {
	console := &Console{
		li:      NewLineInfo(),
		adb:     NewAdb(),
		parser:  filter.NewCmdParser(),
		Cmd:     resolver.NewCmd(),
		sessMap: make(map[string]*Session),
	}
	rand.Seed(time.Now().UnixNano())
	console.ctxP, console.cancel = context.WithCancel(context.Background())
	console.init()
	return console
}

func (con *Console) initCommandList() {
	CommandMap = make(map[string]ci)
	CommandMap[CliLocal] = ci{Usage: "execute local command", Func: con.executeLocalSimpleCmd,
		root: readline.PcItem(CliLocal)}
	CommandMap[CliStart] = ci{Usage: "start adb server", Func: con.startServer,
		root: readline.PcItem(CliStart)}
	CommandMap[CliStop] = ci{Usage: "stop adb server", Func: con.stopServer,
		root: readline.PcItem(CliStop)}
	CommandMap[CliRestart] = ci{Usage: "restart adb server", Func: con.restartServer,
		root: readline.PcItem(CliRestart)}
	CommandMap[CliHelp] = ci{Usage: "display help information", Func: con.help,
		root: readline.PcItem(CliHelp)}
	CommandMap[CliExit] = ci{Usage: "quit current program", Func: con.exit,
		root: readline.PcItem(CliExit)}
	CommandMap[CliVersion] = ci{Usage: "display program version", Func: con.version,
		root: readline.PcItem(CliVersion)}
	CommandMap[CliLogo] = ci{Usage: "display logo", Func: con.logo,
		root: readline.PcItem(CliLogo)}
	CommandMap[CliTime] = ci{Usage: "display current system datetime", Func: con.time,
		root: readline.PcItem(CliTime)}
	CommandMap[CliClear] = ci{Usage: "clear console messages", Func: con.clear,
		root: readline.PcItem(CliClear)}
	CommandMap[CliDashboard] = ci{Usage: "start api server", Func: con.dashboard,
		root: readline.PcItem(CliDashboard)}
	CommandMap[CliList] = ci{Usage: "list active devices", Func: con.list,
		root: readline.PcItem(CliList,
			readline.PcItem(Devices),
			readline.PcItem(Forwards),
			readline.PcItem(Sessions))}
	CommandMap[CliWlan] = ci{Usage: "use WLAN(TCP/IP) instead of USB", Func: con.wlan,
		root: readline.PcItem(CliWlan,
			readline.PcItem(Listen,
				readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber)),
			readline.PcItem(CliConnect,
				readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber)),
			readline.PcItem(Disconnect,
				readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber)),
			readline.PcItem(Stop,
				readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber)))}
	CommandMap[CliConnect] = ci{Usage: "connect to device and create a connection session", Func: con.connect,
		root: readline.PcItem(CliConnect,
			readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber))}
	CommandMap[CliKill] = ci{Usage: "close current connection session", Func: con.kill,
		root: readline.PcItem(CliKill,
			readline.PcItemDynamic(con.listSessions))}
	CommandMap[CliCheckout] = ci{Usage: "switch session to another", Func: con.checkout,
		root: readline.PcItem(CliCheckout,
			readline.PcItemDynamic(con.listSessions))}
	CommandMap[CliAdd] = ci{Usage: "add forwarding rule for device", Func: con.add,
		root: readline.PcItem(CliAdd,
			readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber))}
	CommandMap[CliRm] = ci{Usage: "remove forwarding rule for device", Func: con.rm,
		root: readline.PcItem(CliRm,
			readline.PcItemDynamic(con.adb.getAllDeviceSerialNumber))}
	CommandMap[CliCmd] = ci{Usage: "execute session commands", Func: con.resolverCommand,
		root: readline.PcItem(CliCmd,
			readline.PcItem(internal.Pm, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Fs, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Di, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Net, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Ctrl, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Ms, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Sms, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Contact, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.CallLog, readline.PcItemDynamic(GetSubCommand)),
			readline.PcItem(internal.Phone, readline.PcItemDynamic(GetSubCommand)),
		)}
}

func GetSubCommand(s string) (list []string) {
	l := strings.Fields(s)
	if len(l) >= 2 {
		for _, v := range CmdSubCommandHelpInfo[l[1]] {
			list = append(list, v.Name)
		}
	}
	return
}

func (con *Console) initCommandHelpInfo() {
	// cmd
	CmdCommandHelpInfo = make([]resolver.CommandHelpInfo, 10)
	CmdCommandHelpInfo[0] = resolver.CommandHelpInfo{Name: internal.Pm, Usage: "package manager"}
	CmdCommandHelpInfo[1] = resolver.CommandHelpInfo{Name: internal.Di, Usage: "device manager"}
	CmdCommandHelpInfo[2] = resolver.CommandHelpInfo{Name: internal.Fs, Usage: "file system"}
	CmdCommandHelpInfo[3] = resolver.CommandHelpInfo{Name: internal.Net, Usage: "network manager"}
	CmdCommandHelpInfo[4] = resolver.CommandHelpInfo{Name: internal.Ctrl, Usage: "control manager"}
	CmdCommandHelpInfo[5] = resolver.CommandHelpInfo{Name: internal.Ms, Usage: "media store"}
	CmdCommandHelpInfo[6] = resolver.CommandHelpInfo{Name: internal.Sms, Usage: "SMS"}
	CmdCommandHelpInfo[7] = resolver.CommandHelpInfo{Name: internal.Contact, Usage: "contact"}
	CmdCommandHelpInfo[8] = resolver.CommandHelpInfo{Name: internal.CallLog, Usage: "call log"}
	CmdCommandHelpInfo[9] = resolver.CommandHelpInfo{Name: internal.Phone, Usage: "phone call"}

	// wlan
	WlanCommandHelpInfo = make([]resolver.CommandHelpInfo, 4)
	WlanCommandHelpInfo[0] = resolver.CommandHelpInfo{Name: Listen, Usage: "listen a new port on device via TCP/IP"}
	WlanCommandHelpInfo[1] = resolver.CommandHelpInfo{Name: Connect, Usage: "connect to ip:port"}
	WlanCommandHelpInfo[2] = resolver.CommandHelpInfo{Name: Disconnect, Usage: "disconnect to ip:port"}
	WlanCommandHelpInfo[3] = resolver.CommandHelpInfo{Name: Stop, Usage: "disconnect the device and stop listening on the port"}

	// list
	ListCommandHelpInfo = make([]resolver.CommandHelpInfo, 3)
	ListCommandHelpInfo[0] = resolver.CommandHelpInfo{Name: Devices, Usage: "display all found devices"}
	ListCommandHelpInfo[1] = resolver.CommandHelpInfo{Name: Forwards, Usage: "display all forwards for devices"}
	ListCommandHelpInfo[2] = resolver.CommandHelpInfo{Name: Sessions, Usage: "display all connected sessions for devices"}

	// all resolvers help information
	CmdSubCommandHelpInfo = make(map[string][]resolver.CommandHelpInfo)
	CmdSubCommandHelpInfo[internal.Pm] = resolver.PmHelpList
	CmdSubCommandHelpInfo[internal.Di] = resolver.DiHelpList
	CmdSubCommandHelpInfo[internal.Fs] = resolver.FsHelpList
	CmdSubCommandHelpInfo[internal.Net] = resolver.NetHelpList
	CmdSubCommandHelpInfo[internal.Ctrl] = resolver.CtrlHelpList
	CmdSubCommandHelpInfo[internal.Ms] = resolver.MsHelpList
	CmdSubCommandHelpInfo[internal.Sms] = resolver.SmsHelpList
	CmdSubCommandHelpInfo[internal.Contact] = resolver.ContactHelpList
	CmdSubCommandHelpInfo[internal.CallLog] = resolver.CallLogHelpList
	CmdSubCommandHelpInfo[internal.Phone] = resolver.PhoneHelpList
}

func (con *Console) init() {
	con.initCommandList()
	con.initCommandHelpInfo()

	con.completer = readline.NewPrefixCompleter()
	var items []readline.PrefixCompleterInterface
	for _, v := range CommandMap {
		items = append(items, v.root)
	}
	con.completer.SetChildren(items)

	var err error
	con.instance, err = readline.NewEx(&readline.Config{
		Prompt:          con.prompt(),
		AutoComplete:    con.completer,
		InterruptPrompt: "\n", // Ctrl+C
		EOFPrompt:       "\n", // Ctrl+D
		HistoryLimit:    2000,
		HistoryFile:     util.GetConfig().HistoryFile,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			// Ctrl+Z
			if r == readline.CharCtrlZ {
				return r, false
			}
			return r, true
		},
	})
	if err != nil {
		panic(err)
	}
	con.checkAliveSessions()
	con.startup()
}

func (con *Console) prompt() string {
	return fmt.Sprintf("%s@%s > ", color.GreenString(con.li.User), color.RedString(con.li.Host))
}

func (con *Console) updatePrompt() {
	con.instance.SetPrompt(con.prompt())
}

func (con *Console) Run() {
	for {
		line, err := con.instance.Readline()
		// ignore Ctrl+C and Ctrl+D
		if err == readline.ErrInterrupt || err == io.EOF {
			continue
		}
		con.li.Line = strings.TrimSpace(line)

		if err = con.parser.Parse(con.li.Line); err != nil {
			if err != status.ErrEmptyString {
				util.ErrorBy(err)
			}
			continue
		}
		parts := con.parser.Root().Group()
		if len(parts) == 0 {
			continue
		}
		// execute local command
		if strings.HasPrefix(parts[0], "!") {
			// at least 1 argument
			first := parts[0]
			if len(first) == 1 && len(parts) >= 2 {
				// !, ls, -l
				parts = parts[1:]
				con.executeLocalSimpleCmd(filter.Param{Node: con.parser.Root().Right, Args: parts})
			} else if len(parts) >= 1 {
				// !ls, -l
				parts[0] = first[1:]
				if len(parts[0]) > 0 {
					con.executeLocalSimpleCmd(filter.Param{Node: con.parser.Root().Right, Args: parts})
				}
			}
			continue
		}
		if command, ok := CommandMap[parts[0]]; ok {
			command.Func(filter.Param{Node: con.parser.Root().Right, Args: parts})
			con.updatePrompt()
		} else {
			con.notFoundCommand()
		}
	}
}

// listSessions list of all current connection sessions
func (con *Console) listSessions(string) (list []string) {
	con.mu.Lock()
	defer con.mu.Unlock()

	if con.sessMap == nil {
		return
	}
	for name := range con.sessMap {
		list = append(list, name)
	}
	return
}

// checkAliveSessions start a goroutine for checking current alive sessions,
// and delete those dead sessions
// notice: this will ignore those sessions connected via ip address
func (con *Console) checkAliveSessions() {
	go func() {
		for {
			select {
			case <-con.ctxP.Done():
				// notify main goroutine to exit program safely
				con.wg.Done()
				return
			default:
			}

			con.mu.Lock()
			if con.sessMap != nil {
				mp := con.adb.RefreshDeviceList()
				for sn, session := range con.sessMap {
					var flag bool
					for _, i := range mp {
						if i.SerialNumber == sn {
							flag = true
							break
						}
					}
					if !flag {
						// ignore
						if session.address != "127.0.0.1" {
							flag = true
						}
					}
					if !flag {
						// delete dead session
						con.kill(filter.Param{Args: []string{CliKill, sn}})
					}
				}
			}
			con.mu.Unlock()

			// delay
			time.Sleep(time.Second)
		}
	}()
}

func (con *Console) notFoundCommand() {
	util.Error("command not found: [%s]", con.li.Line)
}

func (con *Console) gracefulExit() {
	// close all alive sessions
	for k := range con.sessMap {
		con.kill(filter.Param{Args: []string{CliKill, k}})
	}
	// cancel all child context goroutines
	con.cancel()
}

// > start start adb server
func (con *Console) startServer(filter.Param) {
	con.adb.StartServer()
}

// > stop stop adb server
func (con *Console) stopServer(filter.Param) {
	con.adb.StopServer()
}

// > restart restart adb server
func (con *Console) restartServer(param filter.Param) {
	con.stopServer(param)
	con.startServer(param)
}

// > help display help information for per command
func (con *Console) help(param filter.Param) {
	table := pt.NewTable()
	table.SetHeader(pt.Header{"Name", "Description"})
	if len(param.Args) == 1 {
		for _, v := range CommandMap {
			name := string(v.root.GetName())
			table.AddRow(pt.Row{util.Green(name), v.Usage})
		}
	} else {
		display := func(list []resolver.CommandHelpInfo) {
			for _, info := range list {
				table.AddRow(pt.Row{util.Green(info.Name), info.Usage})
			}
		}
		cmd := param.Args[1]
		switch cmd {
		case CliCmd:
			if len(param.Args) == 2 {
				display(CmdCommandHelpInfo)
			} else {
				cmd = param.Args[2]
				if v, ok := CmdSubCommandHelpInfo[cmd]; ok {
					for _, info := range v {
						table.AddRow(pt.Row{util.Green(info.Name), info.Usage})
					}
				}
			}
		case CliList:
			display(ListCommandHelpInfo)
		case CliWlan:
			display(WlanCommandHelpInfo)
		}
	}
	table.Filter(param.Node).Print()
}

// startup display some messages before enter interactive shell
func (con *Console) startup() {
	con.logo(filter.EmptyParam)
	con.version(filter.EmptyParam)
	con.time(filter.EmptyParam)
	util.Info("adb path: %s", con.adb.Path)

	for _, addr := range util.GetIpv4Addrs() {
		util.Info("local IP address: %s", addr)
	}
}

func (con *Console) assertDeviceExist(sn string) bool {
	if _, ok := con.adb.CheckDeviceIsExist(sn); !ok {
		util.ErrorBy(status.ErrDeviceNotFound)
		return false
	}
	return true
}

func (con *Console) assertForwardExist(sn string, port int) bool {
	if _, ok := con.adb.CheckPortIsExist(sn, port); !ok {
		util.ErrorBy(status.ErrNotFoundOrNotExisted)
		return false
	}
	return true
}

// > exit quit program gracefully
func (con *Console) exit(filter.Param) {
	_ = con.instance.Close()
	con.wg.Add(1)
	con.gracefulExit()

	util.Info("Bye! Have fun! :)")
	// wait for all child goroutines to exit gracefully
	con.wg.Wait()
	os.Exit(0)
}

// > version display current program version information
func (con *Console) version(filter.Param) {
	util.Info("current version: %s", AppVersion)
}

// > logo display logo information
func (con *Console) logo(filter.Param) {
	rnd := rand.Intn(12)
	util.Print("%s\n\n", util.ColorMapFunc[rnd](util.Logo))
}

// > time display current system datetime
func (con *Console) time(filter.Param) {
	util.Info("current datetime: %s", util.TimeOfNow())
}

// > clear clear console messages
func (con *Console) clear(filter.Param) {
	util.ClearScreen(runtime.GOOS)
}

// > dashboard start web server
func (con *Console) dashboard(filter.Param) {
	if con.curSess == nil || con.curSess.status == unavailable {
		util.Warn("current session not found or unavailable and can not start web server")
		return
	}
	w := api.New(con.curSess.CreateSessionProxy())
	if err := w.Start(); err != nil {
		util.ErrorBy(err)
	}
}

// executeLocalSimpleCmd execute simple command from local system environment
// for unix/linux, you can execute some simple POSIX commands which didn't block current program
// for example, "!ls -lh", "!ps aux"
// for windows, you can execute simple commands as well
// for example, "!dir" "!tasklist" "!notepad"
func (con *Console) executeLocalSimpleCmd(param filter.Param) {
	cmdName := param.Args[0]
	var cmdArgs []string
	if len(param.Args) > 1 {
		cmdArgs = param.Args[1:]
	}
	// execute command util execution timeout
	ctx, cancel := context.WithTimeout(con.ctxP, ExecuteTimeout)
	outChan := make(chan []byte, 1)
	defer cancel()
	go func() {
		con.SetCmdArgs(cmdName, cmdArgs...)
		data, err := con.CommandReadAll()
		if err != nil {
			util.ErrorBy(err)
			data = []byte{}
		}
		outChan <- data
	}()
	select {
	case data := <-outChan:
		filter.PipeOutput(data, param.Node)
	case <-ctx.Done():
		util.Warn("execution timeout")
	}
}

// wlan you need to connect to device via USB firstly,
// then you can use TCP/IP to connect device.
// you can safely unmount USB storage when connect successfully
// to look up listen address you `netstat -antl | grep LISTEN` command
// example:
// > wlan listen emulator-5554 7777 # listen a new port on device via serial number
// > wlan connect emulator-5554 192.168.1.200:7777 # connect to device
// > wlan disconnect emulator-5554 192.168.1.200:7777 # disconnect to device
// > wlan stop emulator-5554 192.168.1.200:7777 # disconnect to device and stop listen port
// usage:
// > wlan listen SERIAL PORT
// > wlan connect SERIAL IP:PORT
// > list devices
// > wlan disconnect SERIAL IP:PORT
// > wlan stop SERIAL IP:PORT
func (con *Console) wlan(param filter.Param) {
	defer util.RecoverIllegalOption()

	sn := param.Args[2]
	if !con.assertDeviceExist(sn) {
		return
	}

	value := param.Args[3]
	switch param.Args[1] {
	case Listen:
		port, err := strconv.Atoi(value)
		if err != nil {
			util.ErrorBy(err)
			return
		}
		con.adb.ListenAtWlan(sn, port)
		util.Info("now you can connect to server through `wlan connect %s <IP>:%d` command", sn, port)
	case Connect:
		con.adb.ConnectAtWlan(sn, value)
	case Disconnect:
		con.adb.DisconnectAtWlan(sn, value)
	case Stop:
		con.adb.StopAtWlan(sn, value)
	}
}

// list display all devices, forwards and sessions information
// > list devices
// > list forwards
// > list sessions
func (con *Console) list(param filter.Param) {
	defer util.RecoverIllegalOption()

	var devices []Device
	var forwards []Forward
	table := pt.NewTable()
	switch param.Args[1] {
	case Devices:
		devices = con.adb.GetAllDevices()
		table.SetHeader(pt.Header{"SerialNumber", "Status"})
		for i := range devices {
			sn := devices[i].SerialNumber
			status := devices[i].Status
			// this is a virtual device which comes from TCP/IP connection
			if strings.Contains(sn, ":") {
				sn = util.Yellow(sn)
			}
			// connected successfully
			if status == "device" {
				status = util.Green(status)
			} else {
				status = util.Yellow(status)
			}
			table.AddRow(pt.Row{sn, status})
		}
	case Forwards:
		forwards = con.adb.GetAllForwards()
		table.SetHeader(pt.Header{"SerialNumber", "Local", "Remote"})
		for i := range forwards {
			table.AddRow(pt.Row{forwards[i].SerialNumber,
				util.Green(strconv.Itoa(forwards[i].LocalPort)),
				util.Yellow(strconv.Itoa(forwards[i].RemotePort))})
		}
	case Sessions:
		table.SetHeader(pt.Header{"Name", "Connection", "Status"})
		for name, session := range con.sessMap {
			var status string
			if session.status == alive {
				status = util.Green("alive")
			} else {
				status = util.Red("unavailable")
			}
			table.AddRow(pt.Row{name, session.conn.Target(), status})
		}
	default:
		table = nil
		con.notFoundCommand()
		return
	}
	table.Filter(param.Node).Print()
}

func (con *Console) resetSession(sess *Session) {
	con.curSess = sess
	con.adb.serialNumber = sess.sn
}

func (con *Console) unsetSession(sn string) {
	if con.curSess != nil && con.curSess.sn == sn {
		con.curSess = nil
	}
}

// connect try to connect to server and open a new session
// > connect SERIAL PORT
func (con *Console) connect(param filter.Param) {
	defer util.RecoverIllegalOption()
	snOrIp := param.Args[1]
	if port, err := strconv.Atoi(param.Args[2]); err != nil {
		util.ErrorBy(err)
	} else if sess, ok := con.sessMap[snOrIp]; ok {
		util.Warn("session [%s] has already exist", sess)
	} else if _, ok = con.adb.CheckPortIsExist(snOrIp, port); !ok {
		// try to create to a new session by ip address and port
		util.Warn("not match forward rule %s:%d", snOrIp, port)
		if err = con.newSession("", snOrIp, port); err != nil {
			util.ErrorBy(err)
		}
	} else {
		// create a new session by serial number and port
		if err = con.newSession(snOrIp, "", port); err != nil {
			util.ErrorBy(err)
		}
	}
}

func (con *Console) newSession(sn, address string, port int) error {
	if sn != "" {
		address = "127.0.0.1"
	} else if address != "" {
		sn = address
	}
	util.Info("try to connect to android server...")
	sess, err := NewSession(con.ctxP, sn, address, port, con.adb)
	if err != nil {
		return err
	}

	// create session by serial number
	if sn != address {
		con.sessMap[sn] = sess
	} else {
		// create session by ip address and port
		con.sessMap[fmt.Sprintf("%s:%d", address, port)] = sess
	}
	con.resetSession(sess)
	return nil
}

func (con *Console) randomGenLocalPort() int {
	const min = 12340
	const max = 12370
	return rand.Intn(max-min+1) + min
}

// NewSessionBy create a connection and open session from command
// NewSessionBy("emulator-5554", 9999, false)
// NewSessionBy("192.168.1.161", 9999, true)
func (con *Console) NewSessionBy(name string, port int, device_address bool) error {
	var err error
	// device
	if device_address {
		// mirror rpc server port
		local := con.randomGenLocalPort()
		if err := con.newDeviceForward(name, local, port); err != nil {
			return err
		}
		err = con.newSession(name, "", local)
	} else {
		// TCP/IP
		err = con.newSession("", name, port)
	}
	return err
}

// kill close the session, notice this method may be called by goroutine used for checking
// > kill SERIAL
func (con *Console) kill(param filter.Param) {
	defer util.RecoverIllegalOption()

	sn := param.Args[1]
	if sess, ok := con.sessMap[sn]; ok {
		delete(con.sessMap, sn)
		if err := sess.CloseSession(); err == nil {
			util.Info("kill session [%s]", sess)
			// reset current session is nil
			con.unsetSession(sess.sn)
		} else {
			util.ErrorBy(err)
		}
	}
}

// checkout switch current session to another session by serial number
// > checkout SERIAL
func (con *Console) checkout(param filter.Param) {
	if len(param.Args) == 1 {
		if con.curSess != nil {
			util.Info("current active session is: [%s]", con.curSess)
		} else {
			util.Warn("there is no active session")
		}
		return
	}
	sn := param.Args[1]
	if sess, ok := con.sessMap[sn]; ok {
		con.resetSession(sess)
		util.Info("switch session to [%s]", con.curSess)
	} else {
		util.Error("sessions list not found: %s", sn)
	}
}

// rm remove the port forwarding rule for device
// > rm SERIAL			# remove all port forwarding rules for device
// > rm SERIAL LPORT	# remove one
func (con *Console) rm(param filter.Param) {
	defer util.RecoverIllegalOption()

	sn := param.Args[1]
	if !con.assertDeviceExist(sn) {
		return
	}
	if len(param.Args[2:]) == 0 {
		mp := con.adb.RefreshForwardList()
		// remove all
		if d, ok := mp[sn]; ok {
			for _, forward := range d.Forwards {
				con.adb.RemoveForward(sn, forward.LocalPort)
			}
			// close session
			con.kill(filter.Param{Args: []string{CliKill, sn}})
		}
	} else if local, err := strconv.Atoi(param.Args[2]); err != nil {
		util.ErrorBy(err)
	} else if con.assertForwardExist(sn, local) {
		con.adb.RemoveForward(sn, local)
		con.kill(filter.Param{Args: []string{CliKill, sn}})
	}
}

// add create a new port forwarding rule and save to table.
// you need to given a serial number, local and remote port.
// the forwarding rule is used for listening a new port on local system
// and forward the messages from local system to Android device via ADB command.
// if you connect to device via TCP/IP you don't need to add port forwarding rule.
// PC <---> ADB <---> Android device <---> Application
// > add SERIAL LPORT RPORT
func (con *Console) add(param filter.Param) {
	defer util.RecoverIllegalOption()

	var local, remote int
	var err error
	if local, err = util.StrToInt(param.Args[2]); err != nil {
		util.ErrorBy(err)
	} else if remote, err = util.StrToInt(param.Args[3]); err != nil {
		util.ErrorBy(err)
	} else {
		if err = con.newDeviceForward(param.Args[1], local, remote); err != nil {
			util.ErrorBy(err)
		}
	}
}

func (con *Console) newDeviceForward(sn string, local, remote int) error {
	mp, ok := con.adb.CheckPortIsExist(sn, local)
	df, ok2 := mp[sn]
	if ok {
		return fmt.Errorf("port forwarding rule %s:%d existed", sn, local)
	} else if ok2 && df != nil && df.Status != "device" {
		return fmt.Errorf("device " + sn + " " + df.Status)
	} else if ok2 {
		con.adb.AddForward(sn, local, remote)
		return nil
	} else {
		return status.ErrNotFoundOrNotExisted
	}
}

func (con *Console) resolverCommand(param filter.Param) {
	if con.curSess == nil {
		util.ErrorBy(status.ErrNoSession)
	} else if con.curSess.status == unavailable {
		util.ErrorBy(status.ErrSessionUnavailable)
	} else {
		con.curSess.executeCommand(param)
	}
}
