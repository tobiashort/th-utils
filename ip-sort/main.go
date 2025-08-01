package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"
)

func sortIPs(ip1, ip2 net.IP) int {
	for i := 0; i < 16; i++ {
		if ip1[i] != ip2[i] {
			return int(ip1[i]) - int(ip2[i])
		}
	}
	return 0
}

func main() {
	var reverse bool
	flag.BoolVar(&reverse, "r", false, "reverse")
	flag.Parse()
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
	if reverse {
		slices.Reverse(ips)
	}
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
