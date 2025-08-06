package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Mask string `clap:"positional,mandatory,description='The mask to convert to CIDR'"`
}

func printInvalid(mask string) {
	fmt.Fprintf(os.Stderr, "Invalid mask '%s'\n", mask)
	os.Exit(1)
}

func main() {
	args := Args{}
	clap.Parse(&args)

	split := strings.Split(args.Mask, ".")
	if len(split) != 4 {
		printInvalid(args.Mask)
		return
	}
	octets := [4]uint32{}
	for idx := range 4 {
		octet, err := strconv.Atoi(split[idx])
		if err != nil || octet < 0 || octet > 255 {
			printInvalid(args.Mask)
			return
		}
		octets[idx] = uint32(octet)
	}
	mask := (octets[0] << 24) | (octets[1] << 16) | (octets[2] << 8) | octets[3]
	cidr := 0
	flipped := false
	for count := range 32 {
		bit := (mask >> count) & 1
		if bit == 0 {
			if flipped {
				printInvalid(args.Mask)
				return
			}
			continue
		}
		flipped = true
		cidr++
	}
	fmt.Println(cidr)
}
