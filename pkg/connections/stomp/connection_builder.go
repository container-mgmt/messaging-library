/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stomp

import (
	"github.com/go-stomp/stomp"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

// ConnectionBuilder is a helper struct for building connections.
type ConnectionBuilder struct {
	BrokerHost   string
	BrokerPort   int
	UserName     string
	UserPassword string
	UseTLS       bool
	InsecureTLS  bool
}

// NewConnection builds and initiate a new connection object.
//
// Example:
//   c, err = stomp.NewConnection(&stomp.ConnectionBuilder{
//   	BrokerHost:   brokerHost,
//   	BrokerPort:   brokerPort,
//   	UserName:     userName,
//   	UserPassword: userPassword,
//   	UseTLS:       useTLS,
//   	InsecureTLS:  insecureTLS,
//   })
//   if err != nil {
//   	glog.Errorf(
//   		"Can't create a new connection to host '%s': %s",
//   		brokerHost,
//  		err.Error(),
//   	)
//   	return
//  }
func NewConnection(b *ConnectionBuilder) (connection client.Connection, err error) {
	// Init Host and port values if found zero values.
	if b.BrokerHost == "" {
		b.BrokerHost = "127.0.0.1"
	}
	if b.BrokerPort == 0 {
		b.BrokerPort = 61613
	}

	// Create the connection object.
	stompConnection := &Connection{
		brokerHost:   b.BrokerHost,
		brokerPort:   b.BrokerPort,
		userName:     b.UserName,
		userPassword: b.UserPassword,
		useTLS:       b.UseTLS,
		insecureTLS:  b.InsecureTLS,
	}

	// Init connection subscriptions.
	stompConnection.subscriptions = make(map[string]*stomp.Subscription, 0)

	connection = stompConnection
	return
}
