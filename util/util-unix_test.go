// +build !windows

package util

import (
	"fmt"
	"testing"
)

func TestOsAwareFilePath(t *testing.T) {
	path := "this/is/mypath"
	fmt.Println(path)
	output := OsAwareFilePath(path)

	if output != "this/is/mypath" {
		t.Errorf("Output is incorrect, got: %s, want: %s.", output, "this/is/mypath")
	}
}

func TestOsAwareFilePathMixed(t *testing.T) {
	path := "this/is\\mypath"
	fmt.Println(path)
	output := OsAwareFilePath(path)

	if output != "this/is/mypath" {
		t.Errorf("Output is incorrect, got: %s, want: %s.", output, "this/is/mypath")
	}
}
