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
	brokerHost   string
	brokerPort   int
	userName     string
	userPassword string

	// Main command:
	rootCmd = &cobra.Command{
		Use:  "messaging-server",
		Long: "Your friendly STOMP broker (YFSB).",
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
		"127.0.0.1",
		"Default hostname for listening for connections.",
	)
	flags.IntVar(
		&brokerPort,
		"port",
		61613,
		"Default port for listening for connections.",
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

	// Register the subcommands:
	rootCmd.AddCommand(serveCmd)
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
