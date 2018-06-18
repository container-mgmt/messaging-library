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

// MessageData is the message payload data type.
type MessageData interface{}

// Message represents a message sent or received by a connection.
// In most cases a message corresponds to a single message sent or received by
// a connection. If, however, the Err field is non-nil, then the message
// corresponds to an ERROR frame, or a connection error between the client
// and the server.
type Message struct {
	// The message body, which is an a string.
	// The ContentType indicates the format of this body.
	Data MessageData // Content of message

	// MIME content type.
	ContentType string // MIME of the message, usually "text/plain"

	// Indicates whether an error was received on the subscription.
	// The error will contain details of the error. If the server
	// sent an ERROR frame, then the Data, ContentType and Header fields
	// will be populated according to the contents of the ERROR frame.
	Err error
}
