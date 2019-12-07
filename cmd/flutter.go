package cmd

import (
  "github.com/befovy/fvm/internal/constants"
  "github.com/befovy/fvm/internal/fileutil"
  "github.com/befovy/fvm/internal/log"
  "github.com/befovy/fvm/internal/tool"
  "github.com/spf13/cobra"
  "os"
)

func init() {
  rootCmd.AddCommand(flutterCommand)
}

var flutterCommand = &cobra.Command{
  Use:   "flutter",
  Short: "Proxies Flutter Commands",
  DisableFlagParsing:true,
  Run: func(cmd *cobra.Command, args []string) {
    link := tool.ProjectFlutterLink()
    if len(link) == 0 || !fileutil.IsSymlink(link) {
      log.Errorf("No FVM config found. Create with <use> command")
    } else {

      dst, err := os.Readlink(link)
      if err != nil {
        log.Errorf("Cannot read link target: %v", err)
        os.Exit(1)
      }
      tool.ProcessRunner(dst, constants.WorkingDirectory(), args...)
    }
  },
}
