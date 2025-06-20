package utils

import (
	"net"
	"sort"
	"strings"

	"github.com/google/gopacket/pcap"
)

func AutoSelectInterface() string {
	ifs, _ := net.Interfaces()

	sort.Slice(ifs, func(i, j int) bool {
		return ifs[i].Index < ifs[j].Index
	})

	for _, iface := range ifs {
		if (iface.Flags&net.FlagUp != 0) && (iface.Flags&net.FlagLoopback == 0) && !strings.HasPrefix(iface.Name, "utun") {
			if isInterfaceConnected(iface.Name) {
				return iface.Name
			}
		}
	}

	return ""
}

func isInterfaceConnected(ifaceName string) bool {
	devices, _ := pcap.FindAllDevs()

	for _, dev := range devices {
		if dev.Name == ifaceName {
			if len(dev.Addresses) > 0 {
				return true
			}
		}
	}

	return false
}
