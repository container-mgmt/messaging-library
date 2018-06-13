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

// Publish a message to a topic
func (c Connection) Publish(m client.Message, topic string) (err error) {
	body := []byte(m.Body)
	contentType := m.ContentType

	err = c.connection.Send(
		topic,
		contentType,
		body,
		stomp.SendOpt.Header("persistent", "true"),
	)

	return
}
