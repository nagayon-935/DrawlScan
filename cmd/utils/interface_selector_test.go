package utils

import (
	"testing"
)

func TestAutoSelectInterface(t *testing.T) {
	iface := AutoSelectInterface()
	t.Logf("AutoSelectInterface() returned: %q", iface)
}

func Test_isInterfaceConnected(t *testing.T) {
	if got := isInterfaceConnected("nonexistent0"); got != false {
		t.Errorf("isInterfaceConnected(nonexistent0) = %v, want false", got)
	}
}
