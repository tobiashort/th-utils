package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/tobiashort/th-utils/pkg/ip"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/th-utils/pkg/cidr"
)

type Args struct {
	File string `clap:"positional,desc='File with IP addresses or CIDRs separated by newlines. Reads from Stdin if not specified.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	var file *os.File
	if args.File != "" {
		var err error
		file, err = os.Open(args.File)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file error: %s: %v", args.File, err)
			os.Exit(1)
		}
	} else {
		file = os.Stdin
	}

	var ips []net.IP
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if strings.Contains(text, "/") {
			expanded, err := cidr.Expand(text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse error: %s: %v", text, err)
				os.Exit(1)
			}
			ips = append(ips, expanded...)
		} else {
			ip := net.ParseIP(text)
			if ip == nil {
				fmt.Fprintf(os.Stderr, "parse error: %s", text)
				os.Exit(1)
			}
			ips = append(ips, ip)
		}
	}
	err := scanner.Err()
	if err != nil {
		log.Fatalln()
		fmt.Fprintf(os.Stderr, "scanner error: %v", err)
		os.Exit(1)
	}

	cidrs := ip.Collapse(ips)
	for _, cidr := range cidrs {
		fmt.Println(cidr)
	}
}
