package svc

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Encryption() gin.HandlerFunc {
	encryptEnable := os.Getenv("ENCRYPT_ENABLE") == "T"
	return func(ctx *gin.Context) {
		if encryptEnable {
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
