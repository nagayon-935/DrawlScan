package layer2

// ARP Layer
	if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
		arp := arpLayer.(*layers.ARP)
		blocks = append(blocks, renderBlock("ARP Packet", []string{
			"Sender MAC: " + net.HardwareAddr(arp.SourceHwAddress).String(),
			"Sender IP: " + net.IP(arp.SourceProtAddress).String(),
			"Target MAC: " + net.HardwareAddr(arp.DstHwAddress).String(),
			"Target IP: " + net.IP(arp.DstProtAddress).String(),
		}, color.New(color.FgHiYellow)))
	}