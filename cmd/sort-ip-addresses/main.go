package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/tobiashort/th-utils/pkg/ip"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	Reverse bool `clap:"description=Reverses sort order"`
}

func main() {
	args := Args{}
	clap.Parse(&args)
	input := string(must.Do2(io.ReadAll(os.Stdin)))
	lines := strings.Split(input, "\n")
	slices.SortFunc(lines, func(line1, line2 string) int {
		trim1 := strings.TrimSpace(line1)
		trim2 := strings.TrimSpace(line2)
		split1 := strings.Split(trim1, " ")
		split2 := strings.Split(trim2, " ")
		ip1 := net.ParseIP(split1[0])
		ip2 := net.ParseIP(split2[0])
		return ip.ToInt(ip1).Cmp(ip.ToInt(ip2))
	})
	if args.Reverse {
		slices.Reverse(lines)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
