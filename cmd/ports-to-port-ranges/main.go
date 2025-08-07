package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tobiashort/clap-go"
)

func groupInts(n []int) [][]int {
	groups := make([][]int, 0)
	if len(n) == 0 {
		return groups
	}
	sort.Ints(n)
	group := make([]int, 0)
	for i, x := range n {
		if i == 0 || n[i]-n[i-1] == 1 {
			group = append(group, x)
		} else {
			groups = append(groups, group)
			group = []int{x}
		}
	}
	groups = append(groups, group)
	return groups
}

type Args struct {
	Ports string `clap:"positional,mandatory,description='Comma separated ports.'"`
}

func main() {
	args := Args{}
	clap.Example(`$ ports-to-port-ranges 1,2,3,4,5,6,11,223,445,555,556,557
1-6,11,223,445,555-557`)
	clap.Parse(&args)

	portStrings := strings.Split(args.Ports, ",")
	ports := make([]int, 0)
	for _, portString := range portStrings {
		port, err := strconv.Atoi(portString)
		if err != nil {
			panic(err)
		}
		ports = append(ports, port)
	}
	portRangeStrings := make([]string, 0)
	for _, group := range groupInts(ports) {
		if len(group) == 1 {
			portRangeStrings = append(portRangeStrings, fmt.Sprintf("%d", group[0]))
		} else if len(group) == 2 {
			portRangeStrings = append(portRangeStrings, fmt.Sprintf("%d", group[0]))
			portRangeStrings = append(portRangeStrings, fmt.Sprintf("%d", group[1]))
		} else {
			portRangeStrings = append(portRangeStrings, fmt.Sprintf("%d-%d", group[0], group[len(group)-1]))
		}
	}
	fmt.Println(strings.Join(portRangeStrings, ","))
}
