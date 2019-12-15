package cmd

import (
  "errors"
  "github.com/befovy/fvm/fvmgo"
  "github.com/logrusorgru/aurora"
  "github.com/spf13/cobra"
)

var local bool

func init() {
  useCommand.Flags().BoolVar(&local, "local", false, "use SDK locally")
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
      ins := fvmgo.Au().Colorize("fvm install <version>", aurora.YellowFg)
      fvmgo.Errorf("Flutter %s is not installed. Please run %v", version, ins)
    } else {
      fvmgo.Infof("Activating")
      if local {
        fvmgo.LinkProjectFlutter(version)
      } else {
        fvmgo.LinkGlobalFlutter(version)
      }
      fvmgo.Infof("%s is active", version)
      fvmgo.Infof("Activating finished")
    }
  },
}
