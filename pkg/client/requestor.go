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

// ResponseHandler is called when a response to a request is received
// m is the response message
// requestID is the id of the request this message responds
type ResponseHandler func(response Message, requestID string) error

// RequestorSpec is a helper struct for building requestors.
type RequestorSpec struct {
	RequestsQueue  string
	ResponsesQueue string
}

// Requestor is a specification of publish/subscribe mechanism
// It allows sending direct response and supply a callback
type Requestor interface {
	Send(request Message, callback ResponseHandler) (requestID string, err error)
	Close() error
}
