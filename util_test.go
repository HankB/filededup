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
