package constants

import (
  "github.com/befovy/fvm/internal/log"
  "os"
  "path"
)

const (
  FlutterRepo = "https://github.com/flutter/flutter.git"
)

var workingDirectory string

func init() {
  workingDirectory, _ = os.Getwd()
}

func FlutterChannels() []string {
  return []string{
    "master", "stable", "dev", "beta",
  }
}

func WorkingDirectory() string {
  return workingDirectory
}

func LocalFlutterLink() string {
  return path.Join(workingDirectory, "fvm")
}

// FVM Home directory
func FvmHome() string {
  home, err := os.UserConfigDir()
  if err != nil {
    log.Errorf("Cannot get home dir: %v", err)
  }
  return path.Join(home, "fvm")
}

func ConfigFile() string {
  return path.Join(FvmHome(), "config.ini")
}
