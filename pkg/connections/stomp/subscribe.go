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
	"github.com/go-stomp/stomp"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

// Subscribe subscribes to a destination
func (c *Connection) Subscribe(destination string, callback func(m client.Message, destination string) error) (err error) {
	var subscription *stomp.Subscription

	// Receive messages:
	subscription, err = c.connection.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		return
	}

	// Wait for messages:
	for message := range subscription.C {
		callback(
			client.Message{
				Body: string(message.Body),
				Err:  message.Err},
			destination)
	}

	return
}
