package cmd

import (
  "errors"
  "github.com/befovy/fvm/fvmgo"
  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(installCommand)
}

var installCommand = &cobra.Command{
  Use:   "install <version>",
  Short: "Installs Flutter SDK Version",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
      return errors.New("need to provide a channel or a version")
    }
    if len(args) > 1 {
      return errors.New("allows only one argument, the version or branch to install")
    }
    return nil
  },
  Run: func(cmd *cobra.Command, args []string) {
    fvmgo.CheckIfGitExists()
    version := args[0]
    if fvmgo.IsValidFlutterChannel(version) {
      fvmgo.FlutterChannelClone(version)
    } else {
      fvmgo.Verbosef("%s is not a valid Flutter channel, presume it's a Flutter version", version)
      fvmgo.FlutterVersionClone(version)
    }
  },
}
