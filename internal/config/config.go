package config

import (
  "fmt"
  "github.com/befovy/fvm/internal/constants"
  "github.com/befovy/fvm/internal/fileutil"
  "github.com/befovy/fvm/internal/log"
  "gopkg.in/ini.v1"
  "os"
  "strings"
)

const configFlutterStoredKey = "cache_path";


func confirmConfigFile() {
  if !fileutil.IsDirectory(constants.FvmHome()) {
    err := os.MkdirAll(constants.FvmHome(), 0755)
    if err != nil {
      log.Errorf("Cannot create the fvm home directory: %v", err)
      os.Exit(1)
    }
  }
  if !fileutil.IsFileExists(constants.ConfigFile()) {
    f, err := os.Create(constants.ConfigFile())
    if err != nil {
      log.Errorf("Cannot create the fvm config file: %v", err)
      os.Exit(1)
    }
    err = f.Close()
    if err != nil {
      log.Errorf("Cannot close the fvm config file: %v", err)
      os.Exit(1)
    }
  }
}

func readConfig() *ini.File {
  confirmConfigFile()
  ctx, err := ini.Load(constants.ConfigFile())
  if err != nil {
    log.Errorf("Cannot load fvm config file: %v", err)
    os.Exit(1)
  }
  return ctx
}

func GetValue(key string) string {
  ctx := readConfig()
  return ctx.Section("").Key(key).String()
}

func SetValue(key, value string) {
  ctx := readConfig()
  ctx.Section("").Key(key).SetValue(value)
  err := ctx.SaveTo(constants.ConfigFile())
  if err != nil {
    log.Errorf("Cannot save fvm config file: %v", err)
    os.Exit(1)
  }
}

func RemoveConfig() {
  if fileutil.IsFileExists(constants.ConfigFile()) {
    err := os.Remove(constants.ConfigFile())
    if err != nil {
      log.Errorf("Cannot remove config file: %v", err)
      os.Exit(1)
    }
  }
}

func SetFlutterStoragePath(p string) {
  if fileutil.IsDirectory(p) {
    SetValue(configFlutterStoredKey, p)
  } else if fileutil.IsNotFound(p) {
    err := os.MkdirAll(p, 0755)
    if err != nil {
      log.Errorf("Cannot create the flutter storage directory: %v", err)
      os.Exit(1)
    }
    SetValue(configFlutterStoredKey, p)
  } else {
    log.Errorf("ExceptionErrorFlutterPath")
    os.Exit(1)
  }
}

func GetFlutterStoragePath() string {
  p := GetValue(configFlutterStoredKey)
  if len(p) == 0 {
    return p
  }
  SetFlutterStoragePath(p)
  return p
}

func AllConfig() string {
  ctx := readConfig()
  sec := ctx.Section("")

  ssb := new(strings.Builder)
  for _, k := range sec.Keys() {
    ssb.WriteString(fmt.Sprintf("%s:%s\n", k.Name(), k.String()))
  }
  return ssb.String()
}
