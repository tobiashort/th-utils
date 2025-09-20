package ip

import (
	"math/big"
	"net"
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
		return net.IP(bs).To4()
	} else {
		return net.IP(bs).To16()
	}
}
