package layer2

import (
	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func PrintEthernetLayer(packet gopacket.Packet) string {
	var blocks []string
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return "No Ethernet layer found"
	}
	eth, _ := ethLayer.(*layers.Ethernet)

	blocks = append(blocks, utils.renderBlock("Ethernet Frame", []string{
		"Src MAC: " + eth.SrcMAC.String(),
		"Dst MAC: " + eth.DstMAC.String(),
		"Type: " + eth.EthernetType.String(),
	}, color.New(color.FgCyan)))

	return "as"
}
