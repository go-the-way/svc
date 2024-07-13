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
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpRequest[T any](method, url string, header map[string]string, body []byte, options ...func(client *http.Client)) (resp0 *http.Response, resp HttpResponse[T], err error) {
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/json"
	var (
		by  io.Reader
		req *http.Request
	)
	if body != nil {
		var encryptStr string
		if EncryptEnable {
			if encryptStr, err = AesEncrypt(body); err != nil {
				return
			}
			by = bytes.NewBufferString(encryptStr)
			header["Encryption"] = "Yes"
		} else {
			by = bytes.NewBuffer(body)
		}
	}
	if req, err = http.NewRequest(method, url, by); err != nil {
		return
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{Timeout: time.Second * 10}
	if len(options) > 0 {
		for _, opt := range options {
			if opt != nil {
				opt(client)
			}
		}
	}
	if resp0, err = client.Do(req); err != nil {
		return
	}
	var bodyBuf []byte
	if bodyBuf, err = io.ReadAll(resp0.Body); err != nil {
		return
	}
	if len(bodyBuf) <= 0 {
		err = errors.New("response body is empty")
		return
	}
	if encryption := strings.EqualFold(resp0.Header.Get("Encryption"), "Yes"); encryption && DecryptEnable {
		if bodyBuf, err = AesDecrypt(bodyBuf); err != nil {
			return
		}
	}
	err = json.Unmarshal(bodyBuf, &resp)
	return
}

type HttpResponse[T any] struct {
	Code uint   `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
