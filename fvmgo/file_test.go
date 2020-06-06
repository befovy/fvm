/*
Copyright © 2020 befovy <befovy@gmail.com>

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
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFileExist(t *testing.T) {
	if runtime.GOOS == "darwin" {
		if !IsFileExists("/bin/bash") {
			t.Fail()
		}
		if IsFileExists("/bin/hello-world") {
			t.Fail()
		}
	}
}

func TestIsDirectory(t *testing.T) {
	if runtime.GOOS == "darwin" {
		if !IsDirectory("/bin") {
			t.Fail()
		}

		if IsDirectory("/bin/bash") {
			t.Fail()
		}
	}
}

func TestIsEmptyDir(t *testing.T) {
	err := os.Mkdir("TestIsEmptyDir", 0644)
	var empty bool
	var f *os.File
	if err != nil {
		empty, err = IsEmptyDir("TestIsEmptyDir")
		if !empty {
			t.Fail()
		}

		f, err = os.Create(filepath.Join("TestIsEmptyDir", "keep"))
		if err != nil {
			t.Fail()
		} else {
			_ = f.Close()
		}
		empty, err = IsEmptyDir("TestIsEmptyDir")
		if empty {
			t.Fail()
		}
	}
	_ = os.RemoveAll("TestIsEmptyDir")
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound("/usr/local/fvm/versions/.DS_Store/.github") {
		t.Fail()
	}
}
