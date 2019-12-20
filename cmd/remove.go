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
  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(removeCommand)
}

var removeCommand = &cobra.Command{
  Use:   "remove <version>",
  Short: "Removes Flutter SDK Version",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
      return errors.New("need to provide a channel or a version")
    }
    if len(args) > 1 {
      return errors.New("allows only one argument, the version or branch to remove")
    }
    return nil
  },
  Run: func(cmd *cobra.Command, args []string) {

    version := args[0]
    isValidInstall := fvmgo.IsValidFlutterInstall(version)
    if !isValidInstall {
      fvmgo.Warnf("Flutter SDK: %s is not installed", version)
    } else {
      fvmgo.Infof("Removing %s", version)
      fvmgo.FlutterSdkRemove(version)
      fvmgo.Infof("Removing %s finished", version)
    }
  },
}
