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
	"flag"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// Global options:
	brokerHost      string
	brokerPort      int
	destinationName string
	userName        string
	userPassword    string
	useTLS          bool
	insecureTLS     bool

	// Main command:
	rootCmd = &cobra.Command{
		Use:  "messaging-tool",
		Long: "A tool that can send and receive messages using a message broker.",
	}
)

func init() {
	// Send logs to the standard error stream by default:
	flag.Set("logtostderr", "true")

	// Register the options that are managed by the 'flag' package, so that they will also be parsed
	// by the 'pflag' package:
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Register the global options:
	flags := rootCmd.PersistentFlags()
	flags.StringVar(
		&brokerHost,
		"host",
		"localhost",
		"The IP address or port number of the message broker.",
	)
	flags.IntVar(
		&brokerPort,
		"port",
		61613,
		"The port number of the message server.",
	)
	flags.StringVar(
		&destinationName,
		"destination",
		"",
		"The name of the destination.",
	)
	flags.StringVar(
		&userName,
		"user",
		"",
		"The name of the user.",
	)
	flags.StringVar(
		&userPassword,
		"password",
		"",
		"The password of the user.",
	)
	flags.BoolVar(
		&useTLS,
		"tls",
		false,
		"Use TLS.",
	)
	flags.BoolVar(
		&insecureTLS,
		"insecure",
		false,
		"Don't check the server TLS certificate and host name.",
	)

	// Register the subcommands:
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(receiveCmd)
}

func main() {
	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been
	// parsed.
	flag.CommandLine.Parse([]string{})

	// Execute the root command:
	rootCmd.SetArgs(os.Args[1:])
	rootCmd.Execute()
}
