package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/clap-go"
)

type Args struct {
	Reverse bool `clap:"description=Reverses sort order"`
}

func sortIPs(ip1, ip2 net.IP) int {
	for i := range 16 {
		if ip1[i] != ip2[i] {
			return int(ip1[i]) - int(ip2[i])
		}
	}
	return 0
}

func main() {
	args := Args{}
	clap.Parse(&args)
	ips := make([]net.IP, 0)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		line = strings.TrimSpace(line)
		ip := net.ParseIP(line)
		if ip == nil {
			fmt.Fprintf(os.Stderr, "parse error: '%s'", line)
		} else {
			ips = append(ips, ip)
		}
	}
	slices.SortFunc(ips, sortIPs)
	if args.Reverse {
		slices.Reverse(ips)
	}
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
