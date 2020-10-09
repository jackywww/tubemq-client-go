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

package constants

import "time"

const (
	SegmentSep = "#"
	AttrSep    = ":"
	GroupSep   = "@"
	LogSegSep  = ";"
	ArraySep   = ","

	DefaultBrokerTlsPort = 8124

	RpcProtocolBeginToken uint64 = 0xFF7FF4FE
	DataBlockLenth        = 8192

	DefaultPayloadLength = 8192
	MaxPayloadLength     = 16 * 1024 * 1024

	// FrameHeadLen defines the length of frame header,
	// you can refer to https://tubemq.apache.org/zh-cn/docs/client_rpc.html for details
	FrameHeaderLenth = 12

	ConfigMaxFilterCondLength = 500
	ConfigHeartBeatPeriod     = 13 * time.Second
	ConfigRegisterRetryTimes = 5
)
