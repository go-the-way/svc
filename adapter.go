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
	FSAdapter struct {
		prefix string
		fs     embed.FS
	}
	FileAdapter struct {
		file fs.File
	}
)

func NewFSAdapter(prefix string, fs embed.FS) http.FileSystem {
	return &FSAdapter{
		prefix: strings.Trim(prefix, "/"),
		fs:     fs,
	}
}

func NewFileAdapter(file fs.File) http.File {
	return &FileAdapter{file: file}
}

func (a *FSAdapter) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(name, "/")
	if a.prefix != "" {
		name = path.Join(a.prefix, name)
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
	return NewFileAdapter(file), nil
}

func (f *FileAdapter) Seek(offset int64, whence int) (int64, error) {
	seeker, ok := f.file.(io.Seeker)
	if !ok {
		return 0, fs.ErrInvalid
	}
	return seeker.Seek(offset, whence)
}

func (f *FileAdapter) Readdir(count int) ([]fs.FileInfo, error) {
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

func (f *FileAdapter) Read(p []byte) (int, error) {
	reader, ok := f.file.(io.Reader)
	if !ok {
		return 0, fs.ErrInvalid
	}
	return reader.Read(p)
}

func (f *FileAdapter) Close() error { return f.file.Close() }

func (f *FileAdapter) Stat() (fs.FileInfo, error) { return f.file.Stat() }
