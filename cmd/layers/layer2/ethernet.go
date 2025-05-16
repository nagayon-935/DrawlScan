package layer2

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

/*
* @Author: nagayon935
 */
func PrintEthernetLayer(packet gopacket.Packet) string {
	eth := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
	return utils.RenderBlock("Ethernet Frame", []string{
		"Src MAC: " + eth.SrcMAC.String(),
		"Dst MAC: " + eth.DstMAC.String(),
		"Type: " + eth.EthernetType.String(),
	}, color.New(color.FgCyan))
}
