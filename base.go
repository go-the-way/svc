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
	"gorm.io/gorm"
)

type (
	PageReq struct {
		Limit int `form:"limit"`
		Page  int `form:"page"`
	}
	PageResp[T any] struct {
		Total int64 `json:"total"`
		List  []T   `json:"list"`
	}
)

func Pagination(db *gorm.DB, req PageReq, count *int64, list any) (err error) {
	if err = db.Count(count).Error; err != nil {
		return
	}
	if *count == 0 {
		return
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	err = db.Offset(req.Limit * (req.Page - 1)).Limit(req.Limit).Find(list).Error

	return
}

func Callback(err0 error, fns ...func()) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			fn()
		}
	}
	return
}

func Callback1[T any](err0 error, t T, fns ...func(t T)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			fn(t)
		}
	}
	return
}

func Callback2[T1, T2 any](err0 error, t1 T1, t2 T2, fns ...func(t1 T1, t2 T2)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			fn(t1, t2)
		}
	}
	return
}

func Callback3[T1, T2, T3 any](err0 error, t1 T1, t2 T2, t3 T3, fns ...func(t1 T1, t2 T2, t3 T3)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			fn(t1, t2, t3)
		}
	}
	return
}

func CallbackErr(err0 error, fns ...func() (err error)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			if err = fn(); err != nil {
				return
			}
		}
	}
	return
}

func Callback1Err[T any](err0 error, t T, fns ...func(t T) (err error)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			if err = fn(t); err != nil {
				return
			}
		}
	}
	return
}

func Callback2Err[T1, T2 any](err0 error, t1 T1, t2 T2, fns ...func(t1 T1, t2 T2) (err error)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			if err = fn(t1, t2); err != nil {
				return
			}
		}
	}
	return
}

func Callback3Err[T1, T2, T3 any](err0 error, t1 T1, t2 T2, t3 T3, fns ...func(t1 T1, t2 T2, t3 T3) (err error)) (err error) {
	if err = err0; err != nil {
		return
	}
	for _, fn := range fns {
		if fn != nil {
			if err = fn(t1, t2, t3); err != nil {
				return
			}
		}
	}
	return
}

func Return[T any](t T, err0 error) (T, error) { return t, err0 }
