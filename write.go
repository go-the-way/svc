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
	"net/http"

	"github.com/gin-gonic/gin"
)

type kv struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func WriteJSON(ctx *gin.Context, code, httpCode int, msg string, err error, data any, encrypts ...bool) {
	dd := kv{Code: code, Msg: msg, Data: data}
	if err != nil {
		dd.Msg = err.Error()
		var cusErr *Error
		if errors.As(err, &cusErr) {
			if cusErr.code > 0 {
				dd.Code = cusErr.code
			}
			if cusErr.httpCode > 0 {
				httpCode = cusErr.httpCode
			}
		}
	}
	encrypt := len(encrypts) > 0 && encrypts[0] && EncryptEnable
	if encrypt {
		marshalBytes, _ := json.Marshal(dd)
		encryptStr, _ := AesEncrypt(marshalBytes)
		ctx.Writer.Header().Set("Encryption", "Yes")
		ctx.String(httpCode, encryptStr)
		return
	}
	ctx.JSON(httpCode, dd)
}

func WriteSuccessJSON(ctx *gin.Context, data any, encrypts ...bool) {
	WriteJSON(ctx, http.StatusOK, http.StatusOK, "", nil, data, encrypts...)
}

func WriteServerErrorJSON(ctx *gin.Context, err error, encrypts ...bool) {
	WriteJSON(ctx, http.StatusInternalServerError, http.StatusInternalServerError, "error", err, nil, encrypts...)
}

func WriteBindError(ctx *gin.Context, err error, encrypts ...bool) {
	WriteJSON(ctx, http.StatusBadRequest, http.StatusBadRequest, "", err, nil, encrypts...)
}

func WriteMessageJSON(ctx *gin.Context, httpCode int, str string, encrypts ...bool) {
	WriteJSON(ctx, httpCode, httpCode, str, nil, nil, encrypts...)
}

func WriteSuccessOrErrorJSON(ctx *gin.Context, err error, encrypts ...bool) {
	if err != nil {
		WriteBindError(ctx, err, encrypts...)
	} else {
		WriteSuccessJSON(ctx, nil, encrypts...)
	}
}
