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
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/go-stomp/stomp"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

// Connection is an implementation of Connection interface
type Connection struct {
	// Global options:
	BrokerHost   string
	BrokerPort   int
	UserName     string
	UserPassword string
	UseTLS       bool
	InsecureTLS  bool
	connection   *stomp.Conn
}

// Open creates a new connection to the messaging broker.
func (c Connection) Open() (err error) {
	// Calculate the address of the server, as required by the Dial methods:
	brokerAddress := fmt.Sprintf("%s:%d", c.BrokerHost, c.BrokerPort)

	// Create the socket:
	var socket io.ReadWriteCloser
	if c.UseTLS {
		socket, err = tls.Dial("tcp", brokerAddress, &tls.Config{
			ServerName:         c.BrokerHost,
			InsecureSkipVerify: c.InsecureTLS,
		})
		if err != nil {
			err = fmt.Errorf(
				"can't create TLS connection to host '%s' and port %d: %s",
				c.BrokerHost,
				c.BrokerPort,
				err.Error(),
			)
			return
		}
	} else {
		socket, err = net.Dial("tcp", brokerAddress)
		if err != nil {
			err = fmt.Errorf(
				"can't create TCP connection to host '%s' and port %d: %s",
				c.BrokerHost,
				c.BrokerPort,
				err.Error(),
			)
			return
		}
	}

	// Prepare the options:
	var options []func(*stomp.Conn) error
	if c.UserName != "" {
		options = append(options, stomp.ConnOpt.Login(c.UserName, c.UserPassword))
	}

	// Create the STOMP connection:
	c.connection, err = stomp.Connect(socket, options...)
	if err != nil {
		err = fmt.Errorf(
			"can't create STOMP connection to host '%s' and port %d: %s",
			c.BrokerHost,
			c.BrokerPort,
			err.Error(),
		)
		return
	}

	return
}

// Close closes the connection, releasing all the resources that it uses. Once closed the
// connection can't be reused.
func (c Connection) Close() (err error) {
	return c.connection.Disconnect()
}

// NewConnectionBuilder gets a ConnectionBuilder for this connection
func (c Connection) NewConnectionBuilder() client.ConnectionBuilder {
	return ConnectionBuilder{}
}
