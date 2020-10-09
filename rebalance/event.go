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

package rebalance

// Event is the event received from the heartbeat packet of the master or broker
type Event struct {
	RebalanceID int64
	Type EventType
	Status EventStatus
	SubscribeStrList []string
}

const (
	// Connect to some broker
	Connect EventType = 1
	// Disconnect from some broker
	Disconnect EventType = 2
	// Report current status
	Report EventType = 3
	// Refresh update whole producer published topic info
	Refresh EventType = 4
	// StopRebalance stop rebalance thread
	StopRebalance EventType = 5
	// OnlyConnect only connect to some broker
	OnlyConnect EventType = 10
	// OnlyDisConnect only disconnect to some broker
	OnlyDisConnect EventType = 20
	// Unknown event type
	UnknownType EventType = -1
)

const (
	// Preparing indicates that the event is being prepared
	Preparing EventStatus = 0
	// Processing indicates that the event is processing
	Processing EventStatus = 1
	// Done indicates that the event is done
	Done EventStatus = 2
	// UnknownStatus indicates unknown event
	UnknownStatus EventStatus = 3
	// Failed indicates failed event
	Failed EventStatus = 4
)

var (
	eventTypeNames = map[EventType]string {
		Connect: "Connect",
		Disconnect: "Disconnect",
		Report: "Report",
		Refresh: "Refresh",
		StopRebalance: "StopRebalance",
		OnlyConnect: "OnlyConnect",
		OnlyDisConnect: "OnlyDisConnect",
		UnknownType: "Unknown",
	}

	eventStatusNames = map[EventStatus]string {
		Preparing : "Preparing",
		Processing : "Processing",
		Done : "Done",
		UnknownStatus : "UnknownStatus",
		Failed : "Failed",
	}
)

// EventType is the type of an Event
type EventType int32

func (e EventType) String() string {
	if name := eventTypeNames[e]; name != "" {
		return name
	}
	return "UnknownType"
}

// EventStatus is the status of an Event
type EventStatus int32

func (e EventStatus) String() string {
	if name := eventStatusNames[e]; name != "" {
		return name
	}

	return "UnknownStatus"
}