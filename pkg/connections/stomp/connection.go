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

// Package stomp contains an implementation of a messaging-library/pkg/client
// connection object used to communicate with a STOMP broker server.
//
// https://godoc.org/github.com/container-mgmt/messaging-library/pkg/client
package stomp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/go-stomp/stomp"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

// Connection represents the logical connection between the program and the messaging system. This
// logical connection may correspond to one or multiple physical connections, depending on the
// underlying protocol and implementation.
//
// The connection may consume expensive resources, like TCP connections, or file descriptors, so it
// is important to reuse it as much as possible, and to close it once it is no longer needed.
//
// Connection is an implementation of Connection interface:
//   https://godoc.org/github.com/container-mgmt/messaging-library/pkg/client#Connection
type Connection struct {
	subscriptions map[string]*stomp.Subscription
	connection    *stomp.Conn
}

// NewConnection builds and initiate a new connection object.
//
// Example:
//   c, err = stomp.NewConnection(&client.ConnectionSpec{
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
func NewConnection(spec *client.ConnectionSpec) (connection client.Connection, err error) {
	// Init Host and port values if found zero values.
	brokerHost := spec.BrokerHost
	if brokerHost == "" {
		brokerHost = "127.0.0.1"
	}
	brokerPort := spec.BrokerPort
	if brokerPort == 0 {
		brokerPort = 61613
	}

	// Create the connection object.
	stompConnection := new(Connection)

	// Init connection subscriptions.
	stompConnection.subscriptions = make(map[string]*stomp.Subscription, 0)

	// Calculate the address of the server, as required by the Dial methods:
	brokerAddress := fmt.Sprintf("%s:%d", brokerHost, brokerPort)

	// Create the socket:
	var socket io.ReadWriteCloser
	if spec.UseTLS {
		socket, err = tls.Dial("tcp", brokerAddress, &tls.Config{
			ServerName:         brokerHost,
			InsecureSkipVerify: spec.InsecureTLS,
		})
		if err != nil {
			err = fmt.Errorf(
				"can't create TLS connection to host '%s' and port %d: %s",
				brokerHost,
				brokerPort,
				err.Error(),
			)
			return
		}
	} else {
		socket, err = net.Dial("tcp", brokerAddress)
		if err != nil {
			err = fmt.Errorf(
				"can't create TCP connection to host '%s' and port %d: %s",
				brokerHost,
				brokerPort,
				err.Error(),
			)
			return
		}
	}

	// Prepare the options:
	var options []func(*stomp.Conn) error
	if spec.UserName != "" {
		options = append(options, stomp.ConnOpt.Login(spec.UserName, spec.UserPassword))
	}

	// Create the STOMP connection:
	stompConnection.connection, err = stomp.Connect(socket, options...)
	if err != nil {
		err = fmt.Errorf(
			"can't create STOMP connection to host '%s' and port %d: %s",
			brokerHost,
			brokerPort,
			err.Error(),
		)
		return
	}

	// Return the created connection object:
	connection = stompConnection

	return
}

// Close closes the connection, releasing all the resources that it uses. Once closed the
// connection can't be reused.
func (c *Connection) Close() (err error) {
	// Sanity check connection.
	if c.connection == nil {
		err = fmt.Errorf("Connection is closed")
		return
	}

	return c.connection.Disconnect()
}
