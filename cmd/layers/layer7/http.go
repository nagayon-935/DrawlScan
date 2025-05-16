package layer7

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