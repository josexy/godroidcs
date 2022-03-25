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
	"io"
	"os"
	"path/filepath"

	"github.com/josexy/godroidcli/android/cli/stream"
	"github.com/josexy/godroidcli/android/internal"
	"github.com/josexy/godroidcli/filter"
	pt "github.com/josexy/godroidcli/prettytable"
	"github.com/josexy/godroidcli/progressbar"
	pb "github.com/josexy/godroidcli/protobuf"
	"github.com/josexy/godroidcli/util"
	"google.golang.org/grpc"
)

var MsHelpList = []CommandHelpInfo{
	{internal.MediaFiles, "display image/audio/video media files"},
	{internal.Thumbnail, "get media file thumbnail"},
	{internal.Download, "download media file"},
	{internal.Delete, "delete media file by uri"},
}

type MediaStore struct {
	*ResolverContext
	resolver pb.MediaStoreResolverClient
}

func NewMediaStore(conn *grpc.ClientConn) *MediaStore {
	pm := &MediaStore{
		resolver: pb.NewMediaStoreResolverClient(conn),
	}
	return pm
}

func (m *MediaStore) SetContext(ctx *ResolverContext) {
	m.ResolverContext = ctx
}

func (m *MediaStore) GetMediaFilesInfo(typ pb.MediaType_Type) (*pb.MediaStoreInfoList, error) {
	return m.resolver.GetMediaFilesInfo(m.ctx, &pb.MediaType{Type: typ})
}

func (m *MediaStore) GetMediaFileThumbnail(uri string, writer io.Writer, fn stream.ProgressCallback) error {
	s, err := m.resolver.GetMediaFileThumbnail(m.ctx, &pb.String{Value: uri})
	return stream.HandleDownloadStream(s, err, writer, fn)
}

func (m *MediaStore) DeleteMediaFile(uri string) (err error) {
	_, err = m.resolver.DeleteMediaFile(m.ctx, &pb.String{Value: uri})
	return
}

func (m *MediaStore) DownloadMediaFile(uri string, writer io.Writer, fn stream.ProgressCallback) error {
	cs, err := m.resolver.DownloadMediaFile(m.ctx, &pb.String{Value: uri})
	return stream.HandleDownloadStream(cs, err, writer, fn)
}

func (m *MediaStore) selectMediaStoreType(typ string) pb.MediaType_Type {
	var t pb.MediaType_Type
	switch typ {
	case Image:
		t = pb.MediaType_IMAGE
	case Video:
		t = pb.MediaType_VIDEO
	case Audio:
		t = pb.MediaType_AUDIO
	case Download:
		t = pb.MediaType_DOWNLOAD
	default:
		t = -1
	}
	return t
}

func (m *MediaStore) dumpThumbnail(uri, filename string) {
	if uri == "" || filename == "" {
		return
	}
	fp, err := os.Create(filename)
	if util.AssertErrorNotNil(err) {
		return
	}
	defer fp.Close()

	pb := progressbar.New(filepath.Base(filename))
	m.Error = m.GetMediaFileThumbnail(uri, fp, func(present, total int64) {
		pb.Update(present, total)
	})
	if util.AssertErrorNotNil(m.Error) {
		return
	}
}

func (m *MediaStore) dumpDownloadMediaFile(uri, filename string) {
	if uri == "" || filename == "" {
		return
	}
	fp, err := os.Create(filename)
	if util.AssertErrorNotNil(err) {
		return
	}
	defer fp.Close()

	pb := progressbar.New(filepath.Base(filename))
	m.Error = m.DownloadMediaFile(uri, fp, func(present, total int64) {
		pb.Update(present, total)
	})
	if util.AssertErrorNotNil(m.Error) {
		return
	}
}

func (m *MediaStore) dumpMediaFilesInfo(typ string) {
	t := m.selectMediaStoreType(typ)
	var list *pb.MediaStoreInfoList
	if t == -1 {
		return
	}
	list, m.Error = m.GetMediaFilesInfo(t)
	if util.AssertErrorNotNil(m.Error) {
		return
	}

	table := pt.NewTable()
	table.SetHeader(pt.Header{
		util.Green("ID"),
		util.Yellow("Name"),
		util.Blue("Size"),
		util.Red("AddDate"),
		util.Red("ModifyDate"),
		util.Cyan("Uri"),
	})
	for _, info := range list.Values {
		table.AddRow(pt.Row{
			util.Green(util.Int32ToStr(info.Id)),
			util.Yellow(info.Name),
			util.Blue(util.CalcFileBytes(info.Size)),
			util.Red(util.TimeOf(info.DateAdd)),
			util.Red(util.TimeOf(info.DateModify)),
			util.Cyan(info.Uri),
		})
	}
	table.Filter(m.Param.Node).Print()
}

// Run
// > cmd ms mediafiles image/audio/video
// > cmd ms thumbnail content://media/external/images/media/259 ./output.jpg
// > cmd ms delete content://media/external/images/media/259
// > cmd ms download content://media/external/images/media/259 ./output
func (m *MediaStore) Run(param filter.Param) bool {
	switch param.Args[0] {
	case internal.MediaFiles:
		m.dumpMediaFilesInfo(param.Args[1])
	case internal.Thumbnail:
		m.dumpThumbnail(util.Trim(param.Args[1]), util.Trim(param.Args[2]))
	case internal.Download:
		m.dumpDownloadMediaFile(util.Trim(param.Args[1]), util.Trim(param.Args[2]))
	case internal.Delete:
		util.AssertErrorNotNil(m.DeleteMediaFile(util.Trim(param.Args[1])))
	default:
		return false
	}
	return true
}
