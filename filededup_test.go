/*
Copyright 2018 Hank Barta

Test code for various routines in filededup.go

*/

package main

import (
	"fmt"
	"log"
	"os/exec"
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

func ExampleGetHash() {
	fmt.Printf("hash %x\n", getHash("sample-files/another file"))
	fmt.Printf("hash %x\n", getHash("sample-files/yet another file"))
	// Output:
	// hash b2ed2fd7ff0dc6de08c32072e40aa6bc
	// hash 2f240ab9499d7988e28288f41967a562
}

func ExampleInsertFile() {
	initDataBase("sqlite3", "test.db")
	insertFile("test file", 33, []byte("abcd"))
	out, err := exec.Command("/usr/bin/sqlite3", "test.db", "select * from files").Output()
	//out, err := exec.Command("/bin/ls", "test.db").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
	closeDataBase()
	// Output: 33|test file|abcd|1
}
