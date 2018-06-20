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

// Package stomp contains the types and functions used to communicate with other services using
// queues and topics.
package stomp

import (
	"encoding/json"
	"fmt"

	"github.com/container-mgmt/messaging-library/pkg/client"
	"github.com/go-stomp/stomp"
	"github.com/golang/glog"
)

// Responder is an implementation of Responder interface
// The stomp responder is a specification of the connection interface
type Responder struct {
	conn          *Connection
	subscription  *stomp.Subscription
	requestsQueue string
	callback      client.RequestHandler
}

// NewResponder created a new responder with a specific destination
func (c *Connection) NewResponder(spec client.ResponderSpec) (r client.Responder, err error) {
	// Check if we already subscibe to this destination,
	// We do not allow for multiple subscriptions for one destination.
	if _, ok := c.subscriptions[spec.RequestsQueue]; ok {
		err = fmt.Errorf("Only one subscription per destination is allowed")
		return
	}

	// Subscribe to receive messages:
	subscription, err := c.connection.Subscribe(spec.RequestsQueue, stomp.AckAuto)
	if err != nil {
		return
	}

	// add the subscription to the connection
	c.subscriptions[spec.RequestsQueue] = subscription

	stompResponder := &Responder{
		conn:          c,
		requestsQueue: spec.RequestsQueue,
		subscription:  subscription,
		callback:      spec.Callback,
	}

	// wait for requests in the background
	go stompResponder.waitForRequests()

	r = stompResponder
	return
}

func (r *Responder) waitForRequests() {
	var data client.MessageData
	for message := range r.subscription.C {
		// Try to unmarshal the byte array coming from the broker into a
		// message body of type map[string]interface{}
		err := json.Unmarshal(message.Body, &data)
		if err != nil {
			// log the error and ignore message
			glog.Warningf(
				"failed to unmarshall message received from destination %s. Ignoring",
				r.requestsQueue)
			continue
		}

		// Validate message is a request
		if kind, ok := data["kind"]; !ok || kind.(string) != "Request" {
			// ignore message
			glog.Warningf(
				"Message of non 'Request' kind received on requests queue %s. Ignoring",
				r.requestsQueue)
			continue
		}

		// Parse requestId
		id, ok := data["requestID"]
		if !ok {
			// ignore message
			glog.Warningf(
				"Request missing 'requestID' field received on requests queue %s. Ignoring",
				r.requestsQueue)
			continue
		}

		// Parse respondTo
		respondTo, ok := data["respondTo"]
		if !ok {
			// ignore message
			glog.Warningf(
				"Request missing 'respondTo' field received on requests queue %s. Ignoring",
				r.requestsQueue)
			continue
		}

		// call callback function
		response, err := r.callback(
			client.Message{
				Data: data,
				Err:  message.Err})

		if err != nil {
			continue
		}

		// Add response kind field to message
		response.Data["kind"] = "Response"

		// Add requestID field to message
		response.Data["requestID"] = id.(string)

		// publish the response
		err = r.conn.Publish(response, respondTo.(string))
		continue
	}
}

// Close closes the Responder
func (r *Responder) Close() (err error) {
	err = r.conn.Unsubscribe(r.requestsQueue)
	r.subscription = nil
	return
}
