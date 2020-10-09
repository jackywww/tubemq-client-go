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

	"github.com/lubanproj/tubemq-client-go/client"
	"github.com/lubanproj/tubemq-client-go/constants"
	"github.com/lubanproj/tubemq-client-go/rebalance"
)

// Consumer defines a basic contract that all consumers need to follow
type Consumer interface {
	// Subscribe defines a subscription to the topic
	Subscribe(topic string, filterConds []string) error
}

// New create a default pull consumer
func New(config *Config) *PullConsumer {
	if err := checkConfig(config); err != nil {
		panic(err)
	}

	clientOpts := []client.Option {
		client.WithTarget(config.MasterAddress),
		client.WithTimeout(config.Timeout),
		client.WithNetwork(constants.DefaultNetworkType),
	}
	client := client.New(clientOpts ...)

	c := &PullConsumer{
		config:           config,
		client:			  client,
		rebalanceEvents:  make(chan *rebalance.Event, 100),
		rebalanceResults: make(chan *rebalance.Event, 100),
		rmtDataCache: &rmtDataCache{
			partitionMap:       &sync.Map{},
			partitionUsedMap:   &sync.Map{},
			partitionOffsetMap: &sync.Map{},
		},
	}

	return c
}

func checkConfig(config *Config) error {
	if config == nil {
		return errors.New("consumer config cannot be nil")
	}

	if config.MasterAddress == "" {
		return errors.New("consumer config err, master address cannot be empty")
	}

	if config.Group == "" {
		return errors.New("consumer config err, group cannot be empty")
	}

	if len(config.Topics) == 0 {
		return errors.New("consumer config err, topic cannot be empty")
	}

	return nil
}
