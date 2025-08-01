package main

import (
	"fmt"
	"os"
	"os/exec"
)

var utils = []string{
	"append",
}

func main() {
	args := os.Args
	if len(args) <= 1 || len(args) > 2 {
		help()
		os.Exit(1)
	}

	switch args[1] {
	case "help":
		help()
		os.Exit(0)
	case "build":
		err := build()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "clean":
		clean()
		os.Exit(0)
	}
}

func build() error {
	os.MkdirAll("build", 0755)

	for _, u := range utils {
		cmd := exec.Command("go", "build", "-o", "build/"+u, "./"+u)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
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
