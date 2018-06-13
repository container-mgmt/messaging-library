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

// Package client contains the types and functions used to communicate with other services using
// queues and topics.
package client

var connections map[string]Connection
var connectionName string

// Register a new client connection
func Register(name string, connection Connection) (err error) {
	if connections == nil {
		connections = make(map[string]Connection, 0)
	}

	// Set last registered connection as default.
	connectionName = name

	// Add a new connection
	connections[name] = connection

	return
}

// Client represent a message broker client
type Client struct {
	Name string
}

// NewConnectionBuilder return a ConnectionBuilder for this connection
func (c *Client) NewConnectionBuilder() ConnectionBuilder {
	if c.Name == "" {
		c.Name = connectionName
	}
	return connections[c.Name].NewConnectionBuilder()
}

// Use specific client by name
func (c *Client) Use(name string) (err error) {
	c.Name = name

	return
}

// Open is
func (c Client) Open() error {
	if c.Name == "" {
		c.Name = connectionName
	}
	return connections[c.Name].Open()
}

// Close is
func (c Client) Close() error {
	if c.Name == "" {
		c.Name = connectionName
	}
	return connections[c.Name].Close()
}

// Publish a message to a topic
func (c Client) Publish(m Message, topic string) error {
	if c.Name == "" {
		c.Name = connectionName
	}
	return connections[c.Name].Publish(m, topic)
}

// Subscribe to a topic
func (c Client) Subscribe(topic string, callback func(m Message, topic string) error) error {
	if c.Name == "" {
		c.Name = connectionName
	}
	return connections[c.Name].Subscribe(topic, callback)
}
