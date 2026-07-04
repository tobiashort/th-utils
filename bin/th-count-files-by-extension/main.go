package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"text/tabwriter"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
	Path string `clap:"positional,desc='A file system directory path, or a file containg paths, or reads from Stdin if not specified'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	paths := []string{}

	walkDir := func(p string) {
		filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && d.Type().IsRegular() {
				paths = append(paths, p)
			}
			return nil
		})
	}

	readFile := func(f *os.File) {
		s := bufio.NewScanner(f)
		for s.Scan() {
			paths = append(paths, s.Text())
		}
		assert.Nil(s.Err(), "scanner error")
	}

	if args.Path != "" {
		stat := must.Do2(os.Stat(args.Path))
		if stat.IsDir() {
			walkDir(args.Path)
		} else {
			f := must.Do2(os.Open(args.Path))
			readFile(f)
		}
	} else {
		readFile(os.Stdin)
	}

	byExt := make(map[string]int)

	for _, p := range paths {
		e := filepath.Ext(p)
		byExt[e] = byExt[e] + 1
	}

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
