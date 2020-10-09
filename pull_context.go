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

package tubeclient

import (
	"github.com/lubanproj/tubemq-client-go/errs"
	"github.com/lubanproj/tubemq-client-go/model"
	"github.com/lubanproj/tubemq-client-go/protocol"
)

// PullContext saved the PullConsumer global context information
type PullContext struct {
	partition      *model.Partition
	usedToken      int64
	lastConsumed   bool
	curOffset      int64
	messageList    []*protocol.Message
	success        bool
	errcode        int32
	errmsg         string
	confirmContext string
}

func (p *PullContext) setSuccResult(messageList []*protocol.Message) {
	p.messageList = messageList
	p.success = true
	p.errcode = errs.Success
	p.errmsg = "ok"
	p.messageList = messageList
}

func (p *PullContext) setFailResult(errcode int32, errmsg string) {
	p.success = false
	p.errcode = errcode
	p.errmsg = errmsg
}
