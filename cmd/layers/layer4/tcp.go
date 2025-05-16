package layer4

if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		blocks = append(blocks, renderBlock("TCP Segment", []string{
			fmt.Sprintf("Src Port: %d", tcp.SrcPort),
			fmt.Sprintf("Dst Port: %d", tcp.DstPort),
		}, color.New(color.FgMagenta)))

		// HTTP (TCP 80)
		if tcp.SrcPort == 80 || tcp.DstPort == 80 {
			if app := packet.ApplicationLayer(); app != nil {
				payload := string(app.Payload())
				if len(payload) > 64 {
					payload = payload[:64] + "..."
				}
				blocks = append(blocks, renderBlock("HTTP", []string{payload}, color.New(color.FgHiWhite)))
			}
		}

		// HTTPS (TCP 443)
		if tcp.SrcPort == 443 || tcp.DstPort == 443 {
			blocks = append(blocks, renderBlock("HTTPS", []string{
				"Encrypted Traffic",
			}, color.New(color.FgHiBlue)))
		}
	}