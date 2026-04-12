package main

import (
	"encoding/hex"
	"net"
	"strings"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	Mac   string `clap:"mandatory,positional,desc='The MAC address of the device'"`
	Debug bool   `clap:"desc='Print debug information'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	if args.Debug {
		clog.Level = clog.LevelDebug
	}

	macStr := args.Mac
	macStr = strings.ReplaceAll(macStr, "-", ":")
	macStr = strings.ToLower(macStr)

	macSplit := strings.Split(macStr, ":")
	assert.True(len(macSplit) == 6, "MAC must be 6 bytes")

	mac := make([]byte, 6)
	for i := 0; i < 6; i++ {
		macPart := macSplit[i]
		if len(macPart) == 1 {
			macPart = "0" + macPart
		}
		assert.True(len(macPart) == 2, "macPart invalid length")
		copy(mac[i:i+1], must.Do2(hex.DecodeString(macPart)))
	}

	packet := make([]byte, 6+16*len(mac))
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 6; i < len(packet); i += 6 {
		copy(packet[i:i+6], mac)
	}

	clog.Debugf("Mac: % x", mac)
	clog.Debugf("Packet: % x", packet)

	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.To4() == nil {
			continue
		}

		ip := ipnet.IP.To4()
		mask := ipnet.Mask

		broadcast := net.IPv4(
			ip[0]|^mask[0],
			ip[1]|^mask[1],
			ip[2]|^mask[2],
			ip[3]|^mask[3],
		)

		udpAddr := net.UDPAddr{
			IP:   broadcast,
			Port: 9,
		}

		conn, err := net.DialUDP("udp", nil, &udpAddr)
		if err != nil {
			clog.Error(err)
			continue
		}
		defer conn.Close()

		clog.Info("Broadcast:", broadcast)
		conn.Write(packet)
	}
}
