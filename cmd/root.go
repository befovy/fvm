package cmd

import (
  "fmt"
  "github.com/befovy/fvm/fvmgo"
  "github.com/spf13/cobra"
  "os"
  "os/signal"
)

var gVerbose bool

func init() {
  rootCmd.PersistentFlags().BoolVarP(&gVerbose, "verbose", "v", false, "Print verbose output")
}

func initFvm() {
  fvmgo.LogColorize()
  if gVerbose {
    fvmgo.LogVerbose()
  }

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
  Use:     "fvm",
  Short:   "Flutter Version Management",
  Long:    "Flutter Version Management: A cli to manage Flutter SDK versions.",
  Version: "0.2.2",
}

// Execute executes the rootCmd
func Execute() {
  cobra.OnInitialize(initFvm)
  if err := rootCmd.Execute(); err != nil {
    fvmgo.Errorf("Command failed: %v", err)
    os.Exit(1)
  }
}
