// Copyright 2024 svc Author. All Rights Reserved.
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
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebTryOption struct {
	FsName             string
	ApiRoutePrefix     string
	IndexHtmlName      string
	StatusNotFoundFunc func(err error, ctx *gin.Context)
}

func WebTry(fs embed.FS, configure ...func(opt *WebTryOption)) func(ctx *gin.Context) {
	wto := WebTryOption{
		"www",
		"/api",
		"/index.html",
		func(err error, ctx *gin.Context) { _ = ctx.AbortWithError(http.StatusNotFound, err) },
	}
	return func(ctx *gin.Context) {
		if len(configure) > 0 {
			if conf := configure[0]; conf != nil {
				conf(&wto)
			}
		}
		ctx.Next()
		w := ctx.Writer
		if w.Status() == http.StatusNotFound {
			filePath := ctx.Request.RequestURI
			if !strings.HasPrefix(filePath, wto.ApiRoutePrefix) {
				if filePath == "/" {
					// override to index.html
					filePath = wto.IndexHtmlName
				}
				// index.html index.css index.js index.svg index.ico index.png
				if strings.Contains(filePath, ".") {
					// 静态资源
				} else {
					filePath = wto.IndexHtmlName
				}
				buf, err, contentType := fileInfo(fs, wto.FsName, filePath)
				if err != nil {
					wto.StatusNotFoundFunc(err, ctx)
					return
				}
				w.Header().Add("Content-Type", contentType)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(buf)
				w.Flush()
			}
		}
	}
}

var extMap = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"ico":  "image/x-icon",
	"svg":  "image/svg+xml",
	"png":  "image/png",
	"json": "application/json",
}

func fileInfo(fs embed.FS, fsName, path string) (buf []byte, err error, contentType string) {
	dot := strings.LastIndexByte(path, '.')
	if dot == -1 {
		err = errors.New("404 Not Found")
		return
	}
	fileExt := strings.ToLower(path[dot+1:])
	contentType = extMap[fileExt]
	buf, err = fs.ReadFile(fsName + path)
	return
}
