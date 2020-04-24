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
	"github.com/befovy/fvm/fvmgo"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(currentCmd)
}

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current Flutter SDK info",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("dose not take argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		current, err := fvmgo.CurrentVersion()
		if err != nil {
			fvmgo.Errorf(err.Error())
		} else if len(current) == 0 {
			ins := fvmgo.YellowV("fvm use <version>")
			fvmgo.Warnf("No active Flutter sdk, please run %v", ins)
		} else {
			ins := fvmgo.YellowV(current)
			fvmgo.Infof("Current active Flutter SDK is %v", ins)

			link := fvmgo.FlutterBin()
			ins = fvmgo.YellowV(link)
			fvmgo.Infof("And its link path is %v", ins)

			dst, err := os.Readlink(link)
			if err != nil {
				fvmgo.Errorf("Cannot read link target: %v", err)
			} else {
				ins = fvmgo.YellowV(dst)
				fvmgo.Infof("Actually path is %v", ins)
			}
		}
	},
}
