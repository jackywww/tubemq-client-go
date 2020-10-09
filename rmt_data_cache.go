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
	"errors"
	"sync"
	"time"

	"github.com/lubanproj/tubemq-client-go/model"
)

// RmtDataCache cache consumption status information for remote messages
type RmtDataCache interface {
}

type rmtDataCache struct {
	indexPartition     chan string
	partitionMap       *sync.Map
	partitionUsedMap   *sync.Map
	partitionOffsetMap *sync.Map
}

func (r *rmtDataCache) selectPartition() (*model.PartitionExt, error) {
	partitionKey, ok := <-r.indexPartition
	if !ok {
		return nil, errors.New("no partition available, please retry later")
	}

	v, ok := r.partitionMap.Load(partitionKey)
	if !ok {
		return nil, errors.New("no valid partition available, please retry later")
	}
	partitionExt, ok := v.(*model.PartitionExt)
	if !ok {
		return nil, errors.New("no valid partition available, please retry later")
	}

	r.partitionUsedMap.Store(partitionKey, time.Now().Unix())
	return partitionExt, nil
}

func (r *rmtDataCache) isPartitionInUse(partitionKey string, usedToken int64) bool {
	if v, ok := r.partitionMap.Load(partitionKey); ok {
		if partitionExt, ok := v.(*model.PartitionExt); ok {
			token, _ := r.partitionUsedMap.Load(partitionExt.Partition.PartitionKey)
			return token != nil && token == usedToken
		}
	}

	return false
}

func (r *rmtDataCache) getPartition(partitionKey string) *model.PartitionExt {
	if v, ok := r.partitionMap.Load(partitionKey); ok {
		if partitionExt, ok := v.(*model.PartitionExt); ok {
			return partitionExt
		}
	}

	return nil
}
