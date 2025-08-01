package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: ports-to-port-ranges COMMA_SEPERATED_PORTS")
		os.Exit(1)
	}
	commaSeperatedPorts := os.Args[1]
	portStrings := strings.Split(commaSeperatedPorts, ",")
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