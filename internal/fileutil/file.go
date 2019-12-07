package fileutil

import (
  "github.com/befovy/fvm/internal/log"
  "os"
)

// IsFileExists checks if a file exists and is not a directory
func IsFileExists(name string) bool {
  info, err := os.Stat(name)
  if os.IsNotExist(err) {
    return false
  }
  return !info.IsDir()
}

// IsDirectory check if path exists and is a directory
func IsDirectory(name string) bool {
  info, err := os.Stat(name)
  if os.IsNotExist(err) {
    return false
  }
  return info.IsDir()
}

func IsSymlink(name string) bool {
  info, err := os.Lstat(name)
  if os.IsNotExist(err) {
    return false
  } else if err != nil {
    log.Warnf("Error when check symlink: %v", err)
    return false
  }
  return (info.Mode() & os.ModeSymlink) != 0
}

func IsNotFound(path string) bool {
  _, err := os.Stat(path)
  return os.IsNotExist(err)
}
