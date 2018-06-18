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
	contentType  string
	messageBody  string
	messageCount int
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends a message to a destination",
	Long:  "Sends a message to a destination.",
	Run:   runSend,
}

type dataPayload struct {
	kind string
	spec map[string]string
}

func init() {
	flags := sendCmd.Flags()
	flags.StringVar(
		&contentType,
		"content-type",
		"text/plain",
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
	flags.IntVar(
		&messageCount,
		"count",
		1,
		"The number of messages to send.",
	)
}

func runSend(cmd *cobra.Command, args []string) {
	var c client.Connection
	var bodyBytes []byte
	var err error

	// Check mandatory arguments:
	if destinationName == "" {
		glog.Errorf("The argument 'destination' is mandatory")
		return
	}

	// Set the clients variables before we can open it.
	c, err = stomp.NewConnection(&stomp.ConnectionBuilder{
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
			"Can't create a new connection to host '%s': %s",
			brokerHost,
			err.Error(),
		)
		return
	}

	// Connect to the messaging service:
	err = c.Open()
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

	// Send a message:
	for i := 0; i < messageCount; i++ {
		data := dataPayload{
			kind: "InfoMessage",
			spec: map[string]string{"message": body},
		}

		m := client.Message{
			Data: data,
		}

		glog.Info(body)
		err = c.Publish(m, destinationName)
		if err != nil {
			glog.Errorf(
				"Can't send message to destination '%s': %s",
				destinationName,
				err.Error(),
			)
			break
		}
		if messageCount > 1 {
			glog.Infof("Message %d sent", i)
		} else {
			glog.Infof("Message sent")
		}
	}
}
