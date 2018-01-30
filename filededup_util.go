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
const priCritcl = 0 // print only if critical messages are displayed
const priWarn = 1   // print errors normally expected to occur
const priInfo = 2   // print everything including normal

func printf(pri int, format string, args ...interface{}) {
	if pri <= len(options.Verbose) {
		fmt.Printf(format, args...)
	}
}
