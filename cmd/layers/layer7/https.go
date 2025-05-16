package layer7

// HTTPS (TCP 443)
		if tcp.SrcPort == 443 || tcp.DstPort == 443 {
			blocks = append(blocks, renderBlock("HTTPS", []string{
				"Encrypted Traffic",
			}, color.New(color.FgHiBlue)))
		}