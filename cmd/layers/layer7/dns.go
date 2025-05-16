	// DNS (UDP 53)
	if udp.SrcPort == 53 || udp.DstPort == 53 {
		if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
			dns := dnsLayer.(*layers.DNS)

			// DNS Queries
			var queries []string
			for _, q := range dns.Questions {
				queries = append(queries, string(q.Name))
			}

			// DNS Responses
			var responses []string
			for _, a := range dns.Answers {
				if a.Type == layers.DNSTypeA { // A Record (IPv4)
					responses = append(responses, fmt.Sprintf("%s record: %s", a.Type.String(), net.IP(a.IP).String()))
				}
			}

			// Add DNS Queries and Responses to blocks
			if len(queries) > 0 {
				blocks = append(blocks, renderBlock("DNS Queries", queries, color.New(color.FgHiYellow)))
			}
			if len(responses) > 0 {
				blocks = append(blocks, renderBlock("DNS Responses", responses, color.New(color.FgHiGreen)))
			}
		}
	}