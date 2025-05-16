package layer7

import (
	"fmt"
	"net"
	"strings"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintDhcpLayer(packet gopacket.Packet) string {
	dhcp := packet.Layer(layers.LayerTypeDHCPv4).(*layers.DHCPv4)

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
	return utils.RenderBlock(fmt.Sprintf("DHCP %s", dhcp.Operation.String()), dhcpInfo, color.New(color.FgHiCyan))
}
