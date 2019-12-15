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

// Au Aurora instance used for colors
func Au() aurora.Aurora {
  if logAu == nil {
    logAu = aurora.NewAurora(false)
  }
  return logAu
}

// Printf print a message with formatting
func Printf(part string, parts ...interface{}) {
  hoverPrint()
  fmt.Println(fmt.Sprintf(part, parts...))
}

// Errorf print a error with formatting (red)
func Errorf(part string, parts ...interface{}) {
  hoverPrint()
  fmt.Println(Au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.RedFg).String())
}

// Warnf print a warning with formatting (yellow)
func Warnf(part string, parts ...interface{}) {
  hoverPrint()
  fmt.Println(Au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.YellowFg).String())
}

// Infof print a information with formatting (green)
func Infof(part string, parts ...interface{}) {
  hoverPrint()
  fmt.Println(Au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.GreenFg).String())
}

// Verbosef print a verbose level information with formatting (cyan)
func Verbosef(part string, parts ...interface{}) {
  if logVerbose {
    hoverPrint()
    fmt.Println(Au().Colorize(fmt.Sprintf(fmt.Sprintf("%v", part), parts...), aurora.WhiteFg).String())
  }
}

func hoverPrint() {
  fmt.Print(Au().Bold(Au().Cyan("fvm: ")).String())
}
