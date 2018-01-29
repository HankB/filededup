/*
Copyright 2018 Hank Barta

Test code for various routines in filededup.go

*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestMin(t *testing.T) {
	if min(3, 4) != 3 {
		t.Fatal("min(3, 4) did not return returned 3")
	}
	if min(5, 4) != 4 {
		t.Fatal("min(3, 4) did not return returned 4")
	}
	if min(-1, -2) != -2 {
		t.Fatal("min(-1, -2) did not return returned '-2'")
	}
	if min(3, 3) != 3 {
		t.Fatal("min(-1, -2) did not return returned '3'")
	}
}

func Example_getHash() {
	fmt.Printf("hash %x\n", getHash("sample-files/another file"))
	fmt.Printf("hash %x\n", getHash("sample-files/yet another file"))
	// Output:
	// hash b2ed2fd7ff0dc6de08c32072e40aa6bc
	// hash 2f240ab9499d7988e28288f41967a562
}

func dumpDatabase() {
	out, err := exec.Command("/usr/bin/sqlite3", "test.db",
		"select length, filename, HEX(hash), linkCount from files").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)
}

func Example_insertFile() {
	initDataBase("sqlite3", "test.db")
	defer closeDataBase()

	insertFile("test file", 33, []byte("abcd"))
	insertFile("test file 3", 333, nil)
	dumpDatabase()

	updateHash("test file", []byte("cba"))
	updateHash("test file 3", []byte("xyz"))
	dumpDatabase()

	// Output:
	// 33|test file|61626364|1
	// 333|test file 3||1
	// 33|test file|636261|1
	// 333|test file 3|78797A|1
}

func TestCompareByteByByte(t *testing.T) {

	if err := exec.Command("/bin/sh", "./prep_compare_files.sh").Run(); err != nil {
		log.Fatal(err)
	}

	if !compareByteByByte("cmpfile.4096-1", "cmpfile.4096-2", 4096) {
		t.Fatal("cmpfile.4096-1, cmpfile.4096-2 do not match")
	}
	if compareByteByByte("cmpfile.4096-1", "cmpfile.4096-3", 4096) {
		t.Fatal("cmpfile.4096-1, cmpfile.4096-3 match")
	}

	if !compareByteByByte("cmpfile.1024-1", "cmpfile.1024-2", 1024) {
		t.Fatal("cmpfile.1024-1, cmpfile.1024-2 do not match")
	}
	if compareByteByByte("cmpfile.1024-1", "cmpfile.1024-3", 1024) {
		t.Fatal("cmpfile.1024-1, cmpfile.1024-3 match")
	}

	if !compareByteByByte("cmpfile.5120-1", "cmpfile.5120-2", 5120) {
		t.Fatal("cmpfile.5120-1, cmpfile.5120-2 do not match")
	}
	if compareByteByByte("cmpfile.5120-1", "cmpfile.5120-3", 5120) {
		t.Fatal("cmpfile.5120-1, cmpfile.5120-3 match")
	}

	if err = exec.Command("/bin/sh", "./rm_compare_files.sh").Run(); err != nil {
		log.Fatal(err)
	}

}

func checkLink(foo, baz string) bool {
	fooInfo, err := os.Stat(foo)
	if err != nil {
		log.Printf("can't Stat %s\n", foo)
	}
	bazInfo, err := os.Stat(baz)
	if err != nil {
		log.Printf("can't Stat %s\n", baz)
	}

	fooStat, ok := fooInfo.Sys().(*syscall.Stat_t)
	if !ok {
		log.Printf("Not a syscall.Stat_t for %s\n", foo)
		return false
	}
	bazStat, ok := bazInfo.Sys().(*syscall.Stat_t)
	if !ok {
		log.Printf("Not a syscall.Stat_t for %s\n", baz)
		return false
	}
	return fooStat.Ino == bazStat.Ino
}

func TestLinkFile(t *testing.T) {

	if err := exec.Command("/bin/sh", "./prep_link_files.sh").Run(); err != nil {
		log.Fatal(err)
	}

	replaceWithLink("a", "b")
	if !checkLink("a", "b") {
		t.Fatal("\"a\" \"b\" not linked\n")
	}

	if err := exec.Command("/bin/sh", "./rm_link_files.sh").Run(); err != nil {
		log.Fatal(err)
	}

}

func Example_findMatch() {

	if err := exec.Command("/bin/sh", "./prep_findmatch_files.sh").Run(); err != nil {
		log.Fatal(err)
	}

	initDataBase("sqlite3", "test.db")
	defer closeDataBase()

	fileNames := []string{"another file", "another file.copy", "empty", "empty-01",
		"README.md", "thing one", "thing two", "yet another file", "x", "y", "z"}

	for _, fname := range fileNames {
		info, err := os.Stat("sample-files/" + fname)
		if err != nil {
			log.Fatal(err)
		}
		matched, matchName, hash := findMatch("sample-files/"+fname, info)
		fmt.Printf("%s, %t, %s, %x .\n", fname, matched, matchName, hash)
		if !matched {
			insertFile("sample-files/"+fname, info.Size(), hash)
		}

	}
	if err := exec.Command("/bin/sh", "./rm_findmatch_files.sh").Run(); err != nil {
		log.Fatal(err)
	}
	dumpDatabase()

	// Output:
	// another file, false, ,  .
	// another file.copy, true, sample-files/another file, b2ed2fd7ff0dc6de08c32072e40aa6bc .
	// empty, false, ,  .
	// empty-01, true, sample-files/empty, d41d8cd98f00b204e9800998ecf8427e .
	// README.md, false, ,  .
	// thing one, false, ,  .
	// thing two, false, , 008ee33a9d58b51cfeb425b0959121c9 .
	// yet another file, false, , 2f240ab9499d7988e28288f41967a562 .
	// x, false, ,  .
	// y, false, ,  .
	// z, true, sample-files/x, 401b30e3b8b5d629635a5c613cdb7919 .
	// 22|sample-files/another file|B2ED2FD7FF0DC6DE08C32072E40AA6BC|1
	// 0|sample-files/empty||1
	// 440|sample-files/README.md||1
	// 64|sample-files/thing one|008EE33A9D58B51CFEB425B0959121C9|1
	// 64|sample-files/thing two|008EE33A9D58B51CFEB425B0959121C9|1
	// 22|sample-files/yet another file|2F240AB9499D7988E28288F41967A562|1
	// 2|sample-files/x||1
	// 2|sample-files/y||1
}

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
