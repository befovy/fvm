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
      fvmgo.ProcessRunner(dst, fvmgo.WorkingDir(), args...)
    }
  },
}
