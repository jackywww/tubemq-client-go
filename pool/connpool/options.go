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

package connpool

import "time"

type Options struct {
	initialCap int   // initial capacity
	maxCap int    // max capacity
	idleTimeout time.Duration
	maxIdle int   // max idle connections
	dialTimeout time.Duration  // dial timeout
}

type Option func(*Options)

func WithInitialCap (initialCap int) Option {
	return func(o *Options) {
		o.initialCap = initialCap
	}
}

func WithMaxCap (maxCap int) Option {
	return func(o *Options) {
		o.maxCap = maxCap
	}
}


func WithMaxIdle (maxIdle int) Option {
	return func(o *Options) {
		o.maxIdle = maxIdle
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = idleTimeout
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *Options) {
		o.dialTimeout = dialTimeout
	}
}