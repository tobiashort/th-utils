//go:build cgo

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tobiashort/th-utils/lib/cfmt"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/clog"
	"github.com/tobiashort/th-utils/lib/must"
	strings2 "github.com/tobiashort/th-utils/lib/strings"
	"github.com/tobiashort/th-utils/lib/worker"
)

type BuildOpt struct {
	Util string `clap:"positional,desc='Only builds the given utitliy.'"`
}

type Args struct {
	Command any      `clap:"cmd,desc='The command to run.'"`
	Clean   any      `clap:"cmdopt,desc='Deletes build path.'"`
	Test    any      `clap:"cmdopt,desc='Runs all tests'"`
	Build   BuildOpt `clap:"cmdopt,desc='Builds binaries.'"`
	Repl    BuildOpt `clap:"cmdopt,desc='Starts REPL.'"`
}

func cleanUp(dir string) {
	must.Do(os.RemoveAll(dir))
}

func ensureDir(dir string) {
	must.Do(os.MkdirAll(dir, 0755))
}

func filepathJoinUncleaned(parts ...string) string {
	return strings.Join(parts, string(filepath.Separator))
}

func listDirs(dir string) []string {
	dirs := make([]string, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs
}

func listBins() []string {
	return listDirs("bin")
}

func listLibs() []string {
	return listDirs("lib")
}

func runClean() {
	cleanUp("build")
}

func runTest() bool {
	pass := true

	testPath := func(pool worker.Pool, pth string) {
		worker := pool.GetWorker()
		worker.Go(func() {
			worker.Printf("#y{%s}", pth)
			cmd := exec.Command("go", "test", pth)
			if err := cmd.Run(); err == nil {
				worker.Logf("#g{PASS} %s", pth)
			} else {
				worker.Logf("#r{FAIL} %s", pth)
				pass = false
			}
		})
	}

	cfmt.Printf("#b{[test %s]}\n", filepathJoinUncleaned(".", "lib", "..."))
	strings.Join([]string{".", "lib", "..."}, string(os.PathSeparator))
	libs := listLibs()
	pool := worker.NewPool(min(len(libs), 5))
	for _, lib := range libs {
		testPath(pool, filepathJoinUncleaned(".", "lib", lib, "..."))
	}
	pool.Wait()

	cfmt.Printf("#b{[test %s]}\n", filepathJoinUncleaned(".", "bin", "..."))
	bins := listBins()
	pool = worker.NewPool(min(len(bins), 5))
	for _, bin := range bins {
		testPath(pool, filepathJoinUncleaned(".", "bin", bin, "..."))
	}
	pool.Wait()

	return pass
}

func runBuild(opt BuildOpt) bool {
	success := true

	ensureDir("build")

	var bins []string
	if opt.Util != "" {
		bins = []string{opt.Util}
	} else {
		bins = listBins()
	}

	buildUtil := func(util string) error {
		executable := filepath.Join("build", util)
		if runtime.GOOS == "windows" {
			executable += ".exe"
		}
		cmd := exec.Command("go", "build")
		if util == "th-utils" {
			cmd.Args = append(cmd.Args, "-ldflags", "-X main.Utils="+strings.Join(bins, ","))
		}
		cmd.Args = append(cmd.Args, "-o", executable)
		cmd.Args = append(cmd.Args, filepathJoinUncleaned(".", "bin", util))
		cmd.Env = os.Environ()
		if runtime.GOOS == "windows" {
			cmd.Env = append(cmd.Env, "CC=zig cc")
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%w: %s", err, string(out))
		}
		return nil
	}

	generateReadmePath := func(pth string) error {
		cmd := exec.Command("go", "run", pth, "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%w: %s", err, string(out))
		}
		if len(out) > 0 {
			out = append([]byte{'`', '`', '`', '\n'}, out...)
			out = append(out, '\n', '`', '`', '`', '\n')
			err = os.WriteFile(filepath.Join(pth, "README.md"), out, 0644)
			if err != nil {
				return err
			}
		}
		return nil
	}

	generateReadmeUtil := func(util string) error {
		return generateReadmePath(filepathJoinUncleaned(".", "bin", util))
	}

	cfmt.Println("#b{[build]}")

	pool := worker.NewPool(min(len(bins), 5))
	for _, util := range bins {
		worker := pool.GetWorker()
		worker.Go(
			func() {
				worker.Printf("#y{%s}", util)
				err := buildUtil(util)
				if err != nil {
					msg := cfmt.Sprintf("#r{ERROR} %s: %v", util, err)
					worker.Logf(msg)
					success = false
					return
				}
				err = generateReadmeUtil(util)
				if err != nil {
					msg := cfmt.Sprintf("#r{ERROR} %s: %v", util, err)
					worker.Logf(msg)
					success = false
					return
				}
				worker.Logf("#g{SUCCESS} %s", util)
			})
	}
	pool.Wait()

	generateReadmePath(".")

	return success
}

func main() {
	args := Args{}
	clap.Parse(&args)

	switch args.Command {
	case nil:
		opt := BuildOpt{}
		clap.Parse(&opt)
		testOk := runTest()
		if !testOk {
			cfmt.Println("#y{====}")
			cfmt.Println("#r{ERROR}")
			os.Exit(1)
		}
		buildOk := runBuild(opt)
		if buildOk {
			cfmt.Println("#y{====}")
			cfmt.Println("#g{SUCCESS}")
			os.Exit(0)
		} else {
			cfmt.Println("#y{====}")
			cfmt.Println("#r{ERROR}")
			os.Exit(1)
		}
	case &args.Clean:
		runClean()
	case &args.Test:
		if ok := runTest(); ok {
			cfmt.Println("#y{====}")
			cfmt.Println("#g{PASS}")
			os.Exit(0)
		} else {
			cfmt.Println("#y{====}")
			cfmt.Println("#r{FAIL}")
			os.Exit(1)
		}
	case &args.Build:
		if ok := runBuild(args.Build); ok {
			cfmt.Println("#y{====}")
			cfmt.Println("#g{SUCCESS}")
			os.Exit(0)
		} else {
			cfmt.Println("#y{====}")
			cfmt.Println("#r{ERROR}")
			os.Exit(1)
		}
	case &args.Repl:
		if pth, err := exec.LookPath("gomacro"); err == nil {
			clog.Infof("gomacro: %s", pth)
			cmd := exec.Command("gomacro")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			must.Do(cmd.Run())
		} else {
			clog.Error(strings2.Dedent(`gomacro: #r{not found}
				                       |go install github.com/cosmos72/gomacro@latest`))
		}
	default:
		clog.Errorf("Unknown command: %v", args.Command)
		os.Exit(1)
		return
	}
}
