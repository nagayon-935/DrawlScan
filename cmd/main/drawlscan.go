package main

import (
	"flag"
	"fmt"
	"os"
)

type analysisOption struct {
	detail  bool
	geoip   bool
	rdns    bool
	summary bool
}

type captureOption struct {
	count   int
	timeout int
}

type generalOption struct {
	help    bool
	version bool
}

type ioOption struct {
	interfaceName string
	outputFile    string
}

type visualizationOption struct {
	ascii   bool
	noAscii bool
}

type options struct {
	analysis      *analysisOption
	capture       *captureOption
	general       *generalOption
	io            *ioOption
	visualization *visualizationOption
}

func helpMessage() string {
	return `Usage: drawlscan [OPTIONS]

OPTIONS:
    -c, --count <NUM>              Capture only a specified number of packets
    -t, --timeout <TIME>           Stop capturing after a specified number of seconds
    -d, --detail                   Show detailed packet information, including header fields and metadata
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -r, --rdns                     Perform reverse DNS lookups on source and destination IP addresses
    -s, --summary                  Display a summary of captured packets by protocol, source, etc
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0)
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
    --ascii                        Enable ASCII-art visualization of packets and traffic (Default is enabled)
    --no-ascii                     Disable ASCII-art output
    -h, --help                     Display this help message
    -v, --version                  Show version information
`
}

func buildFlagSet() (*flag.FlagSet, *options) {
	opts := &options{
		capture:       &captureOption{},
		analysis:      &analysisOption{},
		visualization: &visualizationOption{},
		general:       &generalOption{},
		io:            &ioOption{},
	}

	flags := flag.NewFlagSet("drawlscan", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage()) }

	// Capture options
	flags.IntVar(&opts.capture.count, "c", 0, "Capture only a specified number of packets")
	flags.IntVar(&opts.capture.timeout, "t", 0, "Stop capturing after a specified number of seconds")

	// Analysis options
	flags.BoolVar(&opts.analysis.detail, "d", false, "Show detailed packet information, including header fields and metadata")
	flags.BoolVar(&opts.analysis.geoip, "g", false, "Show GeoIP information for source and destination IP addresses")
	flags.BoolVar(&opts.analysis.rdns, "r", false, "Perform reverse DNS lookups on source and destination IP addresses")
	flags.BoolVar(&opts.analysis.summary, "s", false, "Display a summary of captured packets by protocol, source, etc")

	// Visualization options
	flags.BoolVar(&opts.visualization.ascii, "ascii", true, "Enable ASCII-art visualization of packets and traffic (Default is enable)")
	flags.BoolVar(&opts.visualization.noAscii, "no-ascii", false, "Disable ASCII-art output")

	// General options
	flags.BoolVar(&opts.general.help, "h", false, "Help message")
	flags.BoolVar(&opts.general.version, "v", false, "Version information")

	// IO options
	flags.StringVar(&opts.io.interfaceName, "i", "", "Specify the network interface to capture packets from (e.g., eth0, wlan0)")
	flags.StringVar(&opts.io.outputFile, "o", "", "Save the captured packets to a file in PCAP format")

	return flags, opts
}

func hello() string {
	return "Welcome to DrawlScan!"
}

func goMain(args []string) int {
	fmt.Println(hello())
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
