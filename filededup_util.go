package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Options lists command line arguments
type Options struct {
	// Example of verbosity with level
	Verbose   []bool `short:"v" long:"verbose" description:"Verbose output"`
	Directory string `short:"d" long:"dir" description:"Directory to start" default:"."`
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
