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

import "github.com/gin-gonic/gin"

var engine = gin.New()

func GetApp(middlewares ...gin.HandlerFunc) *gin.Engine { engine.Use(middlewares...); return engine }

func GetAppWithGroup(prefix string, middlewares ...gin.HandlerFunc) *gin.RouterGroup {
	g := GetApp().Group(prefix)
	g.Use(middlewares...)
	return g
}
