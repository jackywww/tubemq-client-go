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

import (
	"context"
	"github.com/lubanproj/tubemq-client-go/balance"
	"github.com/lubanproj/tubemq-client-go/pool/connpool"

	"github.com/lubanproj/tubemq-client-go/codec"
	"github.com/lubanproj/tubemq-client-go/transport"
)

// Client interface provides the ability for clients and servers to communicate
type Client interface {
	Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error
}

type defaultClient struct {
	// Transport is the default ClientTransport
	Transport transport.ClientTransport
	opts *Options
}

// New create a default Client
func New(opt ...Option) Client {

	opts := &Options{}
	for _, o := range opt {
		o(opts)
	}

	return &defaultClient {
		opts: opts,
		Transport: transport.DefaultClientTransport,
	}
}

func (c *defaultClient) Invoke(ctx context.Context, req interface{}, rsp interface{}, opts ...Option) error {

	serialization := codec.GetSerialization(c.opts.serialization)
	payload, err := serialization.Marshal(req)

	request := &codec.RequestWrapper{
		Payload: payload,
	}

	reqbuf, err := codec.DefaultCodec.Encode(request)
	if err != nil {
		return err
	}

	clientTransportOpts := []transport.ClientTransportOption {
		transport.WithClientNetwork(c.opts.network),
		transport.WithTimeout(c.opts.timeout),
		transport.WithClientTarget(c.opts.target),
		transport.WithClientPool(connpool.DefaultPool),
		transport.WithBalancer(balance.GetBalancer(c.opts.balancer)),
	}

	frame, err := c.Transport.Send(ctx, reqbuf, clientTransportOpts ...)
	if err != nil {
		return err
	}
	rspbuf, err := codec.DefaultCodec.Decode(frame)
	if err != nil {
		return err
	}

	return serialization.Unmarshal(rspbuf, rsp)

}