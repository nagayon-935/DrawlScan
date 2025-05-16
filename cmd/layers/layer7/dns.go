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
	dns := packet.Layer(layers.LayerTypeDNS).(*layers.DNS)

	var answerRecoard []string
	if dns.Answers != nil {
		for _, a := range dns.Answers {
			if a.Type == layers.DNSTypeA {
				answerRecoard = append(answerRecoard, fmt.Sprintf("%s -> %s", string(a.Name), net.IP(a.IP).String()))
			}
		}
		return utils.RenderBlock("DNS Packet", answerRecoard, color.New(color.FgHiGreen))
	}
	return ""
}
