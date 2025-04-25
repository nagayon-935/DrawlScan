package main

import "testing"

func Example_wildcherry() {
	goMain([]string{"DrawlScan"})
	// Output:
	// Welcome to DrawlScan!
}

func TestHello(t *testing.T) {
	got := hello()
	want := "Welcome to DrawlScan!"
	if got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}
