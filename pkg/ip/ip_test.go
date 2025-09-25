package ip_test

import (
	"math/big"
	"net"
	"testing"

	"github.com/tobiashort/th-utils/pkg/ip"
)

func TestConvertBytesToIP(t *testing.T) {
	ip := net.IP([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 2, 123})
	actual := ip.String()
	expected := "192.168.2.123"
	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestConvertIntToIP(t *testing.T) {
	i := big.NewInt(3232236033)
	actual := ip.FromInt(i).String()
	expected := "192.168.2.1"
	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
