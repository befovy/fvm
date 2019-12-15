package fvmgo

import (
  "github.com/spf13/viper"
  "os"
  "path"
)

var fvmEnvInited = false

/// check if home is a valid fvm home directory
/// home will be created if not exist,
func initFvmHome(home string) {
  magicFile := path.Join(home, ".fvmhome")
  if IsNotFound(home) {
    err := os.MkdirAll(home, 0755)
    if err != nil {
      Errorf("Can't create fvm home directory %s: %v", home, err)
      os.Exit(1)
    }
  } else if IsDirectory(home) && !IsFileExists(magicFile) {
    Errorf("Invalid fvm home, magic file \".fvmhome\" not exist")
    os.Exit(1)
  } else if IsFileExists(home) || IsSymlink(home) {
    Errorf("Invalid fvm home, %s is not a directory", home)
    os.Exit(1)
  }
}

func confirmConfigFile(filename string) {
  if !IsFileExists(filename) {
    f, err := os.Create(filename)
    if err != nil {
      Errorf("Can't create the fvm config file: %v", err)
      os.Exit(1)
    }
    err = f.Close()
    if err != nil {
      Errorf("Can't close the fvm config file: %v", err)
      os.Exit(1)
    }
  } else if IsDirectory(filename) {
    Errorf("Invalid config file, %s is a directory")
    os.Exit(1)
  } else if IsSymlink(filename) {
    Errorf("Invalid config file, %s is a symlink")
    os.Exit(1)
  }
}

func initFvmEnv() {
  if fvmEnvInited {
    return
  }
  fvmEnvInited = true
  home := os.Getenv("FVM_HOME")
  if len(home) == 0 {
    cfgDir, err := os.UserConfigDir()
    if err != nil {
      Errorf("Cant't get user config dir: %v", err)
      os.Exit(1)
    }
    home = path.Join(cfgDir, "fvm")
  }
  initFvmHome(home)
  cfgFile := path.Join(home, "config.yaml")
  confirmConfigFile(cfgFile)
  viper.SetConfigFile(cfgFile)
  err := viper.ReadInConfig()
  if err != nil {
    Errorf("Cannot load fvm config file: %v", err)
    os.Exit(1)
  }
  viper.Set("FVM_HOME", home)
}

/*
func GetConfigValue(key string) string {
  initFvmEnv()
  return viper.GetString(key)
}

func SetConfigValue(key, value string) {
  initFvmEnv()
  viper.Set(key, value)
  err := viper.WriteConfig()
  if err != nil {
    log.Errorf("Cannot save fvm config file: %v", err)
    os.Exit(1)
  }
}
*/


func createDir(dir, name string) {
  if IsNotFound(dir) {
    err := os.MkdirAll(dir, 0755)
    if err != nil {
      Errorf("Can't create versions dir: %v", err)
      os.Exit(1)
    }
  } else if !IsDirectory(dir) {
    Errorf("Invalid %s path, %s is not a directory", name, dir)
    os.Exit(1)
  }
}



func FvmHome() string {
  initFvmEnv()
  return viper.GetString("FVM_HOME")
}


func VersionsDir() string {
  dir := path.Join(FvmHome(), "versions")
  createDir(dir, "versions")
  return dir
}

func TempDir() string {
  dir := path.Join(FvmHome(), "temp")
  createDir(dir, "temp")
  return dir
}

func WorkingDir() string {
  workingDirectory, err := os.Getwd()
  if err != nil {
    Errorf("Can't get working directory: %v", err)
    os.Exit(1)
  }
  return workingDirectory
}
