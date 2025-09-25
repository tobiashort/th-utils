package main

import (
	"fmt"
	"os"

	"github.com/tobiashort/th-utils/pkg/cidr"

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
	ips := cidr.Expand(input)
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
