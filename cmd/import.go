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

package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/befovy/fvm/fvmgo"
	"github.com/spf13/cobra"
)

var cop bool
var pos string

func init() {
	importCommand.Flags().BoolVarP(&cop, "copy", "c", false, "copy files instead of move")
	importCommand.Flags().StringVar(&pos, "path", "", "special the installed flutter path to import")
	rootCmd.AddCommand(importCommand)
}

var importCommand = &cobra.Command{
	Use:   "import <name>",
	Short: "Import installed flutter into fvm",
	Long: "Import installed flutter into fvm.\n" +
		"If there are more than one flutter detected,\n" +
		"this sub command use the first one in path by default.\n" +
		"Or you can use flags --path to special the path of flutter",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("need to provide a channel or a version or other name as import name, you can use `master` `beta` or `alibaba` `baidu` etc")
		}
		if len(args) > 1 {
			return errors.New("allows only one argument, the name to be imported as")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		dst := args[0]
		choices := fvmgo.FlutterListInstalledSdks()
		if len(choices) > 0 {
			for _, choice := range choices {
				if dst == choice {
					ins := fvmgo.YellowV("%s", dst)
					fvmgo.Errorf("fvm has already installed flutter with channel or version: %v", ins)
					return
				}
			}
		}

		flutters := fvmgo.FlutterOutOfFvm(pos)
		var source string
		if flutters == nil || len(flutters) == 0 {
			fvmgo.Warnf("fvm did not detect any flutter outside of fvm")
		} else if len(flutters) > 1 {
			if len(pos) == 0 {
				fvmgo.Infof("There are %d flutters outside of fvm, you don't set flag --path, so fvm will import te first one", len(flutters))
				source = flutters[0]
			} else {
				var match string
				for _, f := range flutters {
					if strings.HasPrefix(f, pos) {
						fvmgo.Infof("match %s", f)
						if len(match) == 0 {
							match = f
							source = match
						} else if match != f {
							fvmgo.Warnf("More than one installed flutter match the path: %v. You should provide a more detailed path", pos)
							source = ""
							break
						}
					}
				}
			}
		} else {
			source = flutters[0]
		}

		if len(source) > 0 {
			dst = filepath.Join(fvmgo.VersionsDir(), dst)
			source = filepath.Dir(filepath.Dir(source))
			fvmgo.Infof("%s will be imported into fvm", fvmgo.YellowV("%s", source))
			if cop {
				fvmgo.Infof("Copy all files from %v to %v, please wait seconds", source, dst)
			} else {
				fvmgo.Infof("Move all files from %v to %v, please wait seconds", source, dst)
			}
			err := fvmgo.CopyDir(source, dst)
			if err != nil {
				fvmgo.Errorf("Failed to copy files, %v", err)
			}
			if !cop {
				err = os.RemoveAll(source)
				if err != nil {
					fvmgo.Errorf("Failed to delete files, %v", err)
				}
			}
		}
	},
}
