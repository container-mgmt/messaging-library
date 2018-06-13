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
	// In this example we will use the stomp backand,
	// last connection imported will be the default connection for the client.
	// a client can overide default connection by using the Use method.
	_ "github.com/container-mgmt/messaging-library/pkg/connections/stomp"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive messages from a destination",
	Long:  "Receive messages from a destination.",
	Run:   runReceive,
}

func callback(m client.Message, topic string) (err error) {
	glog.Infof(
		"Received message from destination '%s':\n%s",
		topic,
		m.Body,
	)

	return
}

func runReceive(cmd *cobra.Command, args []string) {
	var cli client.Client
	var con client.Connection

	// Check mandatory arguments:
	if destinationName == "" {
		glog.Errorf("The argument 'destination' is mandatory")
		return
	}

	// Build a new connection.
	con = cli.NewConnectionBuilder().
		BrokerHost(brokerHost).
		BrokerPort(brokerPort).
		UserName(userName).
		UserPassword(userPassword).
		UseTLS(useTLS).
		InsecureTLS(insecureTLS).
		Build()

	// Connect to the messaging service:
	err := con.Open()
	if err != nil {
		glog.Errorf(
			"Can't connect to message broker at host '%s' and port %d: %s",
			brokerHost,
			brokerPort,
			err.Error(),
		)
		return
	}
	defer con.Close()
	glog.Errorf(
		"Connected to message broker at host '%s' and port %d",
		brokerHost,
		brokerPort,
	)

	// Receive messages:
	err = con.Subscribe(destinationName, callback)
	if err != nil {
		glog.Errorf(
			"Can't subscribe to destination '%s': %s",
			destinationName,
			err.Error(),
		)
		return
	}
	glog.Infof(
		"Subscribed to destination '%s'",
		destinationName,
	)
	return
}
