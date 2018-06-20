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

// RequestHandler is called when a new request is received
// m is the request message
// requestID is the id of the request
type RequestHandler func(request Message) (response Message, err error)

// ResponderSpec is a helper struct for building responders.
type ResponderSpec struct {
	RequestsQueue string
	Callback      RequestHandler
}

// Responder is a request server interface
type Responder interface {
	Close() error
}
