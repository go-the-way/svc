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
	"sync"

	"github.com/gin-gonic/gin"
)

type (
	tokenCache[T any] interface {
		Add(token string, data ...T)
		Get(token string) (data T)
		GetCtx(ctx *gin.Context) (data T)
		Have(token string) (have bool)
		HaveCtx(ctx *gin.Context) (have bool)
		Delete(token string)
	}
	memTokenCache[T any] struct {
		*sync.RWMutex
		tokenHeader string
		m           map[string]T
	}
)

func newMemTokenCache[T any]() tokenCache[T] {
	return newMemTokenCacheTokenHeader[T]("Token")
}

func newMemTokenCacheTokenHeader[T any](tokenHeader string) tokenCache[T] {
	return &memTokenCache[T]{RWMutex: &sync.RWMutex{}, tokenHeader: tokenHeader, m: map[string]T{}}
}

func (m *memTokenCache[T]) Add(token string, data ...T) {
	m.Lock()
	defer m.Unlock()
	var d T
	if len(data) > 0 {
		d = data[0]
	}
	m.m[token] = d
}

func (m *memTokenCache[T]) Get(token string) (data T) {
	m.RLock()
	defer m.RUnlock()
	data, _ = m.m[token]
	return
}

func (m *memTokenCache[T]) GetCtx(ctx *gin.Context) (data T) { return m.Get(m.ctx0(ctx)) }

func (m *memTokenCache[T]) Have(token string) (have bool) {
	m.RLock()
	defer m.RUnlock()
	_, have = m.m[token]
	return
}

func (m *memTokenCache[T]) ctx0(ctx *gin.Context) string {
	return ctx.GetHeader("Token")
}

func (m *memTokenCache[T]) HaveCtx(ctx *gin.Context) (have bool) { return m.Have(m.ctx0(ctx)) }

func (m *memTokenCache[T]) Delete(token string) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, token)
	return
}
