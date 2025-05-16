package layer3

func printICMPLayer(packet gopacket.Packet) string {
	blocks := []string{}

	// ICMP Layer
	if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		icmp := icmpLayer.(*layers.ICMPv4)
		blocks = append(blocks, renderBlock("ICMP Packet", []string{
			"Type: " + icmp.TypeCode.String(),
		}, color.New(color.FgHiMagenta)))
	}