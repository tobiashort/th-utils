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

func main() {
	args := Args{}
	clap.Parse(&args)

	input := args.CIDR
	ips, err := cidr.Expand(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %s: %s\n", input, err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
