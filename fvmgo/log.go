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
	"fmt"

	"github.com/logrusorgru/aurora"
)

// Colorize chnage the logger to support colors printing.
func LogColorize() {
	logAu = aurora.NewAurora(true)
}

func LogVerbose() {
	logVerbose = true
}

// internal colorized
var logAu aurora.Aurora
var logVerbose bool

// au Aurora instance used for colors
func au() aurora.Aurora {
	if logAu == nil {
		logAu = aurora.NewAurora(false)
	}
	return logAu
}

func YellowV(part string, parts ...interface{}) interface{} {
	return au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.YellowFg)
}

// Printf print a message with formatting
func Printf(part string, parts ...interface{}) {
	hoverPrint()
	fmt.Println(fmt.Sprintf(part, parts...))
}

// Errorf print a error with formatting (red)
func Errorf(part string, parts ...interface{}) {
	hoverPrint()
	fmt.Println(au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.RedFg).String())
}

// Warnf print a warning with formatting (yellow)
func Warnf(part string, parts ...interface{}) {
	hoverPrint()
	fmt.Println(au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.YellowFg).String())
}

// Infof print a information with formatting (green)
func Infof(part string, parts ...interface{}) {
	hoverPrint()
	fmt.Println(au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.GreenFg).String())
}

// Verbosef print a verbose level information with formatting (cyan)
func Verbosef(part string, parts ...interface{}) {
	if logVerbose {
		hoverPrint()
		fmt.Println(au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.WhiteFg).String())
	}
}

func hoverPrint() {
	fmt.Print(au().Bold(au().Cyan("fvm: ")).String())
}
