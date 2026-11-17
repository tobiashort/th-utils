package main

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	File string `clap:"positional,mandatory,desc='The file to preview'"`
}

func run() int {
	args := Args{}
	clap.Parse(&args)

	pth, err := exec.LookPath("pandoc")
	if err != nil {
		clog.Error("pandoc: #r{not in path}")
		return 1
	}
	clog.Infof("pandoc: #g{%s}", pth)

	tmp, err := os.CreateTemp("", "markdown-preview_*.html")
	if err != nil {
		clog.Error("create temp:", err)
		return 1
	}
	defer os.Remove(tmp.Name())

	var cmd *exec.Cmd

	cmd = exec.Command("pandoc", "-o", tmp.Name(), args.File)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return 1
	}

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", tmp.Name())
	case "windows":
		cmd = exec.Command("start", tmp.Name())
	case "darwin":
		cmd = exec.Command("open", tmp.Name())
	default:
		clog.Error("unknown GOOS:", runtime.GOOS)
		return 1
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return 1
	}

	must.Do2(bufio.NewReader(os.Stdin).ReadByte())

	return 0
}

func main() {
	os.Exit(run())
}
