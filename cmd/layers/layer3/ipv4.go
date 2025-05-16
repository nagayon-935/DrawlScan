package layer3

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintIPv4Layer(packet gopacket.Packet) string {
	ip := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	return utils.RenderBlock("IPv4 Packet", []string{
		"Src IP: " + ip.SrcIP.String(),
		"Dst IP: " + ip.DstIP.String(),
		"Protocol: " + ip.Protocol.String(),
	}, color.New(color.FgGreen))
}
