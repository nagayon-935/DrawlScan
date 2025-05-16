package layer7

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