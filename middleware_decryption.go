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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func Decryption() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if DecryptEnable {
			encryption := strings.EqualFold(ctx.Request.Header.Get("Encryption"), "Yes")
			if encryption {
				if ctx.Request.Method == http.MethodGet {
					// decrypt query string
					qs := ctx.Request.URL.Query()
					if qs.Has("encryption_data") {
						encryptionData := qs.Get("encryption_data")
						if decryptBytes, err := AesDecrypt([]byte(encryptionData)); err != nil {
							fmt.Println(err)
						} else {
							qm, _ := url.ParseQuery(string(decryptBytes))
							ctx.Set("have_encryption_data", "Yes")
							ctx.Set("encryption_data_type", "Query")
							ctx.Set("encryption_data", qm)
						}
					}
				} else {
					// decrypt body
					if readAllBytes, err := io.ReadAll(ctx.Request.Body); err != nil {
						fmt.Println(err)
					} else {
						if decryptBytes, dErr := AesDecrypt(readAllBytes); dErr != nil {
							fmt.Println(dErr)
						} else {
							ctx.Set("have_encryption_data", "Yes")
							ctx.Set("encryption_data_type", "Body")
							ctx.Set("encryption_data", decryptBytes)
						}
					}
				}
			}
		}
		ctx.Next()
	}
}
