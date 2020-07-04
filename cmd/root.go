/*
Copyright © 2019 befovy <befovy@gmail.com>

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
package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/befovy/fvm/fvmgo"
	"github.com/spf13/cobra"
)

var gVerbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&gVerbose, "verbose", "v", false, "Print verbose output")
}

func initFvm() {
	fvmgo.LogColorize()
	if gVerbose {
		fvmgo.LogVerbose()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for range c {
			fmt.Println("Killed")
			os.Exit(1)
		}
	}()
}

var rootCmd = &cobra.Command{
	Use:     "fvm",
	Short:   "Flutter Version Management",
	Long:    "Flutter Version Management: A cli to manage Flutter SDK versions.",
	Version: "0.8.0",
}

// Execute executes the rootCmd
func Execute() {
	cobra.OnInitialize(initFvm)
	if err := rootCmd.Execute(); err != nil {
		fvmgo.Errorf("Command failed: %v", err)
		os.Exit(1)
	}
}
