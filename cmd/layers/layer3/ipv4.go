package layer3

// IPv4 Layer
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip := ipLayer.(*layers.IPv4)
		blocks = append(blocks, renderBlock("IPv4 Header", []string{
			"Src IP: " + ip.SrcIP.String(),
			"Dst IP: " + ip.DstIP.String(),
			"Protocol: " + ip.Protocol.String(),
		}, color.New(color.FgGreen)))
	}