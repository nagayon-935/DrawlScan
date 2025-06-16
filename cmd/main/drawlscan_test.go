package main

import (
	"os"
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

func Test_goMain_CI(t *testing.T) {
	os.Setenv("CI", "true")
	args := []string{"drawlscan"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(CI) = %v, want 0", got)
	}
}
