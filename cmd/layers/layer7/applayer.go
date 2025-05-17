package layer7

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/fatih/color"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

func PrintAppLayer(packet gopacket.Packet) string {
	protocol := DetectAppProtocol(packet)
	switch protocol {
	case "HTTP":
		return printHttpInfo(packet)
	case "HTTPS":
		return utils.RenderBlock("HTTPS", []string{"Encrypted Payload"}, color.New(color.FgHiCyan))
	case "QUIC":
		return utils.RenderBlock("QUIC", []string{"Encrypted Payload"}, color.New(color.FgHiMagenta))
	default:
		return ""
	}
}

func DetectAppProtocol(packet gopacket.Packet) string {

	var httpMethods = [][]byte{
		[]byte("GET "), []byte("POST "), []byte("HEAD "), []byte("PUT "),
		[]byte("DELETE "), []byte("OPTIONS "), []byte("TRACE "), []byte("CONNECT "),
	}
	var httpResponse = []byte("HTTP/")

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)

	if udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		if udp.NextLayerType() == gopacket.LayerTypePayload {
			return "QUIC"
		}
	}
	if tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		if tcp.NextLayerType() == layers.LayerTypeTLS {
			return "HTTPS"
		}
		payload := tcp.Payload
		for _, m := range httpMethods {
			if bytes.HasPrefix(payload, m) {
				return "HTTP"
			}
		}
		if bytes.HasPrefix(payload, httpResponse) {
			return "HTTP"
		}
	}
	return "Unknown"
}

// HTTPのメソッド・パス・Hostを抽出して表示
func printHttpInfo(packet gopacket.Packet) string {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return ""
	}
	tcp := tcpLayer.(*layers.TCP)
	payload := tcp.Payload
	if len(payload) == 0 {
		return ""
	}

	scanner := bufio.NewScanner(bytes.NewReader(payload))
	var method, path, host string

	// 1行目: メソッドとパス
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		if len(parts) >= 2 {
			method = parts[0]
			path = parts[1]
		}
	}
	// ヘッダからHostを探す
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Host:") {
			host = strings.TrimSpace(strings.TrimPrefix(line, "Host:"))
			break
		}
	}

	info := []string{}
	if method != "" && path != "" {
		info = append(info, "Method: "+method)
		info = append(info, "Path: "+path)
	}
	if host != "" {
		info = append(info, "Host: "+host)
	}
	return utils.RenderBlock("HTTP", info, color.New(color.FgHiYellow))
}
