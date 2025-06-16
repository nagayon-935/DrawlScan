package main

import (
	"testing"
)

func Test_goMain_Help(t *testing.T) {
	args := []string{"drawlscan", "--help"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(help) = %v, want 0", got)
	}
}

func Test_goMain_Version(t *testing.T) {
	args := []string{"drawlscan", "--version"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(version) = %v, want 0", got)
	}
}

func Test_goMain_InvalidOutputFile(t *testing.T) {
	args := []string{"drawlscan", "--output", "invalid.txt"}
	if got := goMain(args); got == 0 {
		t.Errorf("goMain(invalid output) = %v, want != 0", got)
	}
}
