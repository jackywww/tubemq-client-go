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

package testdata

import (
	"time"

	"github.com/lubanproj/tubemq-client-go"
)

// NewConfig creates a Config with test data
func NewConfig() *tubeclient.Config {
	return &tubeclient.Config{
		Group: "test",
		MasterAddress: "10.215.128.83:8000,10.215.128.83:8000",
		Topics: []string{"test_1"},
		Timeout: 2 * time.Second,
	}
}
