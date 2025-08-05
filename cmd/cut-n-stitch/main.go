package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func must(err error) {
  if err != nil {
    panic(err)
  }
}

func must2[T any](val T, err error) T {
  must(err)
  return val
}

func usage() {
  fmt.Fprintln(os.Stdout, "Usage: cutnstitch DELIMITER FORMAT")
  os.Exit(1)
}

func main() {
  flag.Parse()

  if flag.NArg() != 2 {
    usage()
  }
  
  delimiter := flag.Arg(0)
  format := must2(template.New("").Parse(fmt.Sprintf("%s\n", flag.Arg(1))))
  scanner := bufio.NewScanner(os.Stdin)

  for scanner.Scan() {
    line := scanner.Text()
    cut := strings.Split(line, delimiter)
    must(format.Execute(os.Stdout, cut))
  }
}
