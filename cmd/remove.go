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
  Use: "remove <version>",
  Short:"Removes Flutter SDK Version",
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
