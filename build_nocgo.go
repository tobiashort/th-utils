//go:build !cgo

package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", ".")
	if len(os.Args) > 1 {
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
