package handler

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

type AnalysisOptions struct {
	Geoip  bool
	Filter string
}

type CaptureOptions struct {
	Count int
	Time  int
}

type GeneralOptions struct {
	Help    bool
	Version bool
}

type IOOptions struct {
	InterfaceName string
	OutputFile    string
	ReadFile      string
}

type VisualizationOptions struct {
	NoAscii bool
}

// Config holds the parsed command-line options grouped by concern.
type Config struct {
	Analysis      AnalysisOptions
	Capture       CaptureOptions
	General       GeneralOptions
	IO            IOOptions
	Visualization VisualizationOptions
}

func HelpMessage() string {
	return `Usage: drawlscan [OPTIONS]

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
    -r, --read <FILE>              Read packets from a PCAP file instead of capturing live traffic
    -t, --time <TIME>              Stop capturing after a specified number of seconds
    -v, --version                  Show version information
    --no-ascii                     Disable ASCII-art output
`
}

func buildFlagSet() (*flag.FlagSet, *Config) {
	cfg := &Config{}

	flags := flag.NewFlagSet("drawlscan", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(HelpMessage()) }

	flags.BoolVarP(&cfg.Analysis.Geoip, "geoip", "g", false, "Show GeoIP information for source and destination IP addresses")
	flags.StringVarP(&cfg.Analysis.Filter, "filter", "f", "", "Filter packets")

	flags.IntVarP(&cfg.Capture.Count, "count", "c", -1, "Capture only a specified number of packets")
	flags.IntVarP(&cfg.Capture.Time, "time", "t", -1, "Stop capturing after a specified number of seconds")

	flags.BoolVarP(&cfg.General.Help, "help", "h", false, "Help message")
	flags.BoolVarP(&cfg.General.Version, "version", "v", false, "Version information")

	flags.StringVarP(&cfg.IO.InterfaceName, "interface", "i", "", "Specify the network interface to capture packets from (e.g., eth0, wlan0)")
	flags.StringVarP(&cfg.IO.OutputFile, "output", "o", "", " Save the captured packets to a file in PCAP format")
	flags.StringVarP(&cfg.IO.ReadFile, "read", "r", "", "Read packets from a PCAP file instead of capturing live traffic")

	flags.BoolVar(&cfg.Visualization.NoAscii, "no-ascii", false, "Disable ASCII-art output")
	return flags, cfg
}

// Options parses the given argument list (including argv[0]) and returns the
// resulting configuration.
func Options(optArgs []string) *Config {
	flags, cfg := buildFlagSet()
	flags.Parse(optArgs[1:])
	return cfg
}
