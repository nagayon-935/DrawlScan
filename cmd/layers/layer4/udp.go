package layer4

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintUdpLayer(packet gopacket.Packet) string {
	udp := packet.Layer(layers.LayerTypeUDP).(*layers.UDP)
	return utils.RenderBlock("UDP Datagram", []string{
		fmt.Sprintf("Src Port: %d", udp.SrcPort),
		fmt.Sprintf("Dst Port: %d", udp.DstPort),
	}, color.New(color.FgBlue))
}
