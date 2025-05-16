package layer4

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintTcpLayer(packet gopacket.Packet) string {
	tcp := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)

	return utils.RenderBlock("TCP Packet", []string{
		fmt.Sprintf("Src Port: %d", tcp.SrcPort),
		fmt.Sprintf("Dst Port: %d", tcp.DstPort),
	}, color.New(color.FgMagenta))
}
