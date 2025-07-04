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
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

func Recover(fns ...func()) (err error) {
	for _, fn := range fns {
		if fn != nil {
			func() {
				defer func() {
					if re := recover(); err != nil {
						if err0, ok := re.(error); ok {
							err = err0
						}
					}
				}()
				fn()
			}()
			if err != nil {
				break
			}
		}
	}
	return
}

func MD5(str string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, str)
	buf := hash.Sum(nil)
	return fmt.Sprintf("%x", buf)
}

const (
	fmt1 = "2006-01-02 15:04:05"
	fmt2 = "20060102150405"
)

type number interface {
	/* uint        */ uint8 | uint16 | uint32 | uint | uint64 |
		/* int      */ int8 | int16 | int32 | int | int64 |
		/* float */ float32 | float64
}

func IfFunc(ok bool, fn func()) {
	if ok {
		fn()
	}
}

func IfGt0Func[T number](n T, fn func())   { IfFunc(n > 0, fn) }
func IfNotEmptyFunc(str string, fn func()) { IfFunc(str != "", fn) }
func TimeNow() time.Time                   { return time.Now() }
func TimeNowStr() string                   { return TimeNow().Format(fmt1) }
func TimeNowNum() string                   { return TimeNow().Format(fmt2) }
func TimeNowUnix() string                  { return fmt.Sprintf("%d", TimeNow().Unix()) }
func TimeNowMilli() string                 { return fmt.Sprintf("%d", TimeNow().UnixMilli()) }
func ParseTime(str string) (t time.Time)   { t, _ = time.Parse(fmt1, str); return }
func FormatTime(t time.Time) (str string)  { return t.Format(fmt1) }
