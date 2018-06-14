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
	"github.com/container-mgmt/messaging-library/pkg/client"
)

// ConnectionBuilder is an implementation of ConnectionBuilder interface
type ConnectionBuilder struct {
	connection Connection
}

// Build setter
func (c ConnectionBuilder) Build() client.Connection {
	return c.connection
}

// BrokerHost setter
func (c ConnectionBuilder) BrokerHost(h string) client.ConnectionBuilder {
	c.connection.BrokerHost = h

	return c
}

// BrokerPort setter
func (c ConnectionBuilder) BrokerPort(p int) client.ConnectionBuilder {
	c.connection.BrokerPort = p

	return c
}

// UserName setter
func (c ConnectionBuilder) UserName(n string) client.ConnectionBuilder {
	c.connection.UserName = n

	return c
}

// UserPassword setter
func (c ConnectionBuilder) UserPassword(p string) client.ConnectionBuilder {
	c.connection.UserPassword = p

	return c
}

// UseTLS setter
func (c ConnectionBuilder) UseTLS(t bool) client.ConnectionBuilder {
	c.connection.UseTLS = t

	return c
}

// InsecureTLS setter
func (c ConnectionBuilder) InsecureTLS(t bool) client.ConnectionBuilder {
	c.connection.InsecureTLS = t

	return c
}
