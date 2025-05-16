package main

import (
	"github.com/google/gopacket"
	"github.com/nagayon-935/DrawlScan/cmd/handler"
)

func printPacket(packetSource *gopacket.PacketSource) {
	for packet := range packetSource.Packets() {
		for _, h := range handler.Handlers {
			if packet.Layer(h.LayerType) != nil {
				h.Handler(packet)
			}
		}
	}
}
