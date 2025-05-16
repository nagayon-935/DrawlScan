package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/layers/layer2"
)

func printPacket(packet *gopacket.PacketSource) {
	for packet := range packet.Packets() {
		if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
			layer2.PrintEthernetLayer(packet)
		}
	}

}
