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
	"io"
	"io/fs"
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

var FsHelpList = []CommandHelpInfo{
	{internal.MkDir, "create an empty directory"},
	{internal.RmDir, "remove all directories and files recursively"},
	{internal.Delete, "delete a file"},
	{internal.Create, "create an empty file"},
	{internal.Download, "pull a file from Android device"},
	{internal.Upload, "push a file to Android device"},
	{internal.ForceDownload, "force pull a file from Android device via adb command"},
	{internal.ForceUpload, "force push a file to Android device via adb command"},
	{internal.List, "list the contents of directory"},
	{internal.Cd, "change current directory to another"},
	{internal.Pwd, "display current working directory"},
	{internal.Move, "move directory or file"},
	{internal.Copy, "copy directory or file"},
	{internal.Rename, "rename directory or file"},
	{internal.AppendText, "append text to existing file"},
	{internal.WriteText, "truncate and write new text to file"},
	{internal.ReadText, "read the entire contents of an existing file"},
}

type InternalDirType int

const (
	GeneralFile = iota
	ApkFile
)

const (
	Image    = "image"
	Video    = "video"
	Audio    = "audio"
	Download = "download"
)

type FileSystem struct {
	*ResolverContext
	resolver  pb.FsResolverClient
	remoteDir string
}

func NewFileSystem(conn *grpc.ClientConn) *FileSystem {
	fs := &FileSystem{
		resolver: pb.NewFsResolverClient(conn),
	}
	return fs
}

func (f *FileSystem) SetContext(ctx *ResolverContext) {
	f.ResolverContext = ctx
}

func (f *FileSystem) concat(path string) string {
	if path == "" {
		return path
	}
	if f.remoteDir != "" && path[0] != '/' {
		path = filepath.Join(f.remoteDir, path)
	}
	return path
}

// DownloadFile download file from Android device and redirect the bytes stream to io.Writer
func (f *FileSystem) DownloadFile(src string, writer io.Writer, fn stream.ProgressCallback) error {
	s, err := f.resolver.DownloadGeneralFile(f.ctx, &pb.String{Value: src})
	return stream.HandleDownloadStream(s, err, writer, fn)
}

// DownloadGeneralFile download file from Android device and write to file
func (f *FileSystem) DownloadGeneralFile(src, dest string) error {
	src = f.concat(src)
	baseSrc := filepath.Base(src)
	if dest == "" {
		dest = filepath.Join(dest, baseSrc)
	} else {
		var info fs.FileInfo
		info, f.Error = os.Stat(dest)
		// delete previous file
		if f.Error == nil && info.IsDir() {
			dest = filepath.Join(dest, baseSrc)
		}
	}
	var fp *os.File
	fp, f.Error = os.Create(dest)
	if f.Error != nil {
		return f.Error
	}

	defer func() { _ = fp.Close() }()

	// show download progressbar
	pb := progressbar.New(dest)
	f.Error = f.DownloadFile(src, fp, func(present, total int64) {
		pb.Update(present, total)
	})
	return f.Error
}

// UploadFile upload byte stream from io.Reader to Android device
func (f *FileSystem) UploadFile(reader io.Reader, dest string) error {
	s, err := f.resolver.UploadGeneralFile(f.ctx)
	return stream.HandleUploadStream(s, err, []byte(dest), reader)
}

// UploadGeneralFile upload file to Android device
func (f *FileSystem) UploadGeneralFile(src, dest string) error {
	// concat dir + file
	if dest == "" || strings.HasSuffix(dest, ".") || strings.HasSuffix(dest, "/") {
		dest = filepath.Join(dest, filepath.Base(src))
	}
	dest = f.concat(dest)
	// read local file
	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	err = f.UploadFile(reader, dest)
	reader.Close()
	return err
}

func (f *FileSystem) GetBaseFileTree(path, id, mode string) (*pb.String, error) {
	return f.resolver.GetBaseFileTree(f.ctx, &pb.StringTuple{First: path, Second: id, Third: mode})
}

func (f *FileSystem) ListDir(remote, mode string) (*pb.FileInfoList, error) {
	return f.resolver.ListDir(f.ctx, &pb.StringPair{First: remote, Second: mode})
}

func (f *FileSystem) CreateFile(file string) (err error) {
	_, err = f.resolver.CreateFile(f.ctx, &pb.String{Value: file})
	return
}

func (f *FileSystem) DeleteFile(file string) (err error) {
	_, err = f.resolver.DeleteFile(f.ctx, &pb.String{Value: file})
	return
}

func (f *FileSystem) MkDir(dir string) (err error) {
	_, err = f.resolver.MkDir(f.ctx, &pb.String{Value: dir})
	return
}

func (f *FileSystem) RmDir(dir string) (err error) {
	_, err = f.resolver.RmDir(f.ctx, &pb.String{Value: dir})
	return
}

func (f *FileSystem) Move(src, dest string) (err error) {
	_, err = f.resolver.Move(f.ctx, &pb.StringPair{First: src, Second: dest})
	return err
}

func (f *FileSystem) Rename(src, dest string) (err error) {
	_, err = f.resolver.Rename(f.ctx, &pb.StringPair{First: src, Second: dest})
	return
}

func (f *FileSystem) Copy(src, dest string) (err error) {
	_, err = f.resolver.Copy(f.ctx, &pb.StringPair{First: src, Second: dest})
	return
}

func (f *FileSystem) WriteText(src, text string) (err error) {
	_, err = f.resolver.WriteText(f.ctx, &pb.StringPair{First: src, Second: text})
	return
}

func (f *FileSystem) AppendText(src, text string) (err error) {
	_, err = f.resolver.AppendText(f.ctx, &pb.StringPair{First: src, Second: text})
	return
}

func (f *FileSystem) ReadText(src string) (*pb.Status, error) {
	return f.resolver.ReadText(f.ctx, &pb.String{Value: src})
}

func (f *FileSystem) doOperand(op string, args ...string) {
	var first, second string
	first = f.concat(args[0])
	if len(args) > 1 {
		second = f.concat(args[1])
	}
	switch op {
	case internal.Create:
		f.CreateFile(first)
	case internal.Delete:
		f.DeleteFile(first)
	case internal.MkDir:
		f.MkDir(first)
	case internal.RmDir:
		f.RmDir(first)
	case internal.Move:
		f.Move(first, second)
	case internal.Rename:
		f.Rename(first, second)
	case internal.Copy:
		f.Copy(first, second)
	default:
	}
	util.AssertErrorNotNil(f.Error)
}

func (f *FileSystem) dumpListFiles(dir string) {
	dir = f.concat(dir)
	var list *pb.FileInfoList

	list, f.Error = f.ListDir(dir, "all")
	if util.AssertErrorNotNil(f.Error) {
		return
	}
	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("Name"),
		util.Yellow("Size"),
		util.Blue("LastModifiedTime"),
		"Owner",
		util.Red("Dir"),
		"Readable",
		"Writable",
		"Executable",
	})
	for _, fi := range list.Values {
		table.AddRow(pt.Row{
			util.Green(filepath.Base(fi.Name)),
			util.Yellow(util.CalcFileBytes(fi.Size)),
			util.Blue(util.TimeOf(fi.LastModifiedTime)),
			fi.Owner,
			util.Red(util.BoolToStr(fi.Dir)),
			util.BoolToStr(fi.Readable),
			util.BoolToStr(fi.Writable),
			util.BoolToStr(fi.Executable),
		})
	}
	table.Filter(f.Param.Node).Print()
}

func (f *FileSystem) dumpUploadOrDownload(s1, s2, op string) {
	switch op {
	case internal.Upload:
		f.Error = f.UploadGeneralFile(s1, s2)
		util.AssertErrorNotNil(f.Error)
	case internal.Download:
		f.Error = f.DownloadGeneralFile(s1, s2)
		util.AssertErrorNotNil(f.Error)
	}
}

func (f *FileSystem) dumpForceUploadOrDownload(s1, s2, op string) {
	var data []byte
	switch op {
	case internal.ForceUpload:
		f.cmd.SetArgs("push", s1, f.concat(s2))
		data, f.Error = f.cmd.CommandReadAll()
	case internal.ForceDownload:
		f.cmd.SetArgs("pull", f.concat(s1), s2)
		data, f.Error = f.cmd.CommandReadAll()
	}
	if f.Error == nil {
		fmt.Print(string(data))
	}
}

func (f *FileSystem) dumpWriteText(op, src, text string) {
	src = f.concat(src)
	switch op {
	case internal.AppendText:
		f.AppendText(src, text)
	case internal.WriteText:
		f.WriteText(src, text)
	}
	if util.AssertErrorNotNil(f.Error) {
		return
	}
}

func (f *FileSystem) dumpReadText(src string) {
	var status *pb.Status
	status, f.Error = f.ReadText(f.concat(src))
	if util.AssertErrorNotNil(f.Error) {
		return
	}
	filter.PipeOutput(util.StringToBytes(status.Message), f.Param.Node)
}

// Run
// > cmd fs upload ./app-debug.apk /data/local/tmp/1.apk
// > cmd fs download /data/local/tmp/tmp.apk ./1.apk
// > cmd fs list /storage/emulated/0/Download
// > cmd fs create /storage/emulated/0/Download/1.txt
// > cmd fs delete /storage/emulated/0/Download/1.txt
// > cmd fs force_download /data/local/tmp/tmp.apk ./1.apk
// > cmd fs force_upload ./app-debug.apk /data/local/tmp/1.apk
// > cmd fs mkdir /storage/emulated/0/Download/hello/world
// > cmd fs rm /storage/emulated/0/Download/hello
// > cmd fs cd "/storage/emulated/0"
// > cmd fs pwd
// > cmd fs write /storage/emulated/0/Download/tmp "this is a text"
// > cmd fs read /storage/emulated/0/Download/tmp
// > cmd fs move /storage/emulated/0/Download/tmp /data/local/tmp
// > cmd fs copy /storage/emulated/0/Download/tmp /data/local/tmp
// > cmd fs rename /storage/emulated/0/Download/tmp /data/local/tmp
func (f *FileSystem) Run(param filter.Param) bool {

	switch param.Args[0] {
	case internal.Upload,
		internal.Download:
		f.dumpUploadOrDownload(util.Trim(param.Args[1]), util.Trim(param.Args[2]), param.Args[0])
	case internal.List:
		f.dumpListFiles(util.Trim(param.Args[1]))
	case internal.Create,
		internal.Delete,
		internal.MkDir,
		internal.RmDir,
		internal.Copy,
		internal.Rename,
		internal.Move:
		f.doOperand(param.Args[0], param.Args[1:]...)
	case internal.ForceUpload,
		internal.ForceDownload:
		f.dumpForceUploadOrDownload(util.Trim(param.Args[1]), util.Trim(param.Args[2]), param.Args[0])
	case internal.Cd:
		if s := util.Trim(param.Args[1]); s != "" {
			f.remoteDir = f.concat(s)
		} else {
			f.remoteDir = ""
		}
		util.Info("change current dir to: %s", f.remoteDir)
	case internal.Pwd:
		if f.remoteDir != "" {
			util.Info("current dir: %s", f.remoteDir)
		}
	case internal.AppendText,
		internal.WriteText:
		f.dumpWriteText(param.Args[0], util.Trim(param.Args[1]), util.Trim(param.Args[2]))
	case internal.ReadText:
		f.dumpReadText(util.Trim(param.Args[1]))
	default:
		return false
	}
	return true
}
