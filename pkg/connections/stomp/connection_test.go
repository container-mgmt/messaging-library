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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-stomp/stomp/server"

	"github.com/container-mgmt/messaging-library/pkg/client"
)

//
// Utility functions used for testing.
//

// We are using the internal server to run the tests.
var UseInternalServer = true

// DestinationName generate a random destination name.
func DestinationName() (string, error) {
	data := make([]byte, 10)
	_, err := rand.Read(data)

	return base64.StdEncoding.EncodeToString(data), err
}

// ListenAndServe open a STOMP testing server on address 127.0.0.1:61613.
func ListenAndServe(serverStarted chan bool) {
	s := &server.Server{}
	serverStarted <- true

	fmt.Println("Starting a new STOMP testing server.")

	err := s.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to open new STOMP testing server: %s\n", err.Error())

		// We assume that the reason we can't open the internal server is because
		// we have an extrenal one running.
		fmt.Println("Continue, falling back to external server.")
		UseInternalServer = false
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
	c, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		t.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()
}

func TestPublish(t *testing.T) {
	// Get a unique destination for the test.
	destination, _ := DestinationName()

	// Create and open a connection.
	c, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		t.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()

	// Set a hello world message.
	m := client.Message{
		Data: client.MessageData{
			"value": 42.0,
		},
	}

	// Publish our hello message to the "destination name" destination.
	err = c.Publish(m, destination)
	if err != nil {
		t.Errorf("Fail to publish a message: %s", err.Error())
	}
}

func TestPublishSubscribe(t *testing.T) {
	contentType := "application/json"

	// Get a unique destination for the test.
	destination, _ := DestinationName()
	messageRecieved := make(chan float64)
	// [ When using artimisMQ we can close ]  defer close(messageRecieved)

	// Create and open a connection.
	c, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		t.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()

	// Set a hello world message.
	m := client.Message{
		ContentType: contentType,
		Data: client.MessageData{
			"value": 42.0,
		},
	}

	// Subscribe to the "destination name" destination.
	c.Subscribe(destination, callbackFactory(messageRecieved))
	// [ When using artimisMQ we can Unsubscribe ] defer c.Unsubscribe(destination)

	c.Publish(m, destination)

	r := <-messageRecieved
	if r != 42 {
		t.Errorf("Received %f expected 42", r)
	}
}

func TestPublishSubscribeRunTime(t *testing.T) {
	var m client.Message

	// This test can only run on external server.
	if UseInternalServer {
		t.Skip("skipping test when running using internal server.")
	}

	// Run N times
	runTimes := 1000
	runMaxTime, _ := time.ParseDuration("1.5s")

	// Get a unique destination for the test
	destination, _ := DestinationName()
	messageRecieved := make(chan float64)
	// defer close(messageRecieved)

	// Create a 1k payload to send on the queue.
	thousandBytes := make([]byte, 1000)

	// Create and open a connection.
	c, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		t.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()

	// Subscribe to the "destination name" destination.
	c.Subscribe(destination, callbackFactory(messageRecieved))
	// defer c.Unsubscribe(destination)

	// Start timer
	start := time.Now()

	// Set a 1Kb message 1000 times.
	for n := 0; n < runTimes; n++ {
		m = client.Message{
			Data: client.MessageData{
				"value":        n,
				"bigByteArray": thousandBytes,
			},
		}
		c.Publish(m, destination)
	}

	for n := 0; n < runTimes; n++ {
		r := <-messageRecieved
		if int(r) != n {
			t.Errorf("Received %f expected %d", r, n)
		}
	}

	// Check time elapsed
	elapsed := time.Since(start)
	if elapsed > runMaxTime {
		t.Errorf("Test took %s [max is %s]", elapsed, runMaxTime)
	}
}

//
// Benchmarks.
//

func BenchmarkOpenAndClose(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c, err := NewConnection(&client.ConnectionSpec{})
		if err != nil {
			b.Errorf("Fail to open connection: %s", err.Error())
		}
		c.Close()
	}
}

func BenchmarkPublishAndSubscribe1Kb(b *testing.B) {
	var m client.Message

	// Create a 1k payload to send on the queue.
	thousandBytes := make([]byte, 1000)

	// Get a unique destination for the test.
	destination, _ := DestinationName()
	messageRecieved := make(chan float64, b.N)
	defer close(messageRecieved)

	// Create and open a connection.
	c, err := NewConnection(&client.ConnectionSpec{})
	if err != nil {
		b.Errorf("Fail to open connection: %s", err.Error())
	}
	defer c.Close()

	// Subscribe to the "destination name" destination.
	c.Subscribe(destination, callbackFactory(messageRecieved))
	defer c.Unsubscribe(destination)

	// Publish b.N messages.
	for n := 0; n < b.N; n++ {
		// Send a message with counter as value.
		m = client.Message{
			Data: client.MessageData{
				"value":        n,
				"bigByteArray": thousandBytes,
			},
		}

		c.Publish(m, destination)
	}

	n := 0
loop:
	for r := range messageRecieved {
		// Check response.
		if int(r) != n {
			b.Errorf("Received %f expected %d", r, n)
		}

		// Exit on the b.N message.
		if n++; n == b.N {
			break loop
		}
	}
}
