/*
Copyright 2018 Hank Barta

Test code for various routines in filededup.go

*/

package main

import (
	"testing"
)

func TestMin(t *testing.T) {
	if min(3, 4) == 4 {
		t.Fatal("min() returned wrong value")
	}
}
