package cidr

import (
	"math"
	"math/big"
	"net"

	"github.com/tobiashort/th-utils/pkg/ip"
	. "github.com/tobiashort/utils-go/must"
)

func Expand(cidr string) []net.IP {
	_, subnet := Must3(net.ParseCIDR(cidr))
	start := ip.ToInt(subnet.IP)
	ones, bits := subnet.Mask.Size()
	zeros := bits - ones
	total := math.Pow(2, float64(zeros))
	ret := make([]net.IP, int(total))
	for i := range int(total) {
		ret[i] = ip.FromInt(start)
		start.Add(start, big.NewInt(1))
	}
	return ret
}
