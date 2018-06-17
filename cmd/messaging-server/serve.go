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
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/go-stomp/stomp/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a messages broker server",
	Long:  "Run a messages broker server.",
	Run:   runServe,
}

type myAuthenticator struct {
	userName     string
	userPasscode string
}

// Authenticate based on the given login and passcode, either of which might be nil.
// Returns true if authentication is successful, false otherwise.
func (m myAuthenticator) Authenticate(login, passcode string) bool {
	if m.userName == "" {
		return true
	}
	return login == m.userName && passcode == m.userPasscode
}

func (m myAuthenticator) Message() string {
	if m.userName == "" {
		return "listen to all"
	}
	return fmt.Sprintf(
		"listen only to %s:%s",
		m.userName,
		m.userPasscode,
	)
}

func runServe(cmd *cobra.Command, args []string) {
	brokerAddress := fmt.Sprintf("%s:%d", brokerHost, brokerPort)
	brokerAthenticator := myAuthenticator{
		userName:     userName,
		userPasscode: userPassword}
	brokerServer := server.Server{
		Addr:          brokerAddress,      // TCP address to listen on, DefaultAddr if empty
		Authenticator: brokerAthenticator, // Authenticates login/passcodes. If nil no authentication is performed
	}

	// Indicate we are running.
	glog.Infof(
		"Your friendly STOMP broker (YFSB) server %s: [%s]",
		brokerAddress,
		brokerAthenticator.Message(),
	)

	// ListenAndServe listens on the TCP network address and handle requests on
	// the incoming connections.
	err := brokerServer.ListenAndServe()
	if err != nil {
		glog.Errorf(
			"Can't start server on '%s': %s",
			brokerAddress,
			err.Error(),
		)
	}
}
