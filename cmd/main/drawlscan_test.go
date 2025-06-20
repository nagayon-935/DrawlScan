package main

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

// stdOut is used for testable output redirection.
var stdOut io.Writer = os.Stdout

func Test_goMain_Help(t *testing.T) {
	args := []string{"drawlscan", "--help"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(help) = %v, want 0", got)
	}
}

func Test_goMain_Version(t *testing.T) {
	args := []string{"drawlscan", "--version"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(version) = %v, want 0", got)
	}
}

func Test_processAndPrintPacket(t *testing.T) {
	// 標準出力をキャプチャ
	var buf bytes.Buffer
	old := stdOut
	stdOut = &buf
	defer func() { stdOut = old }()

	// Ethernet + IPv4 + TCPレイヤのダミーパケットを作成
	eth := &layers.Ethernet{
		SrcMAC:       []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		DstMAC:       []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		SrcIP:    []byte{192, 168, 1, 10},
		DstIP:    []byte{192, 168, 1, 20},
		Protocol: layers.IPProtocolTCP,
	}
	tcp := &layers.TCP{
		SrcPort: 12345,
		DstPort: 80,
	}
	_ = tcp.SetNetworkLayerForChecksum(ip)

	gpktBuf := gopacket.NewSerializeBuffer()
	_ = gopacket.SerializeLayers(gpktBuf, gopacket.SerializeOptions{}, eth, ip, tcp)
	packet := gopacket.NewPacket(gpktBuf.Bytes(), layers.LayerTypeEthernet, gopacket.Default)

	// geoip, isAscii両方falseで呼び出し
	processAndPrintPacket(packet, false, false)

	// geoip, isAscii両方trueで呼び出し
	utils.InitGeoIP()
	defer utils.CloseGeoIP()
	processAndPrintPacket(packet, true, true)
}

func Test_goMain_CountAndTimeOut(t *testing.T) {
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping this test in CI environment")
	}
	args := []string{"drawlscan", "--count", "10"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap) = %v, want 0", got)
	}
	args = []string{"drawlscan", "--time", "3"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap) = %v, want 0", got)
	}
}

func Test_goMain_ReadPcap(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping pcap test on Windows runner (wpcap.dll not always available)")
	}
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap) = %v, want 0", got)
	}
}

func Test_goMain_ReadPcap_NoAscii(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping pcap test on Windows runner (wpcap.dll not always available)")
	}
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap", "--no-ascii", "--geoip"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(read pcap no-ascii) = %v, want 0", got)
	}
}

func Test_goMain_InvalidPcapFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping pcap test on Windows runner (wpcap.dll not always available)")
	}
	args := []string{"drawlscan", "--read", "notfound.pcap"}
	if got := goMain(args); got == 0 {
		t.Errorf("goMain(invalid pcap) = %v, want != 0", got)
	}
}

func Test_goMain_InvalidOutputFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping pcap test on Windows runner (wpcap.dll not always available)")
	}
	args := []string{"drawlscan", "--output", "invalid.txt"}
	if got := goMain(args); got == 0 {
		t.Errorf("goMain(invalid output) = %v, want != 0", got)
	}
}

func Test_goMain_Timeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping pcap test on Windows runner (wpcap.dll not always available)")
	}
	args := []string{"drawlscan", "--read", "../../testdata/testdata.pcap"}
	if got := goMain(args); got != 0 {
		t.Errorf("goMain(timeout) = %v, want 0", got)
	}
}
