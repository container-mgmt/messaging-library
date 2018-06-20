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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/go-stomp/stomp/server"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

//
// Utility functions used for testing.
//

// ListenAndServe open a STOMP testing server on address 127.0.0.1:61613.
func ListenAndServe(serverStarted chan bool) {
	s := &server.Server{}
	serverStarted <- true

	err := s.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to open new STOMP testing server: %s\n", err.Error())
	}
}

// Callback for subscribe testing.
func callbackFactory(c chan float64) client.SubscriptionCallback {
	return func(message client.Message, destination string) error {
		if val, ok := message.Data["value"]; ok {
			c <- val.(float64)
		}
		return message.Err
	}
}

//
// Start a STOMP server before runnig the tests.
//

func TestMain(m *testing.M) {
	// Create a chanel
	serverStarted := make(chan bool)
	defer close(serverStarted)

	// Start a testing server on 127.0.0.1:61613.
	log.SetOutput(ioutil.Discard)
	go ListenAndServe(serverStarted)

	// Wait for server to start.
	<-serverStarted

	code := m.Run()
	os.Exit(code)
}

//
// Tests.
//

func TestNewConnection(t *testing.T) {
	// Set the clients variables before we can open it.
	_, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		t.Error("NewConnection exited with an error")
	}
}

func TestOpenAndClose(t *testing.T) {
	// Create a connection.
	c, _ := NewConnection(&client.ConnectionSpec{})

	// Try to open and close connection to server.
	err := c.Open()
	if err != nil {
		t.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()
}

func TestPublish(t *testing.T) {
	// Create and open a connection.
	c, _ := NewConnection(&client.ConnectionSpec{})
	c.Open()
	defer c.Close()

	// Set a hello world message.
	m := client.Message{
		Data: client.MessageData{
			"value": 42.0,
		},
	}

	// Publish our hello message to the "destination name" destination.
	err := c.Publish(m, "destination-name")
	if err != nil {
		t.Errorf("Fail to publish a message: %s", err.Error())
	}
}

func TestPublishSubscribe(t *testing.T) {
	messageRecieved := make(chan float64)

	// Create and open a connection.
	c, _ := NewConnection(&client.ConnectionSpec{})
	c.Open()
	defer c.Close()

	// Set a hello world message.
	m := client.Message{
		Data: client.MessageData{
			"value": 42.0,
		},
	}

	fmt.Printf("Publish: %f\n", m.Data["value"])

	// Subscribe to the "destination name" destination.
	c.Subscribe("destination-name", callbackFactory(messageRecieved))
	// NOTE: The testing server, does not answer to Unsubscribe
	// We should Unsubscribe if using artimisMQ.
	//
	// defer c.Unsubscribe("destination-name")

	c.Publish(m, "destination-name")

	r := <-messageRecieved
	if r != 42 {
		t.Errorf("Received %f expected 42", r)
	}

	fmt.Printf("Received: %f\n", r)
}

//
// Benchmarks.
//

func BenchmarkOpenAndClose(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c, _ := NewConnection(&client.ConnectionSpec{})
		c.Open()
		c.Close()
	}
}

func BenchmarkPublishAndSubscribe(b *testing.B) {
	messageRecieved := make(chan float64, b.N)
	defer close(messageRecieved)

	// Create and open a connection.
	c, _ := NewConnection(&client.ConnectionSpec{})
	c.Open()
	defer c.Close()

	// Subscribe to the "destination name" destination.
	c.Subscribe("destination-name", callbackFactory(messageRecieved))
	defer c.Unsubscribe("destination-name")

	// Set a hello world message.
	m := client.Message{
		Data: client.MessageData{
			"value": 42.0,
		},
	}

	// Publish b.N messages.
	for n := 0; n < b.N; n++ {
		c.Publish(m, "destination-name")
	}

	n := 0
loop:
	for r := range messageRecieved {
		// Check response.
		if r != 42 {
			b.Errorf("Received %f expected 42", r)
		}

		// Exit on the b.N message.
		if n++; n == b.N {
			break loop
		}
	}
}
