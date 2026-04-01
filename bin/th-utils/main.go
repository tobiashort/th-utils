package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/tobiashort/th-utils/lib/choose"
	"github.com/tobiashort/th-utils/lib/clap"
)

var Utils string

type Args struct{}

func main() {
	args := Args{}
	clap.Parse(&args)
	if option, ok := choose.One("Utils:", choose.ToOptions(strings.Split(Utils, ","))); ok {
		cmd := exec.Command(option.Value, "-h")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			os.Exit(0)
		}
	}
	os.Exit(1)
}
