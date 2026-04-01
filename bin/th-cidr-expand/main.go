package main

import (
	"fmt"
	"os"

	"github.com/tobiashort/th-utils/lib/cidr"

	"github.com/tobiashort/th-utils/lib/clap"
)

type Args struct {
	CIDR string `clap:"positional,mandatory,desc='E.g. 192.168.1.0/24'"`
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
