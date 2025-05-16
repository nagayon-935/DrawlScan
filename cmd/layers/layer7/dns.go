package layer7

import (
	"fmt"
	"net"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintDnsLayer(packet gopacket.Packet) string {
	// DNS (UDP 53)
	dns := packet.Layer(layers.LayerTypeDNS).(*layers.DNS)

	if dns.Questions != nil {
		return printDnsQuestion(dns)
	}

	if dns.Answers != nil {
		return printDnsResponse(dns)
	}
	return "None"
}

func printDnsQuestion(dns *layers.DNS) string {
	var queries []string
	for _, q := range dns.Questions {
		queries = append(queries, string(q.Name))
	}
	return utils.RenderBlock("DNS Queries", queries, color.New(color.FgHiYellow))
}

func printDnsResponse(dns *layers.DNS) string {
	var responses []string
	for _, a := range dns.Answers {
		if a.Type == layers.DNSTypeA {
			responses = append(responses, fmt.Sprintf("%s record: %s", a.Type.String(), net.IP(a.IP).String()))
		}
	}
	return utils.RenderBlock("DNS Responses", responses, color.New(color.FgHiGreen))
}
