/**
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package errs

import (
	"fmt"
)

const (
	Success = 200

	Hb_No_Node = 411
	Duplicate_Partition = 412
	Certificate_Failure = 415

	ConsumeGroupForbidden = 450
	ServerConsumeSpeedLimit = 452
	ConsumeConnectForbidden = 455

	TokenInvalid = 501
	PayLoadOverLimit = 502
)


var (
	TokenInvalidError = New(TokenInvalid, "token invalid")
	PayLoadOverLimitError = New(PayLoadOverLimit, "payload overlimit")
)

// Error defines all errors in the framework
type Error struct {
	code int32
	msg string
}

func (e *Error) Error() string {
	if e == nil {
		return "success"
	}
	return fmt.Sprintf("code : %d, msg : %s", e.code, e.msg)
}

// New creates an Error
func New(code int32, msg string) error {
	return &Error {
		code: code,
		msg: msg,
	}
}