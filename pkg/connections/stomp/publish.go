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
	"encoding/json"
	"github.com/go-stomp/stomp"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

// PublishByteArray sends a byte array to the messaging server, which in turn
// sends the message to the specified destination. If the messaging server fails to
// receive the message for any reason, the connection will close.
func (c *Connection) PublishByteArray(contentType string, body []byte, destination string) (err error) {
	err = c.connection.Send(
		destination,
		contentType,
		body,
		stomp.SendOpt.Header("persistent", "true"),
	)
	return
}

// Publish sends a message to the messaging server, which in turn sends the
// message to the specified destination. If the messaging server fails to
// receive the message for any reason, the connection will close.
func (c *Connection) Publish(m client.Message, destination string) (err error) {
	var body []byte

	// Our default contentType is "application/json"
	contentType := m.ContentType
	if contentType == "" {
		contentType = "application/json"
	}

	// Check if we have a byteArray content, if we do, we will overide the
	// object abstraction mechanism, and send the byteArray as a byte array.
	if _, ok := m.Data["byteArray"]; ok {
		// Check the content of byteArray, if it's really a byte array, send it
		// as a byte array to server.
		switch m.Data["byteArray"].(type) {
		case []byte:
			body = m.Data["byteArray"].([]byte)
			err = c.PublishByteArray(contentType, body, destination)
			return
		}
	}

	// Marshal the message body (type: client.MessageData) into a byte array.
	body, err = json.Marshal(m.Data)
	if err != nil {
		return
	}

	err = c.PublishByteArray(contentType, body, destination)
	return
}
