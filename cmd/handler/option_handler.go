package handler

import (
	"fmt"
	"reflect"

	flag "github.com/spf13/pflag"
)

type analysisOption struct {
	Geoip  bool
	Filter string
}

type captureOption struct {
	Count int
	Time  int
}

type generalOption struct {
	Help    bool
	Version bool
}

type ioOption struct {
	InterfaceName string
	OutputFile    string
	ReadFile      string
}

type visualizationOption struct {
	Ascii   bool
	NoAscii bool
}

type options struct {
	Analysis      *analysisOption
	Capture       *captureOption
	General       *generalOption
	Io            *ioOption
	Visualization *visualizationOption
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

func buildFlagSet() (*flag.FlagSet, *options) {
	opts := &options{
		Capture:       &captureOption{},
		Analysis:      &analysisOption{},
		Visualization: &visualizationOption{},
		General:       &generalOption{},
		Io:            &ioOption{},
	}

	flags := flag.NewFlagSet("drawlscan", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(HelpMessage()) }

	flags.BoolVarP(&opts.Analysis.Geoip, "geoip", "g", false, "Show GeoIP information for source and destination IP addresses")
	flags.StringVarP(&opts.Analysis.Filter, "filter", "f", "", "Filter packets")

	flags.IntVarP(&opts.Capture.Count, "count", "c", -1, "Capture only a specified number of packets")
	flags.IntVarP(&opts.Capture.Time, "time", "t", -1, "Stop capturing after a specified number of seconds")

	flags.BoolVarP(&opts.General.Help, "help", "h", false, "Help message")
	flags.BoolVarP(&opts.General.Version, "version", "v", false, "Version information")

	flags.StringVarP(&opts.Io.InterfaceName, "interface", "i", "", "Specify the network interface to capture packets from (e.g., eth0, wlan0)")
	flags.StringVarP(&opts.Io.OutputFile, "output", "o", "", " Save the captured packets to a file in PCAP format")
	flags.StringVarP(&opts.Io.ReadFile, "read", "r", "", "Read packets from a PCAP file instead of capturing live traffic")

	flags.BoolVar(&opts.Visualization.NoAscii, "no-ascii", false, "Disable ASCII-art output")
	return flags, opts
}

func Options(optArgs []string) map[string]interface{} {
	flags, options := buildFlagSet()
	flags.Parse(optArgs[1:])
	optionMap := make(map[string]interface{})
	collectFieldMap(reflect.ValueOf(options), optionMap)
	return optionMap
}

func collectFieldMap(value reflect.Value, optionMap map[string]interface{}) {
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	valueType := value.Type()

	fields := reflect.VisibleFields(valueType)
	for _, field := range fields {
		fieldValue := value.FieldByIndex(field.Index)
		key := field.Name
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Struct) {
			collectFieldMap(fieldValue, optionMap)
		} else {
			optionMap[key] = fieldValue.Interface()
		}
	}
}
