package layer2

import (
	"net"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

// ARP Layer
func PrintARPLayer(packet gopacket.Packet) string {

	arp := packet.Layer(layers.LayerTypeARP).(*layers.ARP)
	return utils.RenderBlock("ARP Packet", []string{
		"Sender MAC: " + net.HardwareAddr(arp.SourceHwAddress).String(),
		"Sender IP: " + net.IP(arp.SourceProtAddress).String(),
		"Target MAC: " + net.HardwareAddr(arp.DstHwAddress).String(),
		"Target IP: " + net.IP(arp.DstProtAddress).String(),
	}, color.New(color.FgHiYellow))
}
