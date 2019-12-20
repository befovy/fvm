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
  "github.com/befovy/fvm/fvmgo"
  "github.com/spf13/cobra"
  "os"
)

func init() {
  rootCmd.AddCommand(flutterCommand)
}

var flutterCommand = &cobra.Command{
  Use:                "flutter",
  Short:              "Proxies Flutter Commands",
  DisableFlagParsing: true,
  Run: func(cmd *cobra.Command, args []string) {
    link := fvmgo.FlutterBin()
    if len(link) == 0 || !fvmgo.IsSymlink(link) {
      fvmgo.Errorf("No enabled Flutter sdk found. Create with <use> command")
    } else {
      dst, err := os.Readlink(link)
      if err != nil {
        fvmgo.Errorf("Cannot read link target: %v", err)
        os.Exit(1)
      }
      err = fvmgo.ProcessRunner(dst, fvmgo.WorkingDir(), args...)
      if err != nil {
        fvmgo.Errorf("Error while run flutter: %v", err)
      }
    }
  },
}
