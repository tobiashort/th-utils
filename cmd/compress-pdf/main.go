package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/clog-go"
)

type Args struct {
	Debug   bool   `clap:"desc='Print debug output'"`
	OutFile string `clap:"desc='The name of the output file'"`
	InFile  string `clap:"mandatory,positional,desc='The PDF to compress'"`
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

	outFile := args.OutFile
	if outFile == "" {
		if i := strings.LastIndex(args.InFile, ".pdf"); i > 0 {
			outFile = args.InFile[0:i]
		}
		outFile += ".compressed.pdf"
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
	cmd.Args = append(cmd.Args, "-sOutputFile="+outFile)
	cmd.Args = append(cmd.Args, args.InFile)
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
