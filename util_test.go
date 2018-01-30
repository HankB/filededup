/*
Copyright 2018 Hank Barta

Test code for various routines in filededup_util.go

*/
package main

import (
	"fmt"
	"os"
)

func checkArgs(args []string) {
	os.Args = args
	parseArgs()
	fmt.Println(len(options.Verbose), options.Directory, options.Trial)
}

func Example_parseArgs() {
	checkArgs([]string{"progname"})
	checkArgs([]string{"progname", "-v"})
	checkArgs([]string{"progname", "-vv"})
	checkArgs([]string{"progname", "--verbose", "-d", "/somedir"})
	checkArgs([]string{"progname", "-d", "/anotherdir"})
	checkArgs([]string{"progname", "-t", "--dir", "."})

	// Output:
	// 0 . false
	// 1 . false
	// 2 . false
	// 1 /somedir false
	// 1 /anotherdir false
	// 1 . true
}
func Example_printf() {
	options.Verbose = []bool{}
	printf(priCritcl, "this is critical output\n")
	printf(priInfo, "this is informational output\n")
	printf(priWarn, "this is warning output\n")
	options.Verbose = []bool{true}
	printf(priCritcl, "this is critical output 2\n")
	printf(priInfo, "this is informational output 2\n")
	printf(priWarn, "this is warning output 2\n")
	options.Verbose = []bool{true, true}
	printf(priCritcl, "this is critical output 3\n")
	printf(priInfo, "this is informational output 3\n")
	printf(priWarn, "this is warning output 3\n")

	// Output:
	// this is critical output
	// this is critical output 2
	// this is warning output 2
	// this is critical output 3
	// this is informational output 3
	// this is warning output 3
}
