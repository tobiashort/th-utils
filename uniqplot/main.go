package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var debug bool
var fillChar string

func columns() int {
	var columns int
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdout
	cmdOut, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	cmdOutStr := string(cmdOut)
	cmdOutStr = strings.TrimSpace(cmdOutStr)
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "exit code: %d: %s", exitCode, cmdOutStr)
		}
		columns = 80
	} else {
		columnsStr := strings.Split(cmdOutStr, " ")[1]
		columns, err = strconv.Atoi(columnsStr)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			columns = 80
		}
	}
	if columns > 80 {
		return 80
	}
	return columns
}

func maxLabelWidth(labels []string) int {
	maxLabelWidth := 0
	for _, label := range labels {
		labelWidth := len(label)
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}
	return maxLabelWidth
}

func maxValue(values []int) int {
	maxValue := 0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func maxValueWidth(values []int) int {
	maxValueWidth := 0
	for _, value := range values {
		valueWidth := len(strconv.Itoa(value))
		if valueWidth > maxValueWidth {
			maxValueWidth = valueWidth
		}
	}
	return maxValueWidth
}

func maxBarWidth(labels []string, values []int) int {
	columns := columns()
	maxLabelWidth := maxLabelWidth(labels)
	maxValueWidth := maxValueWidth(values)
	padding := 5
	// label : ======================== val
	//      ^^^                        ^   ^
	return columns - maxLabelWidth - maxValueWidth - padding
}

func plot(labels []string, values []int) {
	maxLabelWidth := maxLabelWidth(labels)
	maxBarWidth := maxBarWidth(labels, values)
	maxValue := maxValue(values)
	format := fmt.Sprintf("%%%ds : %%s %%d \n", maxLabelWidth)
	for index, label := range labels {
		value := values[index]
		barWidth := int((float64(value) / float64(maxValue)) * float64(maxBarWidth))
		bar := strings.Repeat(fillChar, barWidth)
		fmt.Printf(format, label, bar, value)
	}
}

func main() {
	flag.BoolVar(&debug, "debug", false, "print debug output")
	flag.StringVar(&fillChar, "fillChar", "=", "the fill character using to draw the bar")
	flag.Parse()
	labels := make([]string, 0)
	values := make([]int, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		split := strings.Split(line, " ")
		valueStr := split[0]
		label := split[1]
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			panic(err)
		}
		labels = append(labels, label)
		values = append(values, value)
	}
	plot(labels, values)
}
