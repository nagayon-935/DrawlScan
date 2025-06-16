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

func Test_goMain_ReadPcap(t *testing.T) {
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap) = %v, want 0", got)
	}
}

func Test_goMain_ReadPcap_NoAscii(t *testing.T) {
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap", "--no-ascii"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap no-ascii) = %v, want 0", got)
	}
}

func Test_goMain_ReadPcap_Filter(t *testing.T) {
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap", "--filter", "tcp"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap filter) = %v, want 0", got)
	}
}

func Test_goMain_InvalidPcapFile(t *testing.T) {
	args := []string{"drawlscan", "--read", "notfound.pcap"}
	if got := goMain(args); got == 0 {
		t.Errorf("goMain(invalid pcap) = %v, want != 0", got)
	}
}

func Test_goMain_InvalidOutputFile(t *testing.T) {
	args := []string{"drawlscan", "--output", "invalid.txt"}
	if got := goMain(args); got == 0 {
		t.Errorf("goMain(invalid output) = %v, want != 0", got)
	}
}
func Test_goMain_Timeout(t *testing.T) {
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap", "--time", "1"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(timeout) = %v, want 0", got)
	}
}
