package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"
  "os/signal"

  "github.com/befovy/fvm/internal/log"
)

//var colors bool
var verbose bool

func init() {
  // rootCmd.PersistentFlags().BoolVar(&colors, "colors", true, "Add Colors to log")
  rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output")
}

func initFvm() {
  log.Colorize()
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt, os.Kill)
  go func() {
    for range c {
      fmt.Println("Killed")
      os.Exit(1)
    }
  }()
}

var rootCmd = &cobra.Command{
  Use:   "fvm",
  Short: "Flutter Version Management",
  Long:  "Flutter Version Management: A cli to manage Flutter SDK versions.",
}

// Execute executes the rootCmd
func Execute() {
  cobra.OnInitialize(initFvm)
  if err := rootCmd.Execute(); err != nil {
    log.Errorf("Command failed: %v", err)
    os.Exit(1)
  }
}
