// UDP Layer
if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
	udp := udpLayer.(*layers.UDP)
	blocks = append(blocks, renderBlock("UDP Datagram", []string{
		fmt.Sprintf("Src Port: %d", udp.SrcPort),
		fmt.Sprintf("Dst Port: %d", udp.DstPort),
	}, color.New(color.FgBlue)))

	// DHCP (UDP 67, 68)
	if udp.SrcPort == 67 || udp.DstPort == 67 || udp.SrcPort == 68 || udp.DstPort == 68 {
		if dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4); dhcpLayer != nil {
			dhcp := dhcpLayer.(*layers.DHCPv4)

			// 必要なオプションを直接検索
			var subnetMask, router, domainName, leaseTime, serverIdentifier string
			var domainNameServers []string
			for _, opt := range dhcp.Options {
				switch opt.Type {
				case 1: // Subnet Mask
					if len(opt.Data) == 4 {
						subnetMask = net.IP(opt.Data).String()
					}
				case 3: // Router
					if len(opt.Data) >= 4 {
						router = net.IP(opt.Data[:4]).String()
					}
				case 6: // Domain Name Server
					for i := 0; i+4 <= len(opt.Data); i += 4 {
						domainNameServers = append(domainNameServers, net.IP(opt.Data[i:i+4]).String())
					}
				case 15: // Domain Name
					domainName = string(opt.Data)
				case 51: // Lease Time
					if len(opt.Data) == 4 {
						lease := uint32(opt.Data[0])<<24 | uint32(opt.Data[1])<<16 | uint32(opt.Data[2])<<8 | uint32(opt.Data[3])
						leaseTime = fmt.Sprintf("%d seconds", lease)
					}
				case 54: // DHCP Server Identifier
					if len(opt.Data) == 4 {
						serverIdentifier = net.IP(opt.Data).String()
					}
				}
			}

			// DHCP情報をブロックに追加
			dhcpInfo := []string{
				"Your IP: " + dhcp.YourClientIP.String(),
				"Subnet Mask: " + subnetMask,
				"Router: " + router,
				"Lease Time: " + leaseTime,
				"Server Identifier: " + serverIdentifier,
			}
			if domainName != "" {
				dhcpInfo = append(dhcpInfo, "Domain Name: "+domainName)
			}
			if len(domainNameServers) > 0 {
				dhcpInfo = append(dhcpInfo, "Domain Name Servers: "+strings.Join(domainNameServers, ", "))
			}

			blocks = append(blocks, renderBlock(fmt.Sprintf("DHCP %s", dhcp.Operation.String()), dhcpInfo, color.New(color.FgHiCyan)))
		}
	}

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
}