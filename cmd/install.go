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
package cmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/befovy/fvm/fvmgo"
	"github.com/spf13/cobra"
)

var repo string

func init() {
	installCommand.Flags().StringVar(&repo, "repo", "", "install flutter from unoffical git repo")
	rootCmd.AddCommand(installCommand)
}

func maybeVersion(ver string) bool {
	ver = strings.TrimLeft(ver, "vv")
	splits := strings.Split(ver, ".")
	hasNum := false
	for _, s := range splits {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			hasNum = true
		}
	}
	return hasNum
}

var installCommand = &cobra.Command{
	Use:   "install <version>",
	Short: "Installs Flutter SDK Version",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("need to provide a channel or a version")
		}
		if len(args) > 1 {
			return errors.New("allows only one argument, the version or branch to install")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		flutters := fvmgo.FlutterOutOfFvm("")
		if flutters != nil && len(flutters) > 0 {
			fvmgo.Errorf("You have installed flutter outside of fvm")
			for _, f := range flutters {
				fvmgo.Warnf("-->  %v", f)
			}
			ins := fvmgo.YellowV("fvm import")
			if len(flutters) == 1 {
				fvmgo.Errorf("To import this into fvm, use %v", ins)
			} else {
				fvmgo.Errorf("To import these into fvm, use %v", ins)
			}
		}

		err := fvmgo.CheckIfGitExists()
		if err == nil {
			version := args[0]
			if len(repo) > 0 {
				fvmgo.Infof("Install flutter <%s> from repo %s", version, repo)
				err = fvmgo.FlutterRepoClone(version, repo)
			} else if fvmgo.IsValidFlutterChannel(version) {
				err = fvmgo.FlutterChannelClone(version)
			} else if !maybeVersion(version) {
				fvmgo.Errorf("It seems that you want install a Flutter channel but have a invalid channel")
				channels := fvmgo.YellowV(strings.Join(fvmgo.FlutterChannels(), " "))
				fvmgo.Errorf("Please use one of %v", channels)
			} else {
				fvmgo.Verbosef("%s is not a valid Flutter channel, presume it's a Flutter version", version)
				err = fvmgo.FlutterVersionClone(version)
			}
		}
		if err != nil {
			fvmgo.Errorf(err.Error())
		}
	},
}
