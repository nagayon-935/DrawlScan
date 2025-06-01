package utils

import (
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/google/gopacket/pcap"
)

func AutoSelectInterface() string {
	ifs, err := net.Interfaces()
	if err != nil {
		fmt.Println("Failed to get interfaces:", err)
	}

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

	fmt.Println("No suitable interface found.")
	return ""
}

func isInterfaceConnected(ifaceName string) bool {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return false
	}

	for _, dev := range devices {
		if dev.Name == ifaceName {
			if len(dev.Addresses) > 0 {
				return true
			}
		}
	}

	return false
}
