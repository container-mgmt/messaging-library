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

// Connection represents the logical connection between the program and the messaging system. This
// logical connection may correspond to one or multiple physical connections, depending on the
// underlying protocol and implementation.
//
// The connection may consume expensive resources, like TCP connections, or file descriptors, so it
// is important to reuse it as much as possible, and to close it once it is no longer needed.
type Connection interface {
	// Open creates a new connection to the messaging broker.
	Open() error
	// Close closes the connection, releasing all the resources that it uses. Once closed the
	// connection can't be reused.
	Close() error

	// Publish a message to a topic
	Publish(m Message, topic string) error

	// Subscribe subscribes to a topic
	Subscribe(topic string, callback func(m Message, topic string) error) error
}
