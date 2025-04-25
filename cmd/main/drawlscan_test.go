package main

import "testing"

func Example_wildcherry() {
	goMain([]string{"DrawlScan"})
	// Output:
	// Welcome to WildCherry!
}

func TestHello(t *testing.T) {
	got := hello()
	want := "Welcome to DrawlScan!"
	if got != want {
		t.Errorf("hello() = %q, want %q", got, want)
	}
}
