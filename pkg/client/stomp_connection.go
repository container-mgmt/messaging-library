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

package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/go-stomp/stomp"
)

// StompConnection is an implementation of Connection interface
type StompConnection struct {
	// Global options:
	brokerHost      string
	brokerPort      int
	destinationName string
	userName        string
	userPassword    string
	useTLS          bool
	insecureTLS     bool
	connection      *stomp.Conn
}

// Connect creates a ne connection to the messaging broker.
func (c *StompConnection) Connect() (err error) {
	// Calculate the address of the server, as required by the Dial methods:
	brokerAddress := fmt.Sprintf("%s:%d", c.brokerHost, c.brokerPort)

	// Create the socket:
	var socket io.ReadWriteCloser
	if c.useTLS {
		socket, err = tls.Dial("tcp", brokerAddress, &tls.Config{
			ServerName:         c.brokerHost,
			InsecureSkipVerify: c.insecureTLS,
		})
		if err != nil {
			err = fmt.Errorf(
				"can't create TLS connection to host '%s' and port %d: %s",
				c.brokerHost,
				c.brokerPort,
				err.Error(),
			)
			return
		}
	} else {
		socket, err = net.Dial("tcp", brokerAddress)
		if err != nil {
			err = fmt.Errorf(
				"can't create TCP connection to host '%s' and port %d: %s",
				c.brokerHost,
				c.brokerPort,
				err.Error(),
			)
			return
		}
	}

	// Prepare the options:
	options := make([]func(*stomp.Conn) error, 0)
	if c.userName != "" {
		options = append(options, stomp.ConnOpt.Login(c.userName, c.userPassword))
	}

	// Create the STOMP connection:
	c.connection, err = stomp.Connect(socket, options...)
	if err != nil {
		err = fmt.Errorf(
			"can't create STOMP connection to host '%s' and port %d: %s",
			c.brokerHost,
			c.brokerPort,
			err.Error(),
		)
		return
	}

	return
}

// Close closes the connection, releasing all the resources that it uses. Once closed the
// connection can't be reused.
func (c *StompConnection) Close() (err error) {
	return c.connection.Disconnect()
}
