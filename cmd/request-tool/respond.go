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
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/container-mgmt/messaging-library/pkg/client"
	"github.com/container-mgmt/messaging-library/pkg/connections/stomp"
)

var respondCmd = &cobra.Command{
	Use:   "respond",
	Short: "respond to messages from a queue",
	Long:  "Respond to messages from a queue.",
	Run:   runRespond,
}

func requestHandler(request client.Message) (response client.Message, err error) {
	if request.Err != nil {
		err = request.Err
		glog.Errorf(
			"Received error from queue: %s",
			err.Error(),
		)
		return
	}

	glog.Infof(
		"Received request from destination:\n%v",
		request.Data,
	)

	// create an echo response
	response = client.Message{
		Data: request.Data,
	}

	return
}

func runRespond(cmd *cobra.Command, args []string) {
	var c client.Connection
	var err error

	// Check mandatory arguments:
	if requestsQueue == "" {
		glog.Errorf("The argument 'requests-queue' is mandatory")
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

	// Create Responder
	r, err := c.NewResponder(
		client.ResponderSpec{
			RequestsQueue: requestsQueue,
			Callback:      requestHandler,
		})

	if err != nil {
		glog.Errorf(
			"Failed to create responder on '%s': %s",
			requestsQueue,
			err.Error(),
		)
		return
	}

	defer r.Close()

	glog.Infof(
		"Created responder on  '%s' (Press Ctrl-C to exit)",
		requestsQueue,
	)

	// wait for requests
	done := make(chan bool, 1)
	<-done
	return
}
