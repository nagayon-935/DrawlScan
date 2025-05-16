package layer3

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintIcmpLayer(packet gopacket.Packet) string {
	icmp := packet.Layer(layers.LayerTypeICMPv4).(*layers.ICMPv4)
	return utils.RenderBlock("ICMP Packet", []string{
		"Type: " + icmp.TypeCode.String(),
	}, color.New(color.FgHiMagenta))
}
