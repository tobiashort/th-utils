package main

import (
	"io"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/clog-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	Output string   `clap:"mandatory,desc='The name of the output PDF.'"`
	Files  []string `clap:"positional,mandatory,desc='The files to convert and include in the final PDF.'"`
	Debug  bool     `clap:"desc='Enable debug output'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	if args.Debug {
		clog.Level = clog.LevelDebug
	}

	externalTools := []string{
		"magick",
		"pandoc",
		"pdftk",
	}
	externalToolsOk := true
	for _, tool := range externalTools {
		if path, err := exec.LookPath(tool); err != nil {
			clog.Errorf("%s: not in path", tool)
		} else {
			clog.Infof("%s: %s", tool, path)
		}
	}
	if !externalToolsOk {
		os.Exit(1)
	}

	tempDir := must.Do2(os.MkdirTemp("", "tmp*"))
	clog.Debug("tmp:", tempDir)
	defer os.RemoveAll(tempDir)

	intermediatePdfs := make([]string, 0)
	for _, file := range args.Files {
		intermediatePdf := filepath.Join(tempDir, filepath.Base(file)+".pdf")
		intermediatePdfs = append(intermediatePdfs, intermediatePdf)
		isImage := false
		isPDF := false
		ext := filepath.Ext(file)
		if ext != "" {
			mimeType := mime.TypeByExtension(ext)
			if strings.HasPrefix(mimeType, "image/") {
				isImage = true
			} else if mimeType == "application/pdf" {
				isPDF = true
			}
		}
		if isImage {
			clog.Info("using", file)
			cmd := exec.Command("magick", "convert", file, intermediatePdf)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			must.Do(cmd.Run())
		} else if isPDF {
			clog.Info("using", file)
			fileReader := must.Do2(os.Open(file))
			defer fileReader.Close()
			intermediatePdfWriter := must.Do2(os.Create(intermediatePdf))
			defer intermediatePdfWriter.Close()
			must.Do2(io.Copy(intermediatePdfWriter, fileReader))
		} else {
			clog.Info("using", file)
			cmd := exec.Command("pandoc", "-o", intermediatePdf, file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			must.Do(cmd.Run())
		}
	}

	cmd := exec.Command("pdftk")
	cmd.Args = append(cmd.Args, intermediatePdfs...)
	cmd.Args = append(cmd.Args, "output", args.Output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must.Do(cmd.Run())

	clog.Info("created", args.Output)
}
