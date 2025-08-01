package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2[T any](v T, err error) T {
	must(err)
	return v
}

func main() {
	var wd string

	flag.StringVar(&wd, "wd", "", "working directory")
	flag.Parse()

	if wd == "" {
		wd = must2(os.Getwd())
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
		fmt.Fprintf(w, fmt.Sprintf("%s\t%d\t\n", ext, count))
	}
	w.Flush()
}
