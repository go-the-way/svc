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

import "io"

var (
	ErrNoReturn = io.ErrNoProgress
)

type Error struct {
	error          string
	httpCode, code int
}

func NewError(error string) *Error { return &Error{error: error} }

func NewErrorWithCode(error string, code int) *Error {
	return &Error{error: error, code: code}
}

func NewErrorWithHttpCode(error string, httpCode int) *Error {
	return &Error{error: error, httpCode: httpCode}
}

func NewErrorWithCodes(error string, httpCode int, code int) *Error {
	return &Error{error: error, httpCode: httpCode, code: code}
}

func (e *Error) Error() string { return e.error }
