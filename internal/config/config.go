package config

import (
  "fmt"
  "github.com/befovy/fvm/internal/constants"
  "github.com/befovy/fvm/internal/fileutil"
  "github.com/befovy/fvm/internal/log"
  "github.com/spf13/viper"
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

func readConfig() {
  confirmConfigFile()
  viper.SetConfigFile(constants.ConfigFile())

  err := viper.ReadInConfig()
  if err != nil {
    log.Errorf("Cannot load fvm config file: %v", err)
    os.Exit(1)
  }
}

func GetValue(key string) string {
  readConfig()
  return viper.GetString(key)
}

func SetValue(key, value string) {
  readConfig()
  viper.Set(key, value)
  err := viper.WriteConfig()
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
  readConfig()
  ssb := new(strings.Builder)
  ssb.WriteString("\n")
  for _, k := range viper.AllKeys() {
    v := viper.Get(k)
    if vs, ok := v.(string); ok {
      ssb.WriteString(fmt.Sprintf("%-12s : %s\n", k, vs))
    }
  }
  return ssb.String()
}
