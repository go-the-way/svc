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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-the-way/validator"
)

func Query(ctx *gin.Context, thenFunc noReqNoRespThenFunc, encrypts ...bool) {
	do[noReq, noResp](ctx, noReq{}, nil, nil, nil, noReqNoRespThenFuncWrap(thenFunc), encrypts...)
}

func QueryReq[REQ any](ctx *gin.Context, req REQ, thenFunc reqNoRespThenFunc[REQ], encrypts ...bool) {
	do[REQ, noResp](ctx, req, bindQuery[REQ], validate[REQ], check[REQ], reqNoRespThenFuncWrap[REQ](thenFunc), encrypts...)
}

func QueryResp[RESP any](ctx *gin.Context, thenFunc noReqRespThenFunc[RESP], encrypts ...bool) {
	do[noReq, RESP](ctx, noReq{}, nil, nil, nil, noReqRespThenFuncWrap[RESP](thenFunc), encrypts...)
}

func QueryReqResp[REQ, RESP any](ctx *gin.Context, req REQ, thenFunc thenFunc[REQ, RESP], encrypts ...bool) {
	do[REQ, RESP](ctx, req, bindQuery[REQ], validate[REQ], check[REQ], thenFunc, encrypts...)
}

func Body(ctx *gin.Context, thenFunc noReqNoRespThenFunc, encrypts ...bool) {
	do[noReq, noResp](ctx, noReq{}, nil, nil, nil, noReqNoRespThenFuncWrap(thenFunc), encrypts...)
}

func BodyReq[REQ any](ctx *gin.Context, req REQ, thenFunc reqNoRespThenFunc[REQ], encrypts ...bool) {
	do[REQ, noResp](ctx, req, bindJSON[REQ], validate[REQ], check[REQ], reqNoRespThenFuncWrap[REQ](thenFunc), encrypts...)
}

func BodyResp[RESP any](ctx *gin.Context, thenFunc noReqRespThenFunc[RESP], encrypts ...bool) {
	do[noReq, RESP](ctx, noReq{}, nil, nil, nil, noReqRespThenFuncWrap[RESP](thenFunc), encrypts...)
}

func BodyReqResp[REQ, RESP any](ctx *gin.Context, req REQ, thenFunc thenFunc[REQ, RESP], encrypts ...bool) {
	do[REQ, RESP](ctx, req, bindJSON[REQ], validate[REQ], check[REQ], thenFunc, encrypts...)
}

func Form(ctx *gin.Context, thenFunc noReqNoRespThenFunc, encrypts ...bool) {
	do[noReq, noResp](ctx, noReq{}, nil, nil, nil, noReqNoRespThenFuncWrap(thenFunc), encrypts...)
}

func FormReq[REQ any](ctx *gin.Context, req REQ, thenFunc reqNoRespThenFunc[REQ], encrypts ...bool) {
	do[REQ, noResp](ctx, req, bindForm[REQ], validate[REQ], check[REQ], reqNoRespThenFuncWrap[REQ](thenFunc), encrypts...)
}

func FormResp[RESP any](ctx *gin.Context, thenFunc noReqRespThenFunc[RESP], encrypts ...bool) {
	do[noReq, RESP](ctx, noReq{}, nil, nil, nil, noReqRespThenFuncWrap[RESP](thenFunc), encrypts...)
}

func FormReqResp[REQ, RESP any](ctx *gin.Context, req REQ, thenFunc thenFunc[REQ, RESP], encrypts ...bool) {
	do[REQ, RESP](ctx, req, bindForm[REQ], validate[REQ], check[REQ], thenFunc, encrypts...)
}

func do[REQ, RESP any](ctx *gin.Context, req REQ, bindFunc bindFunc[REQ], validateFunc validateFunc[REQ], checkFunc checkFunc[REQ], thenFunc thenFunc[REQ, RESP], encrypts ...bool) {
	if fn := bindFunc; fn != nil {
		if err := fn(ctx, &req); err != nil {
			WriteBindError(ctx, err, encrypts...)
			return
		}
	}
	if fn := validateFunc; fn != nil {
		if err := fn(&req); err != nil {
			WriteBindError(ctx, err, encrypts...)
			return
		}
	}
	if fn := checkFunc; fn != nil {
		if err := fn(&req); err != nil {
			WriteBindError(ctx, err, encrypts...)
			return
		}
	}

	if fn := thenFunc; fn != nil {
		if resp, err := thenFunc(req); err != nil {
			if errors.Is(err, ErrNoReturn) {
				// ignored
				// no return everything
			} else {
				WriteServerErrorJSON(ctx, err, encrypts...)
			}
		} else {
			encrypt := len(encrypts) > 0 && encrypts[0] && EncryptEnable
			if respType := reflect.TypeOf(resp); respType != nil {
				respTypeKind := respType.Kind()
				switch {
				case respTypeKind == reflect.String: // for func() (str string, err error)
					respStr := fmt.Sprintf("%v", resp)
					if encrypt {
						encryptStr, _ := AesEncrypt([]byte(respStr))
						ctx.String(http.StatusOK, encryptStr)
					} else {
						ctx.String(http.StatusOK, respStr)
					}
				default:
					WriteSuccessJSON(ctx, resp, encrypt)
				}
			}
		}
	}
}

type (
	noReq  struct{}
	noResp struct{}

	bindFunc[REQ any]       func(ctx *gin.Context, req *REQ) (err error)
	validateFunc[REQ any]   func(req *REQ) (err error)
	checkFunc[REQ any]      func(req *REQ) (err error)
	thenFunc[REQ, RESP any] func(req REQ) (resp RESP, err error)

	reqNoRespThenFunc[REQ any]  func(req REQ) (err error)
	noReqRespThenFunc[RESP any] func() (resp RESP, err error)
	noReqNoRespThenFunc         func() (err error)
)

func reqNoRespThenFuncWrap[REQ any](thenFunc reqNoRespThenFunc[REQ]) thenFunc[REQ, noResp] {
	return func(req REQ) (resp noResp, err error) { err = thenFunc(req); return }
}

func noReqRespThenFuncWrap[RESP any](thenFunc noReqRespThenFunc[RESP]) thenFunc[noReq, RESP] {
	return func(req noReq) (resp RESP, err error) { resp, err = thenFunc(); return }
}

func noReqNoRespThenFuncWrap(thenFunc noReqNoRespThenFunc) thenFunc[noReq, noResp] {
	return func(req noReq) (resp noResp, err error) { err = thenFunc(); return }
}

func haveEncryptionData(ctx *gin.Context, dataType string) bool {
	b1 := false
	b2 := false
	{
		if value, ok := ctx.Get("have_encryption_data"); ok {
			if data, okk := value.(string); okk {
				if data == "Yes" {
					b1 = true
				}
			}
		}
	}
	{
		if value, ok := ctx.Get("encryption_data_type"); ok {
			if data, okk := value.(string); okk {
				if data == dataType {
					b2 = true
				}
			}
		}
	}
	return b1 && b2
}

func bindQuery[REQ any](ctx *gin.Context, req *REQ) (err error) {
	if haveEncryptionData(ctx, "Query") {
		if value, ok := ctx.Get("encryption_data"); ok {
			if values, okk := value.(url.Values); okk {
				return mapForm(req, values)
			}
		}
	}
	return ctx.ShouldBindQuery(req)
}

func bindJSON[REQ any](ctx *gin.Context, req *REQ) (err error) {
	if haveEncryptionData(ctx, "Body") {
		if value, ok := ctx.Get("encryption_data"); ok {
			if values, okk := value.([]byte); okk {
				return json.Unmarshal(values, req)
			}
		}
	}
	return ctx.ShouldBindJSON(req)
}

func bindForm[REQ any](ctx *gin.Context, req *REQ) (err error) {
	if haveEncryptionData(ctx, "Body") {
		if value, ok := ctx.Get("encryption_data"); ok {
			if values, okk := value.([]byte); okk {
				return json.Unmarshal(values, req)
			}
		}
	}
	return ctx.ShouldBindWith(req, binding.Form)
}

func validate[REQ any](req *REQ) (err error) {
	v := validator.New(req)
	if vr := v.Validate(); !vr.Passed {
		err = errors.New(vr.Messages())
	}
	return
}

func check[REQ any](req *REQ) (err error) {
loop:
	for _, value := range []reflect.Value{
		reflect.ValueOf(req).MethodByName("Check"),
		reflect.ValueOf(&req).MethodByName("Check"),
	} {
		if value.IsValid() {
			if values := value.Call([]reflect.Value{}); values != nil && len(values) > 0 {
				for _, val := range values {
					if val.CanInterface() {
						if inter := val.Interface(); inter != nil {
							if err = inter.(error); err != nil {
								break loop
							}
						}
					}
				}
			}
		}
	}
	return
}
