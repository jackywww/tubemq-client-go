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

package balance

import (
	"math/rand"
	"strings"
	"time"

	"github.com/lubanproj/tubemq-client-go/constants"
)

// Balancer gets an Address through a given target
type Balancer interface {
	Balance(string) string
}

var balancerMap = make(map[string]Balancer, 0)

const (
	Random = "random"
)

func init() {
	RegisterBalancer(Random, DefaultBalancer)
}

// DefaultBalancer is the default balancer using using a random strategy
var DefaultBalancer = newRandomBalancer()

func RegisterBalancer(name string, balancer Balancer) {
	if balancerMap == nil {
		balancerMap = make(map[string]Balancer)
	}
	balancerMap[name] = balancer
}

// GetBalancer gets a balancer by a given balancer name
func GetBalancer(name string) Balancer {
	if balancer, ok := balancerMap[name]; ok {
		return balancer
	}
	return DefaultBalancer
}

func newRandomBalancer() *randomBalancer {
	return &randomBalancer{}
}

type randomBalancer struct {

}

func (r *randomBalancer) Balance(target string) string {
	nodes := strings.Split(target, constants.ArraySep)

	if target == "" || len(nodes) == 0 {
		return ""
	}

	rand.Seed(time.Now().Unix())
	num := rand.Intn(len(nodes))
	return nodes[num]
}
