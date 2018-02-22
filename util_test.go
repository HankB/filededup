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
	options.Verbose = nil
	options.Directory = "."
	options.Trial = false
	options.Summary = false
	parseArgs()
	fmt.Println(setPrintfPri(priCritcl), options.Verbose, options.Directory, options.Trial, options.Summary)
}

func Example_parseArgs() {
	setPrintfPri(priCritcl)
	checkArgs([]string{"progname"})
	checkArgs([]string{"progname", "-v"})
	checkArgs([]string{"progname", "-vv"})
	checkArgs([]string{"progname", "--verbose", "-d", "/somedir"})
	checkArgs([]string{"progname", "-d", "/anotherdir"})
	checkArgs([]string{"progname", "-d", "/anotherdir", "--summary"})
	checkArgs([]string{"progname", "-d", "xd", "--verbose", "-t"})
	checkArgs([]string{"progname", "-t", "--dir", "."})
	checkArgs([]string{"progname", "-t", "--dir", ".", "-s"})

	// Output:
	// 0 [] . false false
	// 1 [true] . false false
	// 2 [true true] . false false
	// 1 [true] /somedir false false
	// 0 [] /anotherdir false false
	// 0 [] /anotherdir false true
	// 1 [true] xd true false
	// 0 [] . true false
	// 0 [] . true true
}
func Example_printf() {
	setPrintfPri(priCritcl)
	printf(priCritcl, "this is critical output\n")
	printf(priInfo, "this is informational output\n")
	printf(priWarn, "this is warning output\n")
	setPrintfPri(priWarn)
	printf(priCritcl, "this is critical output 2\n")
	printf(priInfo, "this is informational output 2\n")
	printf(priWarn, "this is warning output 2\n")
	setPrintfPri(priInfo)
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
