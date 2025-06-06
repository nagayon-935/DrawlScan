package handler

import (
	"reflect"
	"testing"
)

func TestHelpMessage(t *testing.T) {
	want := `Usage: drawlscan [OPTIONS]

OPTIONS:
    -c, --count <NUM>              Capture only a specified number of packets
    -f, --filter <REGX>            Filter packets using a BPF (Berkeley Packet Filter) expression.
                                   You can specify filters such as:
                                     - ip src 192.168.1.1
                                     - ip dst 192.168.1.2
                                     - ip host 192.168.1.1 and ip host 192.168.1.2
                                     - tcp port 80
                                     - udp port 53
                                     - icmp or icmp6
                                     - vlan 100
                                     - ip host 192.168.1.1 and tcp port 80
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -h, --help                     Display this help message
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0)
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
    -t, --time <TIME>              Stop capturing after a specified number of seconds
    -v, --version                  Show version information
    --no-ascii                     Disable ASCII-art output
`
	got := HelpMessage()
	if got != want {
		t.Errorf("HelpMessage() = %v, want %v", got, want)
	}
}

func Test_buildFlagSet(t *testing.T) {
	flags, opts := buildFlagSet()
	if flags == nil {
		t.Error("buildFlagSet() flags is nil")
	}
	if opts == nil {
		t.Error("buildFlagSet() opts is nil")
	}
	// 簡易的な値チェック
	opts.Analysis.Geoip = true
	if !opts.Analysis.Geoip {
		t.Error("buildFlagSet() opts.Analysis.Geoip not settable")
	}
}

func Test_collectFieldMap(t *testing.T) {
	opts := &options{
		Capture:       &captureOption{Count: 10, Time: 5},
		Analysis:      &analysisOption{Geoip: true, Filter: "tcp"},
		Visualization: &visualizationOption{Ascii: true, NoAscii: false},
		General:       &generalOption{Help: true, Version: false},
		Io:            &ioOption{InterfaceName: "eth0", OutputFile: "out.pcap"},
	}
	optionMap := make(map[string]interface{})
	collectFieldMap(reflect.ValueOf(opts), optionMap)
	if optionMap["Geoip"] != true {
		t.Errorf("collectFieldMap() Geoip = %v, want true", optionMap["Geoip"])
	}
	if optionMap["Count"] != 10 {
		t.Errorf("collectFieldMap() Count = %v, want 10", optionMap["Count"])
	}
	if optionMap["InterfaceName"] != "eth0" {
		t.Errorf("collectFieldMap() InterfaceName = %v, want eth0", optionMap["InterfaceName"])
	}
}

func TestOptions(t *testing.T) {
	args := []string{"drawlscan", "--geoip", "--filter", "tcp", "--count", "5", "--interface", "eth0"}
	got := Options(args)
	if got["Geoip"] != true {
		t.Errorf("Options() Geoip = %v, want true", got["Geoip"])
	}
	if got["Filter"] != "tcp" {
		t.Errorf("Options() Filter = %v, want tcp", got["Filter"])
	}
	if got["Count"] != 5 {
		t.Errorf("Options() Count = %v, want 5", got["Count"])
	}
	if got["InterfaceName"] != "eth0" {
		t.Errorf("Options() InterfaceName = %v, want eth0", got["InterfaceName"])
	}
}
