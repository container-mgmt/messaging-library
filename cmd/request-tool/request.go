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

package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/container-mgmt/messaging-library/pkg/client"
	"github.com/container-mgmt/messaging-library/pkg/connections/stomp"
)

var (
	contentType string
	messageBody string
)

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Sends a request on the requests queue",
	Long:  "Sends a request on the requests queue.",
	Run:   runRequest,
}

var waitResponse chan bool

func init() {
	flags := requestCmd.Flags()
	flags.StringVar(
		&contentType,
		"content-type",
		"application/json",
		"The MIME type of the message body.",
	)
	flags.StringVar(
		&messageBody,
		"body",
		"",
		"The body of the message. If it starts with the @ character then the rest will be "+
			"interpreted as a file name, and the body of the message will be the content "+
			"of that file. If this option isn't given then the body will be taken from "+
			"the standard input.",
	)
}

func responseHandler(response client.Message, requestID string) (err error) {
	if response.Err != nil {
		err = response.Err
		glog.Errorf(
			"Received error from queue: %s",
			err.Error(),
		)
		return
	}

	glog.Infof(
		"Received response to request id %s:\n%v",
		requestID,
		response.Data,
	)

	waitResponse <- true
	return
}

func runRequest(cmd *cobra.Command, args []string) {
	var c client.Connection
	var bodyBytes []byte
	var err error

	// Check mandatory arguments:
	if requestsQueue == "" {
		glog.Errorf("The argument 'requests-queue' is mandatory")
		return
	}

	if responsesQueue == "" {
		glog.Errorf("The argument 'responses-queue' is mandatory")
		return
	}

	// Set the clients variables before we can open it.
	c, err = stomp.NewConnection(&client.ConnectionSpec{
		// Global options:
		BrokerHost:   brokerHost,
		BrokerPort:   brokerPort,
		UserName:     userName,
		UserPassword: userPassword,
		UseTLS:       useTLS,
		InsecureTLS:  insecureTLS,
	})
	if err != nil {
		glog.Errorf(
			"Can't connect to message broker at host '%s' and port %d: %s",
			brokerHost,
			brokerPort,
			err.Error(),
		)
		return
	}
	defer c.Close()
	glog.Infof(
		"Connected to message broker at host '%s' and port %d",
		brokerHost,
		brokerPort,
	)

	// Create the Requestor
	r, err := c.NewRequestor(
		client.RequestorSpec{
			RequestsQueue:  requestsQueue,
			ResponsesQueue: responsesQueue,
		})

	if err != nil {
		glog.Errorf(
			"Failed to create requestor : %s",
			err.Error(),
		)
		return
	}

	defer r.Close()

	glog.Infof(
		"Created requestor to '%s'",
		requestsQueue,
	)

	// Load the message body:
	var body string
	if messageBody == "" {
		glog.Info("Please insert message body:")

		bodyBytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			glog.Errorf(
				"Can't read message body from standard input: %s",
				err.Error(),
			)
			return
		}
		body = string(bodyBytes)
	} else if messageBody[0] == '@' {
		messageFile := messageBody[1:]
		bodyBytes, err = ioutil.ReadFile(messageFile)
		if err != nil {
			glog.Errorf(
				"Can't read message body from file '%s': %s",
				messageFile,
				err.Error(),
			)
			return
		}
		body = string(bodyBytes)
	} else {
		body = messageBody
	}

	// Inform user about the message body.
	glog.Infof("Message: %s", body)

	// Send a message:
	// Create a message with data object to send to the message broker,
	// the data payload must be of type client.MessageData{}.
	m := client.Message{
		ContentType: contentType,
		Data: client.MessageData{
			"kind": "Request",
			"body": body,
		},
	}

	waitResponse = make(chan bool, 1)
	r.Send(m, responseHandler)

	// Wait for response.
	<-waitResponse
}
