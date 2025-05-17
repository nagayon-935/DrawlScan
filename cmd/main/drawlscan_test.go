package main

import "testing"

func Example_drawlscan(t *testing.T) {
	goMain([]string{"DrawlScan"})
	// Output:
	// Welcome to DrawlScan!
}

func TestAutoSelectInterface(t *testing.T) {
	// This test is not implemented yet.
	// You can implement it by mocking the net.Interfaces() and isInterfaceConnected() functions.
	// For now, we will just skip this test.
	t.Skip("Test not implemented yet")
}

func TestIsInterfaceConnected(t *testing.T) {
	// This test is not implemented yet.
	// You can implement it by mocking the pcap.FindAllDevs() function.
	// For now, we will just skip this test.
	t.Skip("Test not implemented yet")
}
