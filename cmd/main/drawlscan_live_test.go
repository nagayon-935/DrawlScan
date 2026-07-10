package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

// fakeLiveHandle implements liveHandle without requiring OS capture
// permissions. It yields the given packets in order, then either closes the
// packet channel with io.EOF or, if block is true, stalls indefinitely
// (simulating a live interface with no traffic) so timeout handling can be
// exercised deterministically.
type fakeLiveHandle struct {
	packets   [][]byte
	idx       int
	block     bool
	filterErr error
	gotFilter string
}

func (f *fakeLiveHandle) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	if f.idx >= len(f.packets) {
		if f.block {
			// Not io.EOF and not a recognized terminal error, so gopacket's
			// packetsToChannel retries after a brief sleep instead of
			// closing the channel.
			return nil, gopacket.CaptureInfo{}, errors.New("no packet available")
		}
		return nil, gopacket.CaptureInfo{}, io.EOF
	}
	data = f.packets[f.idx]
	f.idx++
	ci = gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: len(data),
		Length:        len(data),
	}
	return data, ci, nil
}

func (f *fakeLiveHandle) LinkType() layers.LinkType { return layers.LinkTypeEthernet }

func (f *fakeLiveHandle) SetBPFFilter(expr string) error {
	f.gotFilter = expr
	return f.filterErr
}

// buildTestFramePackets returns n serialized Ethernet+IPv4+TCP frames.
func buildTestFramePackets(t *testing.T, n int) [][]byte {
	t.Helper()
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
	tcp := &layers.TCP{SrcPort: 12345, DstPort: 80}
	if err := tcp.SetNetworkLayerForChecksum(ip); err != nil {
		t.Fatalf("SetNetworkLayerForChecksum() error = %v", err)
	}

	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{}, eth, ip, tcp); err != nil {
		t.Fatalf("SerializeLayers() error = %v", err)
	}
	frame := buf.Bytes()

	packets := make([][]byte, n)
	for i := range packets {
		cp := make([]byte, len(frame))
		copy(cp, frame)
		packets[i] = cp
	}
	return packets
}

func Test_runLiveCapture_PacketsThenChannelClose(t *testing.T) {
	// Arrange: a fake source that yields 3 packets then closes via io.EOF.
	handle := &fakeLiveHandle{packets: buildTestFramePackets(t, 3)}

	// Act
	got := runLiveCapture(handle, "", "", 0, 0, false, false)

	// Assert: the nil packet delivered on channel close (B-2) must not panic
	// and the loop must exit cleanly once the channel closes.
	if got != 0 {
		t.Errorf("runLiveCapture() = %v, want 0", got)
	}
}

func Test_runLiveCapture_CountLimitStopsEarly(t *testing.T) {
	handle := &fakeLiveHandle{packets: buildTestFramePackets(t, 5)}

	got := runLiveCapture(handle, "", "", 2, 0, false, false)

	if got != 0 {
		t.Errorf("runLiveCapture() = %v, want 0", got)
	}
}

func Test_runLiveCapture_Timeout(t *testing.T) {
	// Arrange: a source that never produces a packet or EOF, forcing the
	// select loop to exit via the timeout branch.
	handle := &fakeLiveHandle{block: true}

	got := runLiveCapture(handle, "", "", 0, 10*time.Millisecond, false, false)

	if got != 0 {
		t.Errorf("runLiveCapture() = %v, want 0", got)
	}
}

func Test_runLiveCapture_BPFFilterApplied(t *testing.T) {
	handle := &fakeLiveHandle{packets: buildTestFramePackets(t, 1)}

	got := runLiveCapture(handle, "tcp", "", 0, 0, false, false)

	if got != 0 {
		t.Errorf("runLiveCapture() = %v, want 0", got)
	}
	if handle.gotFilter != "tcp" {
		t.Errorf("SetBPFFilter called with %q, want %q", handle.gotFilter, "tcp")
	}
}

func Test_runLiveCapture_BPFFilterError(t *testing.T) {
	handle := &fakeLiveHandle{filterErr: errors.New("bad filter")}

	got := runLiveCapture(handle, "not a valid filter", "", 0, 0, false, false)

	if got == 0 {
		t.Errorf("runLiveCapture() = %v, want != 0", got)
	}
}

func Test_runLiveCapture_WritesOutputFile(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "out.pcap")
	handle := &fakeLiveHandle{packets: buildTestFramePackets(t, 2)}

	got := runLiveCapture(handle, "", outPath, 0, 0, false, false)

	if got != 0 {
		t.Errorf("runLiveCapture() = %v, want 0", got)
	}
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("output file was not created: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("output file is empty, want header + packet data")
	}
}

func Test_runLiveCapture_InvalidOutputExtension(t *testing.T) {
	handle := &fakeLiveHandle{}

	got := runLiveCapture(handle, "", "out.txt", 0, 0, false, false)

	if got == 0 {
		t.Errorf("runLiveCapture() = %v, want != 0", got)
	}
}

func Test_runLiveCapture_OutputFileCreateFailure(t *testing.T) {
	// A directory that already exists at the target path makes os.Create fail.
	dirAsFile := filepath.Join(t.TempDir(), "out.pcap")
	if err := os.Mkdir(dirAsFile, 0o755); err != nil {
		t.Fatalf("failed to set up test directory: %v", err)
	}
	handle := &fakeLiveHandle{}

	got := runLiveCapture(handle, "", dirAsFile, 0, 0, false, false)

	if got == 0 {
		t.Errorf("runLiveCapture() = %v, want != 0", got)
	}
}

func Test_runLiveCapture_GeoIP(t *testing.T) {
	if err := utils.InitGeoIP(); err != nil {
		t.Skipf("skipping: GeoIP DB unavailable: %v", err)
	}
	defer utils.CloseGeoIP()
	handle := &fakeLiveHandle{packets: buildTestFramePackets(t, 1)}

	got := runLiveCapture(handle, "", "", 0, 0, true, true)

	if got != 0 {
		t.Errorf("runLiveCapture() = %v, want 0", got)
	}
}
