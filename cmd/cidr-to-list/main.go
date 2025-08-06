package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	CIDR string `clap:"positional,mandatory,description='E.g. 192.168.1.0/24'"`
}

func printInvalid(input string) {
	fmt.Fprintf(os.Stderr, "Invalid input '%s'\n", input)
	os.Exit(1)
}

func main() {
	args := Args{}
	clap.Parse(&args)

	input := args.CIDR

	split := strings.Split(input, "/")
	if len(split) != 2 {
		printInvalid(input)
		return
	}

	humanReadableIp := split[0]

	cidr, err := strconv.Atoi(split[1])
	if err != nil || cidr < 0 || cidr > 32 {
		printInvalid(input)
		return
	}

	humanReadableIpParts := strings.Split(humanReadableIp, ".")
	if len(humanReadableIpParts) != 4 {
		printInvalid(input)
		return
	}

	octets := [4]uint32{}
	for idx := range 4 {
		octet, err := strconv.Atoi(humanReadableIpParts[idx])
		if err != nil || octet < 0 || octet > 255 {
			printInvalid(input)
			return
		}
		octets[idx] = uint32(octet)
	}

	mask := uint32(math.Pow(2, 32)-1) << (32 - cidr)

	ip := (octets[0]<<24 | octets[1]<<16 | octets[2]<<8 | octets[3]) & mask

	for range int(math.Pow(2, float64(32-cidr))) {
		octet0 := uint8(ip >> 24)
		octet1 := uint8(ip >> 16)
		octet2 := uint8(ip >> 8)
		octet3 := uint8(ip)
		fmt.Printf("%d.%d.%d.%d\n", octet0, octet1, octet2, octet3)
		ip++
	}
}
