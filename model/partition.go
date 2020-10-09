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
	"strconv"
	"strings"

	"github.com/lubanproj/tubemq-client-go/constants"
	"github.com/lubanproj/tubemq-client-go/protocol"
)

// Partition TODO
type Partition struct {
	// PartitionID is the partition ID
	PartitionID int32
	// PartitionKey is the partition key
	PartitionKey string
	// Topic is the topic of partition
	Topic string
	// Broker holds all the information of a Broker
	Broker *Broker
	// RetryTimes is the times of retries
	RetryTimes int
}

// TopicAndPartition TODO
type TopicAndPartition struct {
	// Topic
	Topic string
	// Partition is the partition of a topic
	Partition int32
}

type partitionTopicInfo struct {
	Topic         string
	Partition     int32
	Buffer        *protocol.Message
	FetchedOffset int64
}

// PartitionExt includes some additional information about the partition
type PartitionExt struct {
	Partition *Partition
	RecvMsgSize int64
	RecvMsgPerMin int64
	IsLastPackConsumed bool
}

// BuildPartition constructs a Partition with a string in a given format
func BuildPartition(broker *Broker, partStr string) *Partition {
	topic := strings.Split(partStr, constants.AttrSep)[0]
	partitionID, _ := strconv.Atoi(strings.Split(partStr, constants.AttrSep)[1])

	return &Partition {
		PartitionID: int32(partitionID),
		Topic: topic,
		Broker: broker,
	}
}

func (p *PartitionExt) GetAndResetLastPackConsumed() bool {
	isLastPackConsumed := p.IsLastPackConsumed
	isLastPackConsumed = false
	return isLastPackConsumed
}