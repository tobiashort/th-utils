package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: subnet-to-list [IP/CIDR]
Reads from STDIN if IP/CIDR is not defined as an argument

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func printInvalid(input string) {
	fmt.Fprintf(os.Stderr, "Invalid input '%s'\n", input)
	os.Exit(1)
}

func main() {
	help := flag.Bool("h", false, "print help")
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	if len(os.Args) > 2 {
		printUsage()
		return
	}
	input := ""
	if len(os.Args) == 2 {
		input = os.Args[1]
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(data))
	}
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
	for idx := 0; idx < 4; idx++ {
		octet, err := strconv.Atoi(humanReadableIpParts[idx])
		if err != nil || octet < 0 || octet > 255 {
			printInvalid(input)
			return
		}
		octets[idx] = uint32(octet)
	}
	mask := uint32(math.Pow(2, 32)-1) << (32 - cidr)
	ip := (octets[0]<<24 | octets[1]<<16 | octets[2]<<8 | octets[3]) & mask
	for count := 0; count < int(math.Pow(2, float64(32-cidr))); count++ {
		octet0 := uint8(ip >> 24)
		octet1 := uint8(ip >> 16)
		octet2 := uint8(ip >> 8)
		octet3 := uint8(ip)
		fmt.Printf("%d.%d.%d.%d\n", octet0, octet1, octet2, octet3)
		ip++
	}
}
