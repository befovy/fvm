package cmd

import (
  "errors"
  "fmt"
  "github.com/befovy/fvm/internal/log"
  "github.com/befovy/fvm/internal/tool"
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
    choices := tool.FlutterListInstalledSdks()
    if len(choices) == 0 {
      log.Warnf("No SDKs have been installed yet.")
    } else {
      for _, c := range choices {
        if tool.IsCurrentVersion(c) {
        c = fmt.Sprintf("%s (current)", c)
        }
        log.Infof(c)
      }
    }
  },
}
