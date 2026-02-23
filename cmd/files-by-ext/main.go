package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	WD string `clap:"positional,desc='The working directory'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	wd := args.WD
	if wd == "" {
		wd = must.Do2(os.Getwd())
	}

	byExt := make(map[string]int)
	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Type().IsRegular() {
			ext := filepath.Ext(path)
			ext = strings.ToLower(ext)
			byExt[ext]++
		}
		return nil
	})

	exts := make([]string, 0)
	for key := range byExt {
		exts = append(exts, key)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)

	slices.Sort(exts)
	for _, ext := range exts {
		count := byExt[ext]
		fmt.Fprintf(w, "%s\t%d\t\n", ext, count)
	}
	w.Flush()
}
