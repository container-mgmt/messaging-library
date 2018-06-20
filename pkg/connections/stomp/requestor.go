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
	"github.com/segmentio/ksuid"
)

// Requestor is an implementation of Requestor interface
// The stomp requestor is a specification of the connection interface
type Requestor struct {
	conn           *Connection
	subscription   *stomp.Subscription
	requestsQueue  string
	responsesQueue string

	// mapping between request ID and it's handler
	pendingRequests map[string]client.ResponseHandler
}

// NewRequestor creates a new requestor API to submit requests
func (c *Connection) NewRequestor(spec client.RequestorSpec) (r client.Requestor, err error) {

	// Check if we already subscibe to this destination,
	// We do not allow for multiple subscriptions for one destination.
	if _, ok := c.subscriptions[spec.ResponsesQueue]; ok {
		err = fmt.Errorf("Only one subscription per destination is allowed")
		return
	}

	// Subscribe to receive messages:
	subscription, err := c.connection.Subscribe(spec.ResponsesQueue, stomp.AckAuto)
	if err != nil {
		return
	}

	// add the subscription to the connection
	c.subscriptions[spec.ResponsesQueue] = subscription

	stompRequestor := &Requestor{
		conn:            c,
		requestsQueue:   spec.RequestsQueue,
		responsesQueue:  spec.ResponsesQueue,
		subscription:    subscription,
		pendingRequests: make(map[string]client.ResponseHandler, 0),
	}

	// wait for responses in the background
	go stompRequestor.waitForResponses()

	r = stompRequestor
	return
}

// Send sends a request to a specific destination
// request is the message request
// callback is the handler that will be called when the response is received
// requestID is returned upon successful send
func (r *Requestor) Send(request client.Message, callback client.ResponseHandler) (requestID string, err error) {

	// generate request uuid
	requestID = ksuid.New().String()

	// Add request fields to message
	request.Data["kind"] = "Request"
	request.Data["requestID"] = requestID
	request.Data["respondTo"] = r.responsesQueue

	// send the message
	// TODO: don't use this function and directly us the stomp API
	err = r.conn.Publish(request, r.requestsQueue)
	if err != nil {
		requestID = ""
	}

	// keep the handler in the pending requests map
	r.pendingRequests[requestID] = callback

	return
}

func (r *Requestor) waitForResponses() {
	var data client.MessageData
	for message := range r.subscription.C {
		// Try to unmarshal the byte array coming from the broker into a
		// message body of type map[string]interface{}
		err := json.Unmarshal(message.Body, &data)
		if err != nil {
			// log the error and ignore message
			glog.Warningf(
				"failed to unmarshall message received from destination %s. Ignoring",
				r.responsesQueue)
			continue
		}

		// Validate message is a response
		if kind, ok := data["kind"]; !ok || kind.(string) != "Response" {
			// ignore message
			glog.Warningf(
				"Message of non 'Response' kind received on response queue %s. Ignoring",
				r.responsesQueue)
			continue
		}

		// Parse requestId
		id, ok := data["requestID"]
		if !ok {
			// ignore message
			glog.Warningf(
				"Response missing 'requestID' field received on response queue %s. Ignoring",
				r.responsesQueue)
			continue
		}

		// Validate requestID
		callback, ok := r.pendingRequests[id.(string)]
		if !ok {
			// ignore message
			glog.Warningf(
				"Received response to non existing request id %s. Ignoring",
				id.(string))
			continue
		}

		// call the relevant response handler
		callback(
			client.Message{
				Data: data,
				Err:  message.Err},
			id.(string))

		// remove the pending request
		// TBD - in the future we would like to keep the resuest to allow
		// a series of responses for a single request.
		// Currently only one response is supported
		delete(r.pendingRequests, id.(string))
	}
}

// Close closes the Requestor
func (r *Requestor) Close() (err error) {
	err = r.conn.Unsubscribe(r.responsesQueue)
	r.subscription = nil
	return
}
