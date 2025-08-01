package main

import (
	"fmt"
	"os"
	"os/exec"
)

var utils = []string{
	"append",
	"prepend",
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		help()
		os.Exit(1)
	}

	switch args[1] {
	case "help":
		help()
		os.Exit(0)
	case "build":
		err := build(args[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "clean":
		clean()
		os.Exit(0)
	}
}

var prefix = "th-"

func build(args []string) error {
	os.MkdirAll("build", 0755)

	buildUtil := func(util string) error {
		cmd := exec.Command("go", "build", "-o", "build/"+prefix+util, "./"+util)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}

	if len(args) > 1 {
		util := args[1]
		err := buildUtil(util)
		if err != nil {
			return err
		}
		return nil
	} else {
		for _, util := range utils {
			err := buildUtil(util)
			if err != nil {
				return nil
			}
		}
		return nil
	}
}

func clean() {
	os.RemoveAll("build")
}

func help() {
	fmt.Print(`usage: th-utils CMD

CMD:
  build       - build all utils
  clean       - remove generated files
  help        - print help
`)
}
