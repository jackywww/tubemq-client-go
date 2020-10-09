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
)

// Broker TODO
type Broker struct {
	// ID is the broker id
	ID int
	// Host e.g. : 127.0.0.1
	Host string
	// Port, e.g. : 6379
	Port int
	// Addr is the address of a broker
	Addr string
	// EnableTls enable tls or not
	EnableTls bool
	// TLSPort
	TLSPort int
}

// BuildBrokerInfo constructs a Broker with a string in a given format
func BuildBroker(brokerStr string) *Broker {
	strList := strings.Split(brokerStr, constants.AttrSep)
	brokerID, _ := strconv.Atoi(strList[0])
	host := strList[1]
	port, _ := strconv.Atoi(strList[2])

	return &Broker {
		ID : brokerID,
		Host: host,
		Port: port,
		TLSPort: port,
	}
}