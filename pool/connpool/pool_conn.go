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

import (
	"errors"
	"net"
	"sync"
	"time"
)

var (
	ErrConnClosed = errors.New("connection closed ... ")
)

type PoolConn struct {
	net.Conn
	c *channelPool
	unusable bool		// if unusable is true, the conn should be closed
	mu sync.RWMutex
	t time.Time  // connection idle time
	dialTimeout time.Duration // connection timeout duration
}

// overwrite conn Close for connection reuse
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
	}

	// reset connection deadline
	p.Conn.SetDeadline(time.Time{})

	return p.c.Put(p)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

func (p *PoolConn) Read(b []byte) (int, error) {
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Read(b)
	if err != nil {
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n, err
}

func (p *PoolConn) Write(b []byte) (int, error) {
	if p.unusable {
		return 0, ErrConnClosed
	}
	n, err := p.Conn.Write(b)
	if err != nil {
		p.MarkUnusable()
		p.Conn.Close()
	}
	return n, err
}

func (c *channelPool) wrapConn(conn net.Conn) *PoolConn {
	p := &PoolConn {
		c : c,
		t : time.Now(),
		dialTimeout: c.dialTimeout,
	}
	p.Conn = conn
	return p
}

