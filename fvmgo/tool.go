package fvmgo

import (
  "bytes"
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

func ProcessRunner(cmd string, dir string, arg ...string) {
  runner := exec.Command(cmd, arg...)
  if len(dir) == 0 {
    cwd, err := os.Getwd()
    if err != nil {
      Errorf("Cannot get work directory: %v", err)
      os.Exit(1)
    }
    runner.Dir = cwd
  } else {
    runner.Dir = dir
  }

  runner.Stderr = os.Stderr
  runner.Stdout = os.Stdout

  err := runner.Run()
  if err != nil {
    Errorf("Command '%s' exited with error: %v", cmd, err)
    os.Exit(runner.ProcessState.ExitCode())
  }
}

/// Returns true if it's a valid Flutter channel
func IsValidFlutterChannel(channel string) bool {
  channels := FlutterChannels()
  return stringSliceContains(channels, channel)
}

/// Returns true if it's a valid Flutter channel
func IsValidFlutterVersion(version string) bool {
  versions := FlutterListAllSdks()
  return stringSliceContains(versions, version)
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
    return path.Join(FvmHome(), "fvmbin", "flutter")
  }
}

func CurrentVersion() string {
  link := FlutterBin()
  if IsNotFound(link) {
    return ""
  }
  dst, err := os.Readlink(link)
  if err != nil {
    Errorf("Cannot read link target: %v", err)
    os.Exit(1)
  }
  return path.Base(path.Dir(path.Dir(dst)))
}

func IsCurrentVersion(version string) bool {
  return CurrentVersion() == version
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

func FlutterChannelClone(channel string) {
  if !IsValidFlutterChannel(channel) {
    Errorf("%s is not a valid flutter channel", channel)
    os.Exit(1)
  }

  Verbosef("%s is a valid flutter channel", channel)
  if checkInstalledCorrectly(channel) {
    Warnf("Flutter channel %s is already installed", channel)
    return
  }
  channelDir := path.Join(VersionsDir(), channel)
  Verbosef("Installing Flutter sdk %s to cache directory %s", channel, channelDir)
  err := os.MkdirAll(channelDir, 0755)
  if err != nil {
    Errorf("Cannot create directory for channel %s: %v", channel, err)
    os.Exit(1)
  }
  ProcessRunner("git", channelDir, "clone", "-b", channel, FlutterRepo, ".")
  Infof("Successfully installed flutter channel %s", channel)
}

func FlutterVersionClone(version string) {
  if !IsValidFlutterVersion(version) {
    Errorf("%s is not a valid version", version)
    os.Exit(1)
  }
  Verbosef("%s is a valid flutter version", version)
  if checkInstalledCorrectly(version) {
    Warnf("Flutter version %s is already installed", version)
    return
  }

  versionDir := path.Join(VersionsDir(), version)
  Verbosef("Installing Flutter sdk %s to cache directory %s", version, versionDir)

  err := os.MkdirAll(versionDir, 0755)
  if err != nil {
    Errorf("Cannot creat directory for version %s: %v", version, err)
    os.Exit(1)
  }
  ProcessRunner("git", versionDir, "clone", "-b", version, FlutterRepo, ".")
  Infof("Successfully installed flutter channel %s", version)
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
func CheckIfGitExists() {
  runner := exec.Command("git", "--version")
  Verbosef("Running `git --version` to check if git is available")
  err := runner.Run()
  if err != nil {
    Errorf("You need git installed to run fvm. Go to https://git-scm.com/downloads")
    os.Exit(1)
  }
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
        tags = append(tags, version[1])
      }
    }
  }
  return tags
}

func FlutterListInstalledSdks() []string {
  if IsNotFound(VersionsDir()) || !IsDirectory(VersionsDir()) {
    return []string{}
  }

  fis, err := ioutil.ReadDir(VersionsDir())
  if err != nil {
    Errorf("Cannot list installed versions: %v", err)
    os.Exit(1)
  }

  versions := make([]string, 0, len(fis))
  for _, fi := range fis {
    v := fi.Name()
    if checkInstalledCorrectly(v) {
      versions = append(versions, fi.Name())
    }
  }
  return versions
}

func projectFlutterLink(dir string, depth int) string {
  if depth == 0 {
    return ""
  }
  var link string
  if len(dir) == 0 {
    dir = WorkingDir()
  }
  link = path.Join(dir, ".fvmbin", "flutter")

  if IsSymlink(link) {
    return link
  } else if path.Dir(link) == link {
    return ""
  }

  depth -= 1
  return projectFlutterLink(path.Dir(dir), depth)
}

func linkFlutter(linkDir, version string) {
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

func LinkGlobalFlutter(version string) {
  linkPath := path.Join(FvmHome(), "fvmbin")
  linkFlutter(linkPath, version)

  osPath := os.Getenv("PATH")

  var paths []string
  if runtime.GOOS == "windows" {
    paths = strings.Split(osPath, ";")
  } else {
    paths = strings.Split(osPath, ":")
  }

  if !stringSliceContains(paths, linkPath) {
    Infof("add %s to path to enable flutter use", linkPath)
    if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
      Infof("export PATH=\"%s:$PATH\"", linkPath)
    } else {
      Infof("Add %s to PATH", linkPath)
    }
  }
}

func LinkProjectFlutter(version string) {
  linkPath := path.Join(WorkingDir(), ".fvmbin")
  linkFlutter(linkPath, version)
}
