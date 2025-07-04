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

var (
	_ = FSFileAdapter
	_ = FSDirAdapter
	_ = AesEncrypt
	_ = AesDecrypt
	_ = GetAppWithGroup
	_ = Callback
	_ = Callback1[any]
	_ = Callback2[any, any]
	_ = Callback3[any, any, any]
	_ = CallbackErr
	_ = Callback1Err[any]
	_ = Callback2Err[any, any]
	_ = Callback3Err[any, any, any]
	_ = TransformCallback[any]
	_ = Cors
	_ = Decryption
	_ = WebTry
	_ = Pagination
	_ = Recover
	_ = MD5
	_ = IfGt0Func[byte]
	_ = IfNotEmptyFunc
	_ = TimeNowStr
	_ = TimeNowNum
	_ = TimeNowUnix
	_ = TimeNowMilli
	_ = ParseTime
	_ = FormatTime
	_ = Plugins
	_ = Return[any]
	_ = Uri
	_ = UriReq[any]
	_ = UriResp[any]
	_ = UriReqResp[any, any]
	_ = Query
	_ = QueryReq[any]
	_ = QueryResp[any]
	_ = QueryReqResp[any, any]
	_ = Body
	_ = BodyReq[any]
	_ = BodyResp[any]
	_ = BodyReqResp[any, any]
	_ = Form
	_ = FormReq[any]
	_ = FormResp[any]
	_ = FormReqResp[any, any]
	_ = ValidatorLangSupport
	_ = ValidatorLangFunc
	_ = MemTokenCache[any]
	_ = WriteMessageJSON
	_ = WriteSuccessOrErrorJSON
)
