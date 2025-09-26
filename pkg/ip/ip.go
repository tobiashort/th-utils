package ip

import (
	"fmt"
	"math"
	"math/big"
	"net"
	"slices"
	"strings"
)

func ToInt(ip net.IP) *big.Int {
	return big.NewInt(0).SetBytes(ip)
}

func FromInt(n *big.Int) net.IP {
	if len(n.Bytes()) > 16 {
		panic("too many bytes")
	}

	padding := 16 - len(n.Bytes())
	bs := make([]byte, 16)
	for i, b := range n.Bytes() {
		bs[i+padding] = b
	}

	v4 := true
	for i := range 10 {
		if bs[i] != 0 {
			v4 = false
			break
		}
	}

	if v4 {
		bs[11] = 255
		bs[10] = 255
		return net.IP(bs).To4()
	} else {
		return net.IP(bs).To16()
	}
}

func Collapse(ips []net.IP) []string {
	ipsAsInt := make([]*big.Int, len(ips))
	for i, ip := range ips {
		ipsAsInt[i] = ToInt(ip)
	}
	slices.SortFunc(ipsAsInt, func(a, b *big.Int) int { return a.Cmp(b) })

	var ranges [][]*big.Int
	var curRange []*big.Int
	for _, ip := range ipsAsInt {
		if len(curRange) == 0 {
			curRange = append(curRange, ip)
		} else if big.NewInt(0).Sub(ip, curRange[len(curRange)-1]).Cmp(big.NewInt(1)) == 0 {
			curRange = append(curRange, ip)
		} else {
			ranges = append(ranges, curRange)
			curRange = []*big.Int{ip}
		}
	}
	if len(curRange) > 0 {
		ranges = append(ranges, curRange)
	}

	var collapse func([]*big.Int, []string) []string
	collapse = func(ipRange []*big.Int, cidrs []string) []string {
		// Alignment constraint (by trailing zeros):
		// A CIDR block of size 2^k must start at an address whose last k bits are zero.
		// The largest such block you can place at start is limited by the number of trailing zero bits in start.
		// If tz = trailingZeros(start), then alignment allows at most 2^tz addresses.
		start := ipRange[0]
		end := ipRange[len(ipRange)-1]
		tz := start.TrailingZeroBits()
		maxPerAlignment := int(math.Pow(2, float64(tz)))

		// Range-length constraint:
		// You also can’t make the block larger than what remains in the range.
		// Let length = end - start + 1. The largest power-of-two not exceeding length is 2^floor(log2(length)).
		length := float64(new(big.Int).Sub(end, start).Int64() + 1)
		maxPerLength := int(math.Pow(2, math.Floor(math.Log2(length))))

		// The maximum block if IP addresses
		max := min(maxPerAlignment, maxPerLength)

		// the prefix
		var prefix int
		if strings.Contains(FromInt(start).String(), ":") {
			// IPv6
			prefix = 128
		} else {
			// IPv4
			prefix = 32
		}
		prefix -= int(math.Log2(float64(max)))

		cidrs = append(cidrs, fmt.Sprintf("%s/%d", FromInt(start), prefix))
		if max < len(ipRange) {
			cidrs = collapse(ipRange[max:], cidrs)
		}
		return cidrs
	}

	var cidrs []string
	for _, r := range ranges {
		cidrs = collapse(r, cidrs)
	}
	return cidrs
}
