/*
Copyright © 2019 befovy <befovy@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package fvmgo

import (
  "io/ioutil"
  "os"
)

// IsFileExists checks if a file exists and is not a directory
func IsFileExists(name string) bool {
  info, err := os.Stat(name)
  if err != nil {
    return false
  }
  return !info.IsDir()
}

// IsDirectory check if path exists and is a directory
func IsDirectory(name string) bool {
  info, err := os.Stat(name)
  if err != nil {
    return false
  }
  return info.IsDir()
}

func IsEmptyDir(name string) (bool, error) {
  entries, err := ioutil.ReadDir(name)
  if err != nil {
    return false, err
  }
  return len(entries) == 0, nil
}

func IsSymlink(name string) bool {
  info, err := os.Lstat(name)
  if os.IsNotExist(err) {
    return false
  } else if err != nil {
    Warnf("Error when check symlink: %v", err)
    return false
  }
  return (info.Mode() & os.ModeSymlink) != 0
}

func IsNotFound(path string) bool {
  _, err := os.Lstat(path)
  return os.IsNotExist(err)
}
