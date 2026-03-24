package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	Out  string `clap:"desc='The output file'"`
	File string `clap:"mandatory,positional,desc='The video to convert'"`
}

func run() int {
	args := Args{}
	clap.Parse(&args)

	if ffmpeg, err := exec.LookPath("ffmpeg"); err != nil {
		clog.Error("gs: #r{not found}")
		return 1
	} else {
		clog.Infof("ffmpeg: #g{%s}", ffmpeg)
	}

	out := args.Out
	if out == "" {
		ext := filepath.Ext(args.File)
		if i := strings.LastIndex(args.File, ext); i > 0 {
			out = args.File[0:i]
		}
		out += ".gif"
	}

	tmp := must.Do2(os.MkdirTemp("", "video-to-gif_*"))
	defer os.RemoveAll(tmp)

	cmd := exec.Command("ffmpeg")
	cmd.Args = append(cmd.Args, "-y")
	cmd.Args = append(cmd.Args, "-i", args.File)
	cmd.Args = append(cmd.Args, "-vf", "palettegen")
	cmd.Args = append(cmd.Args, filepath.Join(tmp, "palette.png"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 1
	}

	cmd = exec.Command("ffmpeg")
	cmd.Args = append(cmd.Args, "-y")
	cmd.Args = append(cmd.Args, "-i", args.File)
	cmd.Args = append(cmd.Args, "-i", filepath.Join(tmp, "palette.png"))
	cmd.Args = append(cmd.Args, "-filter_complex", "paletteuse")
	cmd.Args = append(cmd.Args, "-r", "10")
	cmd.Args = append(cmd.Args, out)
	cmd.Stdin = os.Stdin
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
