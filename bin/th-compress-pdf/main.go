package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
)

type Args struct {
	Debug bool   `clap:"desc='Print debug output'"`
	Out   string `clap:"desc='The name of the output file'"`
	File  string `clap:"mandatory,positional,desc='The PDF to compress'"`
}

func run() int {
	args := Args{}
	clap.Parse(&args)

	if gs, err := exec.LookPath("gs"); err != nil {
		clog.Error("gs: #r{not found}")
		return 1
	} else {
		clog.Infof("gs: #g{%s}", gs)
	}

	out := args.Out
	if out == "" {
		if i := strings.LastIndex(args.File, ".pdf"); i > 0 {
			out = args.File[0:i]
		}
		out += ".compressed.pdf"
	}

	cmd := exec.Command("gs")
	cmd.Args = append(cmd.Args, "-sDEVICE=pdfwrite")
	cmd.Args = append(cmd.Args, "-dCompatibilityLevel=1.4")
	// -dPDFSETTINGS=/screen   — Low quality and small size at 72dpi.
	// -dPDFSETTINGS=/ebook    — Slightly better quality but also a larger file size at 150dpi.
	// -dPDFSETTINGS=/prepress — High quality and large size at 300 dpi.
	// -dPDFSETTINGS=/default  — System chooses the best output, which can create larger PDF files.
	cmd.Args = append(cmd.Args, "-dPDFSETTINGS=/screen")
	cmd.Args = append(cmd.Args, "-dNOPAUSE")
	cmd.Args = append(cmd.Args, "-dQUIET")
	cmd.Args = append(cmd.Args, "-dBATCH")
	if args.Debug {
		cmd.Args = append(cmd.Args, "-dDEBUG")
	}
	cmd.Args = append(cmd.Args, "-sOutputFile="+out)
	cmd.Args = append(cmd.Args, args.File)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
