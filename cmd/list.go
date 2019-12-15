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
  },
}
