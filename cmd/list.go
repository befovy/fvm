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
	"errors"
	"fmt"

	"github.com/befovy/fvm/fvmgo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists installed Flutter SDK Version",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("dose not take argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		choices := fvmgo.FlutterListInstalledSdks()
		if len(choices) == 0 {
			fvmgo.Warnf("No Flutter SDKs have been installed yet.")
		} else {
			for _, c := range choices {
				if fvmgo.IsCurrentVersion(c) {
					c = fmt.Sprintf("%s (current)", c)
				}
				fvmgo.Infof(c)
			}
		}
		flutters := fvmgo.FlutterOutOfFvm("")
		if flutters != nil && len(flutters) > 0 {
			fvmgo.Errorf("You have installed flutter outside of fvm")
			for _, f := range flutters {
				fvmgo.Warnf("-->  %v", f)
			}
			ins := fvmgo.YellowV("fvm import")
			if len(flutters) == 1 {
				fvmgo.Errorf("To import this into fvm, use %v", ins)
			} else {
				fvmgo.Errorf("To import these into fvm, use %v", ins)
			}
		}
	},
}
