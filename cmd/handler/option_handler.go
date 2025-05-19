package handler

import (
	"fmt"
	"reflect"

	flag "github.com/spf13/pflag"
)

type analysisOption struct {
	Geoip    bool
	Rdns     bool
	Protocol string
	Port     string
}

type captureOption struct {
	Count   int
	Timeout int
}

type generalOption struct {
	Help    bool
	Version bool
}

type ioOption struct {
	InterfaceName string
	OutputFile    string
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

func helpMessage() string {
	return `Usage: drawlscan [OPTIONS]

OPTIONS:
    -c, --count <NUM>              Capture only a specified number of packets
    -g, --geoip                    Show GeoIP information for source and destination IP addresses
    -h, --help                     Display this help message
    -i, --interface <INTERFACE>    Specify the network interface to capture packets from (e.g., eth0, wlan0)
    -o, --output <FILE>            Save the captured packets to a file in PCAP format
	-p, --protocol <PROTOCOL>        Filter packets by protocol (e.g., TCP, UDP, ICMP)
	-P, --port <PORT>              Filter packets by port number (e.g., 80, 443)
    -r, --rdns                     Perform reverse DNS lookups on source and destination IP addresses
    -t, --timeout <TIME>           Stop capturing after a specified number of seconds
    -v, --version                  Show version information
    --ascii                        Enable ASCII-art visualization of packets and traffic (Default is enabled)
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
	flags.Usage = func() { fmt.Println(helpMessage()) }

	// Analysis options
	flags.BoolVarP(&opts.Analysis.Geoip, "geoip", "g", false, "Show GeoIP information for source and destination IP addresses")
	flags.BoolVarP(&opts.Analysis.Rdns, "rdns", "r", false, "Perform reverse DNS lookups on source and destination IP addresses")
	flags.StringVarP(&opts.Analysis.Protocol, "protocol", "p", "all", "Filter packets by protocol (e.g., TCP, UDP, ICMP)")
	flags.StringVarP(&opts.Analysis.Port, "port", "P", "all", "Filter packets by port number (e.g., 80, 443)")

	// Capture options
	flags.IntVarP(&opts.Capture.Count, "count", "c", 0, "Capture only a specified number of packets")
	flags.IntVarP(&opts.Capture.Timeout, "timeout", "t", 0, "Stop capturing after a specified number of seconds")

	// General options
	flags.BoolVarP(&opts.General.Help, "help", "h", false, "Help message")
	flags.BoolVarP(&opts.General.Version, "version", "v", false, "Version information")

	// IO options
	flags.StringVarP(&opts.Io.InterfaceName, "interface", "i", "", "Specify the network interface to capture packets from (e.g., eth0, wlan0)")
	flags.StringVarP(&opts.Io.OutputFile, "output", "o", "", " Save the captured packets to a file in PCAP format")

	// Visualization options
	flags.BoolVar(&opts.Visualization.Ascii, "ascii", true, "Enable ASCII-art visualization of packets and traffic (Default is enable)")
	flags.BoolVar(&opts.Visualization.NoAscii, "no-ascii", false, "Disable ASCII-art output")

	return flags, opts
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

// 使い方例
func OptionHandler(optArgs []string) map[string]interface{} {
	flags, options := buildFlagSet()
	flags.Parse(optArgs[1:])

	optionMap := make(map[string]interface{})
	collectFieldMap(reflect.ValueOf(options), optionMap)
	return optionMap
}
