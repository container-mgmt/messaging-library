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
	"os"
	"testing"
	"time"

	"github.com/container-mgmt/messaging-library/pkg/client"

	"github.com/go-stomp/stomp/server"
)

//
// Utility functions used for testing
//

// Declare a chaneles
var messagesRecieved chan int
var serverStarted chan bool

// ListenAndServe open a STOMP testing server on address 127.0.0.1:61613.
func ListenAndServe() {
	s := &server.Server{}
	serverStarted <- true

	err := s.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to open STOMP testing server %s", err.Error())
		os.Exit(1)
	}
}

// Callback for subscribe testing
func callback(message client.Message, destination string) (err error) {
	err = message.Err

	if val, ok := message.Data["number"]; ok {
		messagesRecieved <- int(val.(float64))
	}
	return
}

//
// Start a STOMP server before runnig the tests
//

func TestMain(m *testing.M) {
	// Create a chanel
	messagesRecieved = make(chan int)
	serverStarted = make(chan bool)

	// Start a testing server on 127.0.0.1:61613
	go ListenAndServe()

	// Wait for server to start
	<-serverStarted

	code := m.Run()

	// Close channels
	close(serverStarted)
	close(messagesRecieved)

	os.Exit(code)
}

//
// Tests
//

func TestNewConnection(t *testing.T) {
	// Set the clients variables before we can open it.
	_, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		t.Error("NewConnection exited with an error")
	}
}

func TestOpenAndClose(t *testing.T) {
	// Create a connection
	c, _ := NewConnection(&client.ConnectionSpec{})

	// Try to open and close connection to server
	err := c.Open()
	if err != nil {
		t.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()
}

func TestPublish(t *testing.T) {
	// Create and open a connection
	c, _ := NewConnection(&client.ConnectionSpec{})
	c.Open()
	defer c.Close()

	// Set a hello world message
	m := client.Message{
		Data: client.MessageData{
			"number": 42,
		},
	}

	// Publish our hello message to the "destination name" destination.
	err := c.Publish(m, "destination-name")
	if err != nil {
		t.Errorf("Fail to publish a message: %s", err.Error())
	}
}

func TestPublishSubscribe(t *testing.T) {
	// Create and open a connection
	c, _ := NewConnection(&client.ConnectionSpec{})
	c.Open()
	defer c.Close()

	// Set a hello world message
	m := client.Message{
		Data: client.MessageData{
			"number": 42,
		},
	}

	// Subscribe to the "destination name" destination.
	go c.Subscribe("destination-name", callback)

loop:
	// Send messages until callback is called
	for {
		select {
		case r := <-messagesRecieved:
			if r == 42 {
				break loop
			}
		case <-time.After(1 * time.Microsecond):
			c.Publish(m, "destination-name")
		}
	}
}

//
// Benchmarks
//

func BenchmarkOpenAndClose(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c, _ := NewConnection(&client.ConnectionSpec{})
		c.Open()
		c.Close()
	}
}

func BenchmarkPublishSubscribe(b *testing.B) {
	// Create and open a connection
	c, _ := NewConnection(&client.ConnectionSpec{})
	c.Open()
	defer c.Close()

	// Subscribe to the "destination name" destination.
	go c.Subscribe("destination-name", callback)

	// Set a hello world message
	m := client.Message{
		Data: client.MessageData{
			"number": 42,
		},
	}

	for n := 0; n < b.N; n++ {
		fmt.Printf("Loop %d\n", n)
	loop:
		// Send messages until callback is called
		for {
			select {
			case r := <-messagesRecieved:
				if r == 42 {
					break loop
				}
			case <-time.After(10 * time.Microsecond):
				c.Publish(m, "destination-name")
			}
		}
	}
}
