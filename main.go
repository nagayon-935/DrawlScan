package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	snapshotLen = 1024
	promiscuous = true
	timeout     = pcap.BlockForever
)

func main() {
	var iface string
	flag.StringVar(&iface, "i", "", "network interface to capture packets from")
	flag.Parse()

	if iface == "" {
		iface = autoSelectInterface()
		if iface == "" {
			log.Fatal("No suitable interface found")
		}
		fmt.Printf("Using interface: %s\n", iface)
	}

	handle, err := pcap.OpenLive(iface, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketSummary(packet)
	}
}

func autoSelectInterface() string {
	ifs, err := net.Interfaces()
	if err != nil {
		log.Fatal("Failed to get interfaces:", err)
	}

	// インターフェイスをインデックスの昇順でソート
	sort.Slice(ifs, func(i, j int) bool {
		return ifs[i].Index < ifs[j].Index
	})

	for _, iface := range ifs {

		// インターフェイスが有効で、ループバックやトンネルでないものを選択
		if (iface.Flags&net.FlagUp != 0) && (iface.Flags&net.FlagLoopback == 0) && !strings.HasPrefix(iface.Name, "utun") {
			// 接続状態を確認
			if isInterfaceConnected(iface.Name) {
				return iface.Name
			}
		}
	}

	// 適切なインターフェイスが見つからない場合
	fmt.Println("No suitable interface found.")
	return ""
}

func isInterfaceConnected(ifaceName string) bool {
	// pcap.FindAllDevs を使用して、利用可能なインターフェイスを取得
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return false
	}

	// 指定されたインターフェイスがデバイスリストに含まれているか確認
	for _, dev := range devices {
		if dev.Name == ifaceName {
			// アドレスが存在する場合、接続されているとみなす
			if len(dev.Addresses) > 0 {
				return true
			}
			return false
		}
	}

	return false
}

func printPacketSummary(packet gopacket.Packet) {
	var blocks []string

	// Ethernet Layer
	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth := ethLayer.(*layers.Ethernet)
		blocks = append(blocks, renderBlock("Ethernet Frame", []string{
			"Src MAC: " + eth.SrcMAC.String(),
			"Dst MAC: " + eth.DstMAC.String(),
			"Type: " + eth.EthernetType.String(),
		}, color.New(color.FgCyan)))
	}

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

	// IPv4 Layer
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip := ipLayer.(*layers.IPv4)
		blocks = append(blocks, renderBlock("IPv4 Header", []string{
			"Src IP: " + ip.SrcIP.String(),
			"Dst IP: " + ip.DstIP.String(),
			"Protocol: " + ip.Protocol.String(),
		}, color.New(color.FgGreen)))
	}

	// ICMP Layer
	if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		icmp := icmpLayer.(*layers.ICMPv4)
		blocks = append(blocks, renderBlock("ICMP Packet", []string{
			"Type: " + icmp.TypeCode.String(),
		}, color.New(color.FgHiMagenta)))
	}

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
					"Op: " + dhcp.Operation.String(),
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

				blocks = append(blocks, renderBlock("DHCP", dhcpInfo, color.New(color.FgHiCyan)))
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
						responses = append(responses, a.IP.String())
					} else {
						responses = append(responses, fmt.Sprintf("Type: %d", a.Type))
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

	// TCP Layer
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

	printHorizontalBlocks(blocks)
}

func renderBlock(title string, lines []string, c *color.Color) string {
	var b strings.Builder

	// 各ブロックの幅を動的に計算（色付けを無視して計算）
	maxWidth := len(title)
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	maxWidth += 2 // 両端のスペース分を追加

	// ブロックの上部（色なし）
	b.WriteString(c.Sprintf("+-%s-+\n", strings.Repeat("-", maxWidth)))

	// タイトル行（色付き）
	b.WriteString(c.Sprintf("| %-*s |\n", maxWidth, title))

	// 内容行（色付き）
	for _, line := range lines {
		b.WriteString(c.Sprintf("| %-*s |\n", maxWidth, line))
	}

	// ブロックの下部（色なし）
	b.WriteString(c.Sprintf("+-%s-+\n", strings.Repeat("-", maxWidth)))

	return b.String()
}

// func renderBlockAsJSON(title string, lines []string) string {
// 	var b strings.Builder
// 	b.WriteString("{\n")
// 	b.WriteString(fmt.Sprintf("  \"title\": \"%s\",\n", title))
// 	b.WriteString("  \"lines\": [\n")
// 	for i, line := range lines {
// 		b.WriteString(fmt.Sprintf("    \"%s\"", line))
// 		if i < len(lines)-1 {
// 			b.WriteString(",")
// 		}
// 		b.WriteString("\n")
// 	}
// 	b.WriteString("  ]\n")
// 	b.WriteString("}\n")
// 	return b.String()
// }

func printHorizontalBlocks(blocks []string) {
	if len(blocks) == 0 {
		return
	}

	// 各ブロックの行と幅を計算
	lines := make([][]string, len(blocks))
	maxWidths := make([]int, len(blocks))
	maxLines := 0

	for i, b := range blocks {
		blockLines := strings.Split(strings.TrimSuffix(b, "\n"), "\n")
		lines[i] = blockLines
		if len(blockLines) > maxLines {
			maxLines = len(blockLines)
		}
		for _, line := range blockLines {
			// ANSIエスケープシーケンスを無視して幅を計算
			visibleLength := len(stripANSI(line))
			if visibleLength > maxWidths[i] {
				maxWidths[i] = visibleLength
			}
		}
	}

	// 各行を描画
	for l := 0; l < maxLines; l++ {
		for i := 0; i < len(lines); i++ {
			if l < len(lines[i]) {
				// 行を描画（色付きのまま）
				fmt.Print(lines[i][l])
				// 空白を追加して幅を揃える
				fmt.Print(strings.Repeat(" ", maxWidths[i]-len(stripANSI(lines[i][l]))))
			} else {
				// 空行の場合は空白を出力
				fmt.Print(strings.Repeat(" ", maxWidths[i]))
			}
			fmt.Print("  ") // ブロック間のスペース
		}
		fmt.Println()
	}
	fmt.Println()
}

// ANSIエスケープシーケンスを取り除く正規表現
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(input string) string {
	return ansiEscape.ReplaceAllString(input, "")
}
