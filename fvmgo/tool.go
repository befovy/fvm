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
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const FlutterRepo = "https://github.com/flutter/flutter.git"

func FlutterChannels() []string {
	return []string{
		"master", "stable", "dev", "beta",
	}
}

func stringSliceContains(slice []string, target string) bool {
	contains := false
	for _, v := range slice {
		if v == target {
			contains = true
			break
		}
	}
	return contains
}

func ProcessRunner(cmd string, dir string, arg ...string) error {
	runner := exec.Command(cmd, arg...)
	if len(dir) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot get work directory: %v", err))
		}
		runner.Dir = cwd
	} else {
		runner.Dir = dir
	}

	runner.Stderr = os.Stderr
	runner.Stdout = os.Stdout

	err := runner.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("Command '%s' exited with error: %v", cmd, err))
	}
	return nil
}

/// Returns true if it's a valid Flutter channel
func IsValidFlutterChannel(channel string) bool {
	channels := FlutterChannels()
	return stringSliceContains(channels, channel)
}

/// Returns true if it's a valid Flutter channel
func IsValidFlutterVersion(version string) bool {
	initFvmEnv()
	versions := viper.GetStringSlice("FLUTTER_REMOTE_TAGS")
	if stringSliceContains(versions, version) {
		return true
	} else {
		versions = FlutterListAllSdks()
		return stringSliceContains(versions, version)
	}
}

func IsValidFlutterInstall(version string) bool {
	versions := FlutterListInstalledSdks()
	return stringSliceContains(versions, version)
}

func FlutterBin() string {
	projectBin := projectFlutterLink("", 50)
	if len(projectBin) > 0 {
		return projectBin
	} else {
		return path.Join(FvmHome(), "current")
	}
}

func CurrentVersion() (string, error) {
	link := FlutterBin()
	if IsNotFound(link) {
		return "", nil
	}
	dst, err := os.Readlink(link)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot read link target: %v", err))
	}
	return path.Base(path.Dir(path.Dir(dst))), nil
}

func IsCurrentVersion(version string) bool {
	c, _ := CurrentVersion()
	return c == version
}

func FlutterSdkRemove(version string) {
	versionDir := path.Join(VersionsDir(), version)
	if !IsNotFound(versionDir) {
		err := os.RemoveAll(versionDir)
		if err != nil {
			Errorf("Cannot remove flutter version %s: %v", version, err)
			os.Exit(1)
		}
	}
}

func checkInstalledCorrectly(version string) bool {
	versionDir := path.Join(VersionsDir(), version)
	gitDir := path.Join(versionDir, ".github")
	binDir := path.Join(versionDir, "bin")

	if IsNotFound(versionDir) {
		return false
	}
	if IsNotFound(gitDir) || IsNotFound(binDir) {
		Warnf("%s exists but was not setup correctly. Doing cleanup...", version)
		FlutterSdkRemove(version)
		return false
	}
	return true
}

func FlutterChannelClone(channel string) error {
	if !IsValidFlutterChannel(channel) {
		return errors.New(fmt.Sprintf("%s is not a valid flutter channel", channel))
	}

	Verbosef("%s is a valid flutter channel", channel)
	if checkInstalledCorrectly(channel) {
		Warnf("Flutter channel %s is already installed", channel)
		return nil
	}
	channelDir := path.Join(VersionsDir(), channel)
	Verbosef("Installing Flutter sdk %s to cache directory %s", channel, channelDir)
	err := os.MkdirAll(channelDir, 0755)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create directory for channel %s: %v", channel, err))
	}
	err = ProcessRunner("git", channelDir, "clone", "-b", channel, FlutterRepo, ".")
	if err != nil {
		return err
	}
	Infof("Successfully installed flutter channel %s", channel)
	return nil
}

func FlutterVersionClone(version string) error {
	if !IsValidFlutterVersion(version) {
		return errors.New(fmt.Sprintf("%s is not a valid version", version))
	}
	Verbosef("%s is a valid flutter version", version)
	if checkInstalledCorrectly(version) {
		Warnf("Flutter version %s is already installed", version)
		return nil
	}

	versionDir := path.Join(VersionsDir(), version)
	Verbosef("Installing Flutter sdk %s to cache directory %s", version, versionDir)

	err := os.MkdirAll(versionDir, 0755)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot creat directory for version %s: %v", version, err))
	}
	err = ProcessRunner("git", versionDir, "clone", "-b", version, FlutterRepo, ".")
	if err != nil {
		return err
	}
	Infof("Successfully installed flutter channel %s", version)
	return nil
}

func gitGetVersion(p string) string {

	runner := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	runner.Dir = p
	b := new(bytes.Buffer)
	runner.Stdout = b

	err := runner.Run()

	out := string(b.Bytes())
	if strings.TrimSpace(out) == "HEAD" {
		runner = exec.Command("git", "tag", "--points-at", "HEAD")
		runner.Dir = p
		b.Reset()
		runner.Stdout = b

		err = runner.Run()
		if err != nil {
			Errorf("Cannot get git repo version: %v", err)
		}
		out = string(b.Bytes())
	}
	return strings.TrimSpace(out)
}

func flutterSdkVersion(branch string) string {
	branchDir := path.Join(VersionsDir(), branch)
	if IsNotFound(branchDir) {
		Errorf("Could not get version from SDK that is not installed")
		os.Exit(1)
	}
	return gitGetVersion(branchDir)
}

// CheckIfGitExists checks if git command is available
func CheckIfGitExists() error {
	runner := exec.Command("git", "--version")
	Verbosef("Running `git --version` to check if git is available")
	err := runner.Run()
	if err != nil {
		return errors.New("You need git installed to run fvm. Go to https://git-scm.com/downloads")
	}
	return nil
}

func FlutterListAllSdks() []string {
	runner := exec.Command("git", "ls-remote", "--tags", FlutterRepo)
	var b bytes.Buffer
	runner.Stdout = &b

	err := runner.Run()
	if err != nil {
		Errorf("Cannot list remote tags: %v", err)
		os.Exit(1)
	}

	tags := make([]string, 0)
	var tag string
	for {
		tag, err = b.ReadString('\n')
		if io.EOF == err {
			break
		} else if err != nil {
			Errorf("Cannot get exec runner output tag: %v", err)
			os.Exit(1)
		} else {
			version := strings.Split(tag, "refs/tags/")
			if len(version) > 1 {
				Verbosef("list remote tag: %s", strings.TrimSpace(version[1]))
				tags = append(tags, strings.TrimSpace(version[1]))
			}
		}
	}

	viper.Set("FLUTTER_REMOTE_TAGS", tags)
	err = viper.WriteConfig()
	if err != nil {
		Errorf("Can't write remote tags to config cache: %v", err)
	}
	return tags
}

func FlutterListInstalledSdks() []string {
	dir := VersionsDir()
	if IsNotFound(dir) || !IsDirectory(dir) {
		Verbosef("version directory %s is not found or is not a directory", dir)
		return []string{}
	}

	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		Errorf("Cannot list installed versions: %v", err)
		return []string{}
	} else {
		versions := make([]string, 0, len(fis))
		for _, fi := range fis {
			v := fi.Name()
			if checkInstalledCorrectly(v) {
				versions = append(versions, fi.Name())
			}
		}
		return versions
	}
}

func projectFlutterLink(dir string, depth int) string {
	if depth == 0 {
		return ""
	}
	var link string
	if len(dir) == 0 {
		dir = WorkingDir()
	}
	link = path.Join(dir, ".fvmbin", "current")

	if IsSymlink(link) {
		return link
	} else if path.Dir(link) == link {
		return ""
	}

	depth -= 1
	return projectFlutterLink(path.Dir(dir), depth)
}

func linkFlutterBin(linkDir, version string) {
	if !IsDirectory(linkDir) && !IsNotFound(linkDir) {
		Errorf("The path fvm used to make link exists but is not a directory")
		os.Exit(1)
	}

	if IsNotFound(linkDir) {
		err := os.MkdirAll(linkDir, 0755)
		if err != nil {
			Errorf("Can't make directory %s: %v", linkDir, err)
			os.Exit(1)
		}
	}

	versionBin := path.Join(VersionsDir(), version, "bin", "flutter")
	destLink := path.Join(linkDir, "flutter")

	if !IsNotFound(destLink) {
		err := os.RemoveAll(destLink)
		if err != nil {
			Errorf("Cannot remove link file: %v", err)
			os.Exit(1)
		}
	}

	err := os.Symlink(versionBin, destLink)
	if err != nil {
		Errorf("Cannot link flutter to global: %v", err)
		os.Exit(1)
	}
}

func linkFlutterDir(linkDir, version string) {
	versionDir := path.Join(VersionsDir(), version)

	if !IsNotFound(linkDir) {
		err := os.RemoveAll(linkDir)
		if err != nil {
			Errorf("Cannot remove link file: %v", err)
			os.Exit(1)
		}
	}

	err := os.Symlink(versionDir, linkDir)
	if err != nil {
		Errorf("Cannot link flutter to dest %s: %v", versionDir, err)
		os.Exit(1)
	}
}

func envPaths() []string {
	osPath := os.Getenv("PATH")
	var paths []string
	if runtime.GOOS == "windows" {
		paths = strings.Split(osPath, ";")
	} else {
		paths = strings.Split(osPath, ":")
	}
	return paths
}

func hasFlutterBin(name string) (bool, string) {
	if IsDirectory(name) && IsDirectory(path.Join(name, "bin")) {
		name = path.Join(name, "bin", "flutter")
	} else if IsDirectory(name) {
		name = path.Join(name, "flutter")
	}
	if IsSymlink(name) {
		dst, err := os.Readlink(name)
		if err != nil {
			Errorf("Cannot read link target: %v", err)
		} else {
			name = dst
		}
	}
	return IsExecutable(name), name
}

func FlutterOutOfFvm(install string) []string {
	paths := envPaths()
	res := make([]string, 0)
	if len(install) > 0 {
		paths = append(paths, install)
	}
	for _, p := range paths {
		has, name := hasFlutterBin(p)
		if has && !strings.HasPrefix(p, FvmHome()) {
			res = append(res, name)
		}
	}
	return res
}

func LinkGlobalFlutter(version string) {
	linkPath := path.Join(FvmHome(), "fvmbin")
	linkFlutterBin(linkPath, version)

	currentPath := path.Join(FvmHome(), "current")
	linkFlutterDir(currentPath, version)
	paths := envPaths()

	if !stringSliceContains(paths, currentPath) {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			cmd := YellowV("    export PATH=\"%s:$PATH\"", currentPath)
			Infof("Add %s to path to make sure you can use flutter from terminal\n%v", currentPath, cmd)
		} else {
			Warnf("Add %s to path to make sure you can use flutter from terminal", currentPath)
		}
	} else {
		Infof("linkpath: %v", linkPath)
	}
}

func LinkProjectFlutter(version string) {
	linkPath := path.Join(WorkingDir(), ".fvmbin")
	linkFlutterBin(linkPath, version)

	currentPath := path.Join(WorkingDir(), ".fvmbin", "current")
	linkFlutterDir(currentPath, version)
}
