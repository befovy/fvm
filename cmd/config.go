package cmd

import (
  "errors"
  "github.com/befovy/fvm/internal/config"
  "github.com/befovy/fvm/internal/log"
  "github.com/spf13/cobra"
)

var (
  cfg_cache_path string
  cfg_list       bool
)

func init() {
  configCommand.Flags().BoolVar(&cfg_list, "ls", false, "Lists all config options")
  configCommand.Flags().StringVarP(&cfg_cache_path, "cache-path", "c", "", "Path to store Flutter cached versions")
  rootCmd.AddCommand(configCommand)
}

var configCommand = &cobra.Command{
  Use:   "config",
  Short: "Config fvm options",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 0 {
      return errors.New("dose not take argument")
    }
    return nil
  },

  Run: func(cmd *cobra.Command, args []string) {
    if len(cfg_cache_path) > 0 {
      config.SetFlutterStoragePath(cfg_cache_path)
    }
    if cfg_list {
      all := config.AllConfig()
      if len(all) > 0 {
        log.Infof(all)
      } else {
        log.Warnf("No configuration has been set")
      }
    }
  },
}
