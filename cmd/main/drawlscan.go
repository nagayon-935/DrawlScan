package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type analysisOption struct {
	detail   bool
	geoip    bool
	rdns     bool
	summary  bool
	protocol string
	port     string
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
    -d, --detail                   Show detailed packet information, including header fields and metadata
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -h, --help                     Display this help message
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0)
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
	-p, --protocol <PROTOCOL>        Filter packets by protocol (e.g., TCP, UDP, ICMP)
	-P, --port <PORT>              Filter packets by port number (e.g., 80, 443)
    -s, --summary                  Display a summary of captured packets by protocol, source, etc
    -r, --rdns                     Perform reverse DNS lookups on source and destination IP addresses
    -t, --timeout <TIME>           Stop capturing after a specified number of seconds
    -v, --version                  Show version information
    --ascii                        Enable ASCII-art visualization of packets and traffic (Default is enabled)
    --no-ascii                     Disable ASCII-art output
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

	// Analysis options
	flags.BoolVarP(&opts.analysis.detail, "detail", "d", false, "Show detailed packet information, including header fields and metadata")
	flags.BoolVarP(&opts.analysis.geoip, "geoip", "g", false, "Show GeoIP information for source and destination IP addresses")
	flags.BoolVarP(&opts.analysis.rdns, "rdns", "r", false, "Perform reverse DNS lookups on source and destination IP addresses")
	flags.BoolVarP(&opts.analysis.summary, "summary", "s", false, "Display a summary of captured packets by protocol, source, etc")
	flags.StringVarP(&opts.analysis.protocol, "protocol", "p", "", "Filter packets by protocol (e.g., TCP, UDP, ICMP)")
	flags.StringVarP(&opts.analysis.port, "port", "P", "", "Filter packets by port number (e.g., 80, 443)")

	// Capture options
	flags.IntVarP(&opts.capture.count, "count", "c", 0, "Capture only a specified number of packets")
	flags.IntVarP(&opts.capture.timeout, "timeout", "t", 0, "Stop capturing after a specified number of seconds")

	// General options
	flags.BoolVarP(&opts.general.help, "help", "h", false, "Help message")
	flags.BoolVarP(&opts.general.version, "version", "v", false, "Version information")

	// IO options
	flags.StringVarP(&opts.io.interfaceName, "interface", "i", "", "Specify the network interface to capture packets from (e.g., eth0, wlan0)")
	flags.StringVarP(&opts.io.outputFile, "output", "o", "", " Save the captured packets to a file in PCAP format")

	// Visualization options
	flags.BoolVar(&opts.visualization.ascii, "ascii", true, "Enable ASCII-art visualization of packets and traffic (Default is enable)")
	flags.BoolVar(&opts.visualization.noAscii, "no-ascii", false, "Disable ASCII-art output")

	return flags, opts
}

func hello() string {
	return "Welcome to DrawlScan!"
}

func goMain(args []string) int {
	fmt.Println(hello())
	flags, opts := buildFlagSet()
	flags.Parse(args[1:])
	fmt.Println("Parsed options:", opts.io.interfaceName)
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
