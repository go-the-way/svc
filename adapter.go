// Copyright 2025 svc Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svc

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type (
	fsAdapter struct {
		name, prefix string
		dir          bool
		fs           embed.FS
	}
	fileAdapter struct{ file fs.File }
)

// The fs adapter usage:
//
// //go:embed testdata.json
// var fileFS embed.FS
//
// //go:embed scripts/
// var dirFS embed.FS
//
// func init() {
// 	svc.GetApp().StaticFileFS("/testdata.json", "testdata.json", svc.FSFileAdapter("testdata.json", fileFS))
// 	svc.GetApp().StaticFS("/scripts", svc.FSDirAdapter("scripts", dirFS))
// }

func FSFileAdapter(name string, fs embed.FS) http.FileSystem {
	return &fsAdapter{name: name, fs: fs}
}

func FSDirAdapter(prefix string, fs embed.FS) http.FileSystem {
	return &fsAdapter{prefix: strings.Trim(prefix, "/"), dir: true, fs: fs}
}

func newFileAdapter(file fs.File) http.File { return &fileAdapter{file} }

func (a *fsAdapter) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/")
	if a.dir {
		if a.prefix != "" {
			name = path.Join(strings.Trim(a.prefix, "/"), name)
		}
	}
	f, err := a.fs.Open(name)
	if err != nil {
		return nil, err
	}
	file, ok := f.(fs.File)
	if !ok {
		f.Close()
		return nil, fs.ErrInvalid
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	if stat.IsDir() {
		file.Close()
		return nil, errors.New("is a directory")
	}
	return newFileAdapter(file), nil
}

func (f *fileAdapter) Seek(offset int64, whence int) (int64, error) {
	seeker, ok := f.file.(io.Seeker)
	if !ok {
		return 0, fs.ErrInvalid
	}
	return seeker.Seek(offset, whence)
}

func (f *fileAdapter) Readdir(count int) ([]fs.FileInfo, error) {
	dir, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, fs.ErrInvalid
	}
	entries, err := dir.ReadDir(count)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err0 := entry.Info()
		if err0 != nil {
			return nil, err0
		}
		fileInfos = append(fileInfos, info)
	}
	return fileInfos, nil
}

func (f *fileAdapter) Read(p []byte) (int, error) {
	reader, ok := f.file.(io.Reader)
	if !ok {
		return 0, fs.ErrInvalid
	}
	return reader.Read(p)
}

func (f *fileAdapter) Close() error { return f.file.Close() }

func (f *fileAdapter) Stat() (fs.FileInfo, error) { return f.file.Stat() }
