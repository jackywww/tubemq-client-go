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

package client

import "time"

// Options defines the client call parameters
type Options struct {
	network string // network type, e.g.:  tcp、udp
	timeout time.Duration // timeout
	serialization string // serialization type,
	target string  // target address, e.g.: 127.0.0.1:8000
	balancer string  // balancer name, e.g.: random
}

// Option provides the modification entry for Options
type Option func(*Options)

// WithNetwork set the network type for client
func WithNetwork(network string) Option {
	return func(o *Options) {
		o.network = network
	}
}

// WithTimeout set the timeout for client
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

// WithSerialization set the serialization for client
func WithSerialization(serialization string) Option {
	return func(o *Options) {
		o.serialization = serialization
	}
}

// WithTarget set the target for client
func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

// WithBalancer set the balancer for client
func WithBalancer(balancer string) Option {
	return func(o *Options) {
		o.balancer = balancer
	}
}