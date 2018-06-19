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
//
// In this library we use the term destinations to describe both queues and topics.
package client

// SubscriptionCallback is the callback function type used for subscription callback.
// The callback function is used when subscribing to a destination, it is the
// function that will triger in the event of a message or an error frame.
//
// For example:
//   func callback(message client.Message, destination string) (err error) {
//  	if message.Err != nil {
//  		err = message.Err
//  		glog.Errorf(
//  			"Received error from destination '%s': %s",
//  			destinationName,
//  			err.Error(),
//  		)
//  		return
//  	}
//
//  	glog.Infof(
//  		"Received message from destination '%s':\n%v",
//  		destination,
//  		message.Data,
//  	)
//  	return
//  }
type SubscriptionCallback func(m Message, destination string) error

// Connection represents the logical connection between the program and the messaging system. This
// logical connection may correspond to one or multiple physical connections, depending on the
// underlying protocol and implementation.
//
// The connection may consume expensive resources, like TCP connections, or file descriptors, so it
// is important to reuse it as much as possible, and to close it once it is no longer needed.
//
// For implementation example see:
//   https://godoc.org/github.com/container-mgmt/messaging-library/pkg/connections/stomp
type Connection interface {
	// Open creates a new connection to the messaging broker.
	Open() error

	// Close closes the connection, releasing all the resources that it uses. Once closed the
	// connection can't be reused.
	Close() error

	// Publish sends a message to the messaging server, which in turn sends the
	// message to the specified destination. If the messaging server fails to
	// receive the message for any reason, the connection will close.
	// e.g.
	//   // The next lines will send a MessageData{} object to the server.
	//   err = c.Publish(
	//     client.Message{
	//       Data: client.MessageData{
	//         "some-key": "some-value",
	//         "other-key": "other-value",
	//       },
	//       ContentType: "application/json", // Default is "application/json"
	//     },
	//     "queue-name",
	//   )
	//
	// the function check if we have a byteArray key, if we do we will overide the
	// object abstraction mechanism, and send the byteArray to the server as is.
	// e.g.
	//   // The next lines will send a byte array to the server.
	//   data := client.MessageData{
	//     "byteArray": []byte("\"data\": \"message\""),
	//   }
	Publish(m Message, destination string) error

	// Subscribe creates a subscription on the messaging server.
	// The subscription has a destination, and messages sent to that destination
	// will be received by this subscription.
	//
	// Once a message or an error is received, the callback function will be trigered.
	Subscribe(destination string, callback SubscriptionCallback) error

	// Unsubscribe unsubscribes from a destination
	Unsubscribe(destination string) error
}
