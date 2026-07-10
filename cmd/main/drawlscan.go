package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/nagayon-935/DrawlScan/cmd/handler"
	"github.com/nagayon-935/DrawlScan/cmd/utils"
)

// liveHandle is the subset of *pcap.Handle used by runLiveCapture. It lets
// tests inject a fake packet source instead of requiring OS capture
// permissions.
type liveHandle interface {
	gopacket.PacketDataSource
	LinkType() layers.LinkType
	SetBPFFilter(expr string) error
}

func processAndPrintPacket(packet gopacket.Packet, geoip bool, isAscii bool) {
	var blocks []string
	for _, h := range handler.Handlers {
		if packet.Layer(h.LayerType) != nil {
			block := h.Handler(packet)
			if block != "" && block != "invisible" {
				blocks = append(blocks, block)
			}
		}
	}
	if geoip {
		if netLayer := packet.NetworkLayer(); netLayer != nil {
			src, dst := netLayer.NetworkFlow().Endpoints()
			srcGeo := utils.LookupCountry(src.String())
			if srcGeo != "" && srcGeo != "invisible" {
				blocks = append(blocks, srcGeo)
			}
			dstGeo := utils.LookupCountry(dst.String())
			if dstGeo != "" && dstGeo != "invisible" {
				blocks = append(blocks, dstGeo)
			}
		}
	}
	if isAscii {
		utils.PrintHorizontalBlocks(blocks)
	} else {
		fmt.Println(packet)
	}
}

func goMain(args []string) int {
	cfg := handler.Options(args)
	var (
		count         = cfg.Capture.Count
		filter        = cfg.Analysis.Filter
		geoip         = cfg.Analysis.Geoip
		help          = cfg.General.Help
		iface         = cfg.IO.InterfaceName
		writeFilePath = cfg.IO.OutputFile
		readFilePath  = cfg.IO.ReadFile
		timeSec       = cfg.Capture.Time
		version       = cfg.General.Version
		isAscii       = !cfg.Visualization.NoAscii
	)

	if help {
		fmt.Println(handler.HelpMessage())
		return 0
	}

	if version {
		fmt.Println("Version: " + VERSION)
		return 0
	}

	if geoip {
		if err := utils.InitGeoIP(); err != nil {
			fmt.Println("Failed to initialize GeoIP database:", err)
			return 1
		}
		defer utils.CloseGeoIP()
	}

	if readFilePath != "" {
		handle, err := pcap.OpenOffline(readFilePath)
		if err != nil {
			fmt.Println("Failed to open pcap file:", err)
			return 1
		}
		defer handle.Close()
		if filter != "" {
			if err := handle.SetBPFFilter(filter); err != nil {
				fmt.Println("Failed to set BPF filter: ", err)
				return 1
			}
		}
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		received := 0
		for packet := range packetSource.Packets() {
			processAndPrintPacket(packet, geoip, isAscii)
			received++
			if count > 0 && received >= count {
				break
			}
		}
		fmt.Printf("Read %d packets from %s\n", received, readFilePath)
		return 0
	}

	if iface == "" {
		iface = utils.AutoSelectInterface()
		if iface == "" {
			fmt.Println("No suitable interface found")
			return 1
		}
	}

	handle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer handle.Close()

	fmt.Println("Using interface: ", iface)

	var timeout time.Duration
	if timeSec > 0 {
		timeout = time.Duration(timeSec) * time.Second
	}

	return runLiveCapture(handle, filter, writeFilePath, count, timeout, geoip, isAscii)
}

// runLiveCapture applies the BPF filter (if any), optionally mirrors captured
// packets to writeFilePath, and processes packets from handle until the
// packet channel closes, count packets have been received, or timeout
// elapses (timeout <= 0 disables the timeout).
func runLiveCapture(handle liveHandle, filter string, writeFilePath string, count int, timeout time.Duration, geoip bool, isAscii bool) int {
	if filter != "" {
		if err := handle.SetBPFFilter(filter); err != nil {
			fmt.Println("Failed to set BPF filter: ", err)
			return 1
		}
	}

	var distFile *os.File
	var pcapw *pcapgo.Writer
	if writeFilePath != "" {
		ext := strings.ToLower(filepath.Ext(writeFilePath))
		if ext != ".pcap" && ext != ".pcapng" {
			fmt.Println("Output file must have .pcap or .pcapng extension: ", writeFilePath)
			return 1
		}
		var err error
		distFile, err = os.Create(writeFilePath)
		if err != nil {
			fmt.Println("Failed to create output file: ", err)
			return 1
		}
		defer distFile.Close()
		pcapw = pcapgo.NewWriter(distFile)
		if err := pcapw.WriteFileHeader(1600, handle.LinkType()); err != nil {
			fmt.Println("WriteFileHeader: ", err)
			return 1
		}
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetChan := packetSource.Packets()

	var timeCh <-chan time.Time
	if timeout > 0 {
		timeCh = time.After(timeout)
	}
	received := 0

	done := false
	start := time.Now()

	for !done {
		select {
		case packet, ok := <-packetChan:
			if !ok {
				done = true
				continue
			}
			if distFile != nil {
				if err := pcapw.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
					fmt.Println("pcap.WritePacket(): ", err)
					return 1
				}
			}
			processAndPrintPacket(packet, geoip, isAscii)
			received++
		case <-timeCh:
			done = true
		}
		if count > 0 && received >= count {
			done = true
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Captured %d packets\n", received)
	fmt.Printf("Capture duration: %.2f seconds\n", elapsed.Seconds())
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
