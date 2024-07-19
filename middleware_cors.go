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
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	mo   = http.MethodOptions
	mh   = http.MethodHead
	aCAO = "Access-Control-Allow-Origin"
	aCAH = "Access-Control-Allow-Headers"
	aCAM = "Access-Control-Allow-Methods"
	aCEH = "Access-Control-Expose-Headers"
	aCAC = "Access-Control-Allow-Credentials"
)

type CorsOption struct {
	AccessControlAllowOrigin      string
	AccessControlAllowHeaders     string
	AccessControlAllowMethods     string
	AccessControlExposeHeaders    string
	AccessControlAllowCredentials string
	PreflightCond                 func(req *http.Request) (ok bool)
	PreflightFunc                 func(ctx *gin.Context)
}

func Cors(configure ...func(opt *CorsOption)) gin.HandlerFunc {
	dco := CorsOption{
		"*",
		"*",
		"POST, GET, OPTIONS, PUT, DELETE",
		"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
		"true",
		func(req *http.Request) (ok bool) { m := req.Method; return m == mo || m == mh },
		func(ctx *gin.Context) { ctx.AbortWithStatus(http.StatusNoContent) },
	}
	return func(ctx *gin.Context) {
		if len(configure) > 0 {
			if conf := configure[0]; conf != nil {
				conf(&dco)
			}
		}
		ctx.Header(aCAO, dco.AccessControlAllowOrigin)
		ctx.Header(aCAH, dco.AccessControlAllowHeaders)
		ctx.Header(aCAM, dco.AccessControlAllowMethods)
		ctx.Header(aCEH, dco.AccessControlExposeHeaders)
		ctx.Header(aCAC, dco.AccessControlAllowCredentials)
		if dco.PreflightCond(ctx.Request) {
			dco.PreflightFunc(ctx)
			return
		}
		ctx.Next()
	}
}
