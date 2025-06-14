package main

import (
	"os"
	"testing"
)

func Test_goMain_Help(t *testing.T) {
	os.Setenv("CI", "")
	args := []string{"drawlscan", "--help"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(help) = %v, want 0", got)
	}
}

func Test_goMain_Version(t *testing.T) {
	os.Setenv("CI", "")
	args := []string{"drawlscan", "--version"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(version) = %v, want 0", got)
	}
}

func Test_goMain_CI(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping test in CI environment")
	}
	args := []string{"drawlscan"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(CI) = %v, want 0", got)
	}
	os.Setenv("CI", "")
}
