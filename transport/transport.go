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
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/lubanproj/tubemq-client-go/constants"
	"github.com/lubanproj/tubemq-client-go/errs"
)

// ClientTransport defines the criteria that all client transport layers
// need to support
type ClientTransport interface {
	// send requests
	Send(context.Context, []byte, ...ClientTransportOption) ([]byte, error)
}

type clientTransport struct {
	opts *ClientTransportOptions
}

type connWrapper struct {
	net.Conn
	framer Framer
}

func wrapConn(rawConn net.Conn) *connWrapper {
	return &connWrapper{
		Conn:   rawConn,
		framer: NewFramer(),
	}
}

// Framer defines the reading of data frames from a data stream
type Framer interface {
	// read a full frame
	ReadFrame(net.Conn) ([]byte, error)
}

type framer struct {
	buffer  []byte
	counter int // to prevent the dead loop
}

// Create a Framer
func NewFramer() Framer {
	return &framer{
		buffer: make([]byte, constants.DefaultPayloadLength),
	}
}

func (f *framer) Resize() {
	f.buffer = make([]byte, len(f.buffer)*2)
}

func (f *framer) ReadFrame(conn net.Conn) ([]byte, error) {
	frameHeader := make([]byte, constants.FrameHeaderLenth)
	if num, err := io.ReadFull(conn, frameHeader); num != constants.FrameHeaderLenth || err != nil {
		return nil, err
	}

	token := binary.BigEndian.Uint64(frameHeader[0:4])
	if token != constants.RpcProtocolBeginToken {
		return nil, errs.TokenInvalidError
	}

	listSize := binary.BigEndian.Uint32(frameHeader[8:12])
	dataLen := listSize * constants.DefaultPayloadLength
	if dataLen > constants.MaxPayloadLength {
		return nil, errs.PayLoadOverLimitError
	}
	for uint32(len(f.buffer)) < dataLen && f.counter <= 12 {
		f.buffer = make([]byte, len(f.buffer)*2)
		f.counter++
	}
	if num, err := io.ReadFull(conn, f.buffer[:dataLen]); num != constants.FrameHeaderLenth || err != nil {
		return nil, err
	}

	return append(frameHeader, f.buffer[:dataLen]...), nil
}


// DefaultClientTransport creates a default ClientTransport
var DefaultClientTransport = New()

// New Use the singleton pattern to create a ClientTransport
var New = func() ClientTransport {
	return &clientTransport{
		opts: &ClientTransportOptions{},
	}
}

func (c *clientTransport) Send(ctx context.Context, req []byte, opts ...ClientTransportOption) ([]byte, error) {

	for _, o := range opts {
		o(c.opts)
	}

	addr := c.opts.Balancer.Balance(c.opts.Target)

	conn, err := c.opts.Pool.Get(ctx, c.opts.Network, addr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	sendNum := 0
	num := 0
	for sendNum < len(req) {
		num, err = conn.Write(req[sendNum:])
		if err != nil {
			return nil, err
		}
		sendNum += num

		if err = isDone(ctx); err != nil {
			return nil, err
		}
	}

	// parse frame
	wrapperConn := wrapConn(conn)
	frame, err := wrapperConn.framer.ReadFrame(conn)
	if err != nil {
		return nil, err
	}

	return frame, err
}

func isDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}
