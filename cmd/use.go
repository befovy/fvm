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

var local bool

func init() {
  useCommand.Flags().BoolVarP(&local, "local", "l", false, "use SDK locally")
  rootCmd.AddCommand(useCommand)
}

var useCommand = &cobra.Command{
  Use:   "use <version>",
  Short: "Which Flutter SDK Version you would like to use",
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
      ins := fvmgo.YellowV("fvm install <version>")
      fvmgo.Errorf("Flutter %s is not installed. Please run %v", version, ins)
    } else {
      if local {
        fvmgo.LinkProjectFlutter(version)
      } else {
        fvmgo.LinkGlobalFlutter(version)
      }
      fvmgo.Infof("%s is active", version)
    }
  },
}
