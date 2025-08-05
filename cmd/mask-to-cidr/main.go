package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func printUsage() {
	fmt.Fprint(os.Stderr, `Usage: mask-to-cidr [MASK]
Read from STDIN when MASK is not provided as an argument.

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
	flag.Usage = printUsage
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
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(bytes))
	}
	split := strings.Split(input, ".")
	if len(split) != 4 {
		printInvalid(input)
		return
	}
	octets := [4]uint32{}
	for idx := 0; idx < 4; idx++ {
		octet, err := strconv.Atoi(split[idx])
		if err != nil || octet < 0 || octet > 255 {
			printInvalid(input)
			return
		}
		octets[idx] = uint32(octet)
	}
	mask := (octets[0] << 24) | (octets[1] << 16) | (octets[2] << 8) | octets[3]
	cidr := 0
	flipped := false
	for count := 0; count < 32; count++ {
		bit := (mask >> count) & 1
		if bit == 0 {
			if flipped {
				printInvalid(input)
				return
			}
			continue
		}
		flipped = true
		cidr++
	}
	fmt.Println(cidr)
}
