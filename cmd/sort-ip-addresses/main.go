package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/clap-go"
	. "github.com/tobiashort/utils-go/must"
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
	input := string(Must2(io.ReadAll(os.Stdin)))
	lines := strings.Split(input, "\n")
	slices.SortFunc(lines, func(line1, line2 string) int {
		split1 := strings.Split(line1, " ")
		split2 := strings.Split(line2, " ")
		ip1 := net.ParseIP(split1[0])
		if ip1 == nil {
			return 1
		}
		ip2 := net.ParseIP(split2[0])
		if ip2 == nil {
			return -1
		}
		return sortIPs(ip1, ip2)
	})
	if args.Reverse {
		slices.Reverse(lines)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
