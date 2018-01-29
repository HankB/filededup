/*
Copyright 2018 Hank Barta

Test code for various routines in filededup_util.go

*/
package main

import (
	"fmt"
	"os"
)

func Example_psrseArgs() {
	os.Args = []string{"progname"}
	parseArgs()
	fmt.Println(len(options.Verbose), options.Directory)

	os.Args = []string{"progname", "-v"}
	parseArgs()
	fmt.Println(len(options.Verbose), options.Directory)

	os.Args = []string{"progname", "-vv"}
	parseArgs()
	fmt.Println(len(options.Verbose), options.Directory)

	os.Args = []string{"progname", "--verbose", "-d", "/somedir"}
	parseArgs()
	fmt.Println(len(options.Verbose), options.Directory)

	os.Args = []string{"progname", "-d", "/anotherdir"}
	parseArgs()
	fmt.Println(len(options.Verbose), options.Directory)

	// Output:
	// 0 .
	// 1 .
	// 2 .
	// 1 /somedir
	// 1 /anotherdir
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
