package main

import (
	"fmt"
	"os"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	CIDR int `clap:"positional,mandatory,description='The CIDR to be converted to subnet mask'"`
}

func printInvalid(input string) {
}

func main() {
	args := Args{}
	clap.Parse(&args)

	if args.CIDR < 0 || args.CIDR > 32 {
		fmt.Fprintf(os.Stderr, "Invalid input '%d'. Must be inbetween 0 and 32.\n", args.CIDR)
		os.Exit(1)
	}

	mask := uint32(0b11111111_11111111_11111111_11111111)
	mask = (mask << (32 - args.CIDR)) & mask
	octet0 := uint8(mask >> 24)
	octet1 := uint8(mask >> 16)
	octet2 := uint8(mask >> 8)
	octet3 := uint8(mask)

	fmt.Printf("%d.%d.%d.%d\n", octet0, octet1, octet2, octet3)
}
