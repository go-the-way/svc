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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type (
	HttpClientOpt struct {
		Timeout            time.Duration
		InsecureSkipVerify bool
	}
	HttpResponse[T any] struct {
		Code uint   `json:"code"`
		Msg  string `json:"msg"`
		Data T      `json:"data"`

		rawResponse *http.Response `json:"-"`
	}
)

func NewHttpClient(opts ...func(opt *HttpClientOpt)) *http.Client {
	defOpt := &HttpClientOpt{
		Timeout:            time.Second * 3,
		InsecureSkipVerify: true,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			if opt != nil {
				opt(defOpt)
			}
		}
	}

	return &http.Client{
		Timeout: defOpt.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: defOpt.InsecureSkipVerify,
			},
		},
	}
}

func HttpDo[REQ, RESP any](method, url string, header map[string]string, req REQ, opts ...func(client *http.Client)) (resp0 HttpResponse[RESP], resp RESP, err error) {
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/json"
	var (
		by      io.Reader
		req0    *http.Request
		rawResp *http.Response
		buf     []byte
	)
	if reflect.ValueOf(req).IsValid() {
		if buf, err = json.Marshal(req); err != nil {
			return
		}
		var encryptStr string
		if EncryptEnable {
			if encryptStr, err = AesEncrypt(buf); err != nil {
				return
			}
			by = bytes.NewBufferString(encryptStr)
			header["Encryption"] = "Yes"
		} else {
			by = bytes.NewBuffer(buf)
		}
	}
	if req0, err = http.NewRequest(method, url, by); err != nil {
		return
	}
	if header != nil {
		for k, v := range header {
			req0.Header.Set(k, v)
		}
	}

	client := NewHttpClient(func(opt *HttpClientOpt) { opt.Timeout = time.Second * 10 })

	if len(opts) > 0 {
		for _, opt := range opts {
			if opt != nil {
				opt(client)
			}
		}
	}

	if rawResp, err = client.Do(req0); err != nil {
		return
	}
	var bodyBuf []byte
	if bodyBuf, err = io.ReadAll(rawResp.Body); err != nil {
		return
	}
	if len(bodyBuf) <= 0 {
		err = errors.New("response body is empty")
		return
	}
	if encryption := strings.EqualFold(rawResp.Header.Get("Encryption"), "Yes"); encryption && DecryptEnable {
		if bodyBuf, err = AesDecrypt(bodyBuf); err != nil {
			return
		}
	}
	if err = json.Unmarshal(bodyBuf, &resp0); err != nil {
		return
	}
	resp0.rawResponse = rawResp
	resp = resp0.Data
	if resp0.Code != http.StatusOK {
		err = errors.New(resp0.Msg)
	}
	return
}
