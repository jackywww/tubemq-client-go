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

package model

import (
	"strings"

	"github.com/lubanproj/tubemq-client-go/constants"
)

// SubscribeInfo holds the consumer's subscription information
type SubscribeInfo struct {
	ConsumerID string
	Group string
	Partition *Partition
	OverTls bool
}

// BuildSubscribeInfo constructs a SubscribeInfo with a string in a given format
func BuildSubscribeInfo(subStr string) *SubscribeInfo {
	consumerStr := strings.Split(subStr, constants.SegmentSep)[0]
	brokerStr := strings.Split(subStr, constants.SegmentSep)[1]
	partStr := strings.Split(subStr, constants.SegmentSep)[2]

	consumerID := strings.Split(consumerStr, constants.GroupSep)[0]
	group := strings.Split(consumerStr, constants.GroupSep)[1]

	broker := BuildBroker(brokerStr)
	partition := BuildPartition(broker, partStr)

	return &SubscribeInfo{
		ConsumerID: consumerID,
		Group: group,
		Partition: partition,
	}
}
