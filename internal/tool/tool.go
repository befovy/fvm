package tool

import (
  "bytes"
  "github.com/befovy/fvm/internal/config"
  "github.com/befovy/fvm/internal/constants"
  "github.com/befovy/fvm/internal/fileutil"
  "github.com/befovy/fvm/internal/log"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "path"
  "strings"
)

func ProcessRunner(cmd string, dir string, arg ...string) {
  runner := exec.Command(cmd, arg...)
  if len(dir) == 0 {
    cwd, err := os.Getwd()
    if err != nil {
      log.Errorf("Cannot get work directory: %v", err)
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
    log.Errorf("Command '%s' exited with error: %v", cmd, err)
    os.Exit(runner.ProcessState.ExitCode())
  }
}

/// Returns true if it's a valid Flutter channel
func IsValidFlutterChannel(channel string) bool {
  channels := constants.FlutterChannels()
  valid := false
  for _, v := range channels {
    if v == channel {
      valid = true
      break
    }
  }
  return valid
}

func IsValidFlutterVersion(version string) bool {
  versions := FlutterListAllSdks()

  valid := false
  for _, v := range versions {
    v = strings.TrimSpace(v)
    if version == v {
      valid = true
      break
    }
  }
  return valid
}

func IsValidFlutterInstall(version string) bool {
  versions := FlutterListInstalledSdks()
  installed := false
  for _, v := range versions {
    if v == version {
      installed = true
      break
    }
  }
  return installed
}

func projectFlutterLink(dir string, depth int) string {
  if depth == 0 {
    return ""
  }
  var link string
  if len(dir) == 0 {
    dir = constants.WorkingDirectory()
  }
  link = path.Join(dir, "fvm")

  if fileutil.IsSymlink(link) {
    return link
  } else if path.Dir(link) == link {
    return ""
  }

  depth -= 1
  return projectFlutterLink(path.Dir(dir), depth)
}

func ProjectFlutterLink() string {
  return projectFlutterLink(constants.WorkingDirectory(), 20)
}

func IsCurrentVersion(version string) bool {
  link := ProjectFlutterLink()
  if fileutil.IsNotFound(link) {
    return false
  }
  dst, err := os.Readlink(link)
  if err != nil {
    log.Errorf("Cannot read link target: %v", err)
    os.Exit(1)
  }
  current := path.Base(path.Dir(path.Dir(dst)))
  return current == version
}

func versionsDir() string {
  flutterPath := config.GetFlutterStoragePath()
  if len(flutterPath) != 0 {
    return flutterPath
  }

  return path.Join(constants.FvmHome(), "versions")
}

func FlutterSdkRemove(version string) {
  versionDir := path.Join(versionsDir(), version)
  if !fileutil.IsNotFound(versionDir) {
    err := os.RemoveAll(versionDir)
    if err != nil {
      log.Errorf("Cannot remove flutter version %s: %v", version, err)
      os.Exit(1)
    }
  }
}

func checkInstalledCorrectly(version string) bool {
  versionDir := path.Join(versionsDir(), version)
  gitDir := path.Join(versionDir, ".github")
  binDir := path.Join(versionDir, "bin")

  if fileutil.IsNotFound(versionDir) {
    return false
  }
  if fileutil.IsNotFound(gitDir) || fileutil.IsNotFound(binDir) {
    log.Warnf("%s exists but was not setup correctly. Doing cleanup...", version)
    FlutterSdkRemove(version)
    return false
  }
  return true
}

func FlutterChannelClone(channel string) {
  if !IsValidFlutterChannel(channel) {
    log.Errorf("%s is not a invalid channel", channel)
    os.Exit(1)
  }

  if checkInstalledCorrectly(channel) {
    return
  }
  channelDir := path.Join(versionsDir(), channel)
  err := os.MkdirAll(channelDir, 0755)
  if err != nil {
    log.Errorf("Cannot creat directory for channel %s: %v", channel, err)
    os.Exit(1)
  }
  ProcessRunner("git", channelDir, "clone", "-b", channel, constants.FlutterRepo, ".")
}

func FlutterVersionClone(version string) {
  if !IsValidFlutterVersion(version) {
    log.Errorf("%s is not a valid version", version)
    os.Exit(1)
  }

  if checkInstalledCorrectly(version) {
    return
  }

  versionDir := path.Join(versionsDir(), version)
  err := os.MkdirAll(versionDir, 0755)
  if err != nil {
    log.Errorf("Cannot creat directory for version %s: %v", version, err)
    os.Exit(1)
  }
  ProcessRunner("git", versionDir, "clone", "-b", version, constants.FlutterRepo, ".")
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
      log.Errorf("Cannot get git repo version: %v", err)
    }
    out = string(b.Bytes())
  }
  return strings.TrimSpace(out)
}

func flutterSdkVersion(branch string) string {
  branchDir := path.Join(versionsDir(), branch)
  if fileutil.IsNotFound(branchDir) {
    log.Errorf("Could not get version from SDK that is not installed")
    os.Exit(1)
  }
  return gitGetVersion(branchDir)
}

func CheckIfGitExists() {

  runner := exec.Command("git", "--version")
  err := runner.Run()
  if err != nil {
    log.Errorf("You need Git Installed to run fvm. Go to https://git-scm.com/downloads")
    os.Exit(1)
  }
}

func FlutterListAllSdks() []string {
  //constants.FlutterRepo

  runner := exec.Command("git", "ls-remote", "--tags", constants.FlutterRepo)

  var b bytes.Buffer
  runner.Stdout = &b

  err := runner.Run()
  if err != nil {
    log.Errorf("Cannot list remote tags: %v", err)
    os.Exit(1)
  }

  tags := make([]string, 0)
  var tag string
  for {
    tag, err = b.ReadString('\n')
    if io.EOF == err {
      break
    } else if err != nil {
      log.Errorf("Cannot get exec runner output tag: %v", err)
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
  if fileutil.IsNotFound(versionsDir()) || !fileutil.IsDirectory(versionsDir()) {
    return []string{}
  }

  fis, err := ioutil.ReadDir(versionsDir())
  if err != nil {
    log.Errorf("Cannot list installed versions: %v", err)
    os.Exit(1)
  }

  versions := make([]string, 0, len(fis))
  for _, fi := range fis {
    versions = append(versions, fi.Name())
  }
  return versions
}

func LinkProjectFlutterDir(version string) {
  versionBin := path.Join(versionsDir(), version, "bin", "flutter")
  localLink := constants.LocalFlutterLink()
  if !fileutil.IsNotFound(localLink) {
    err := os.RemoveAll(localLink)
    if err != nil {
      log.Errorf("Cannot remove local link file: %v", err)
      os.Exit(1)
    }
  }

  err := os.Symlink(versionBin, localLink)
  if err != nil {
    log.Errorf("Cannot link flutter to local: %v", err)
    os.Exit(1)
  }
}
