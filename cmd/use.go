package cmd

import (
  "errors"
  "github.com/befovy/fvm/internal/log"
  "github.com/befovy/fvm/internal/tool"
  "github.com/logrusorgru/aurora"
  "github.com/spf13/cobra"
)

func init() {
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
    isValidInstall := tool.IsValidFlutterInstall(version)
    if !isValidInstall {
      ins := log.Au().Colorize("fvm install <version>", aurora.YellowFg)
      log.Errorf("Flutter %s is not installed. Please run %v", version, ins)
    } else {
      log.Infof("Activating")
      tool.LinkProjectFlutterDir(version)
      log.Infof("%s is active", version)
      log.Infof("Activating finished")
    }
  },
}
