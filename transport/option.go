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

package transport

import (
	"time"

	"github.com/lubanproj/tubemq-client-go/balance"
	"github.com/lubanproj/tubemq-client-go/pool/connpool"
)

// ClientTransportOptions includes all ClientTransport parameter options
type ClientTransportOptions struct {
	Target      string
	Network     string
	Pool        connpool.Pool
	Timeout     time.Duration
	Balancer 	balance.Balancer
}

// ClientTransportOption Use the Options mode to wrap the ClientTransportOptions
type ClientTransportOption func(*ClientTransportOptions)

// WithClientTarget returns a ClientTransportOption which sets the value for target
func WithClientTarget(target string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Target = target
	}
}

// WithClientNetwork returns a ClientTransportOption which sets the value for network
func WithClientNetwork(network string) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Network = network
	}
}

// WithClientPool returns a ClientTransportOption which sets the value for pool
func WithClientPool(pool connpool.Pool) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Pool = pool
	}
}

// WithTimeout returns a ClientTransportOption which sets the value for timeout
func WithTimeout(timeout time.Duration) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Timeout = timeout
	}
}

// WithBalancer returns a ClientTransportOption which sets the value for balancer
func WithBalancer(balancer balance.Balancer) ClientTransportOption {
	return func(o *ClientTransportOptions) {
		o.Balancer = balancer
	}
}

