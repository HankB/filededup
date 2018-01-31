package main

/** Miscellaneous code not otherwise related to the core functionality
 */

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Options lists command line arguments
type Options struct {
	// Example of verbosity with level
	Verbose   []bool `short:"v" long:"verbose" description:"Verbose output"`
	Directory string `short:"d" long:"dir" description:"Directory to start" default:"."`
	Trial     bool   `short:"t" long:"trial" description:"report actions instead of performing operations"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func parseArgs() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

// printf() provides prioritized output using fmt.Printf
// Three level priority 0 => critical, 1 => warnings, 2 => info
type pri int

const (
	priCritcl pri = iota // print only if critical messages are displayed
	priWarn              // print errors normally expected to occur
	priInfo              // print everything including normal
)

func printf(p pri, format string, args ...interface{}) {
	if int(p) <= len(options.Verbose) {
		fmt.Printf(format, args...)
	}
}
