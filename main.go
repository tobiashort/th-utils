package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/worker-go"
)

type Args struct {
	Prefix string `clap:"default-value='th-',description='the prefix each binary will be given'"`
	Util   string `clap:"positional,description='only compiles and installes the given utitliy'"`
	Clean  bool   `clap:"description='delete installation path'"`
}

func cleanUp(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}
}

func ensureDir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
}

func filepathJoinUncleaned(parts ...string) string {
	return strings.Join(parts, string(filepath.Separator))
}

func runTests() {
	fmt.Println("Running tests...")
	cmd := exec.Command("go", "test", "./...")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		os.Exit(1)
	} else {
		fmt.Println("Tests ok.")
	}
}

func main() {
	args := Args{}
	clap.Parse(&args)
	prefix := args.Prefix

	if args.Clean {
		cleanUp("build")
	}

	ensureDir("build")

	runTests()

	buildUtil := func(util string) error {
		executable := filepath.Join("build", prefix+util)
		if runtime.GOOS == "windows" {
			executable += ".exe"
		}
		cmd := exec.Command("go", "build", "-o", executable, filepathJoinUncleaned(".", "cmd", util))
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

	generateReadme := func(util string) error {
		cmd := exec.Command("go", "run", filepathJoinUncleaned(".", "cmd", util), "-h")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%w: %s", err, string(out))
		}
		if len(out) > 0 {
			out = append([]byte{'`', '`', '`', '\n'}, out...)
			out = append(out, '\n', '`', '`', '`', '\n')
			err = os.WriteFile(filepath.Join("cmd", util, "README.md"), out, 0644)
			if err != nil {
				return err
			}
		}
		return nil
	}

	utils := make([]string, 0)

	if args.Util != "" {
		utils = append(utils, args.Util)
	} else {
		entries, err := os.ReadDir("cmd")
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				utils = append(utils, entry.Name())
			}
		}
	}

	fmt.Println("Building utils...")

	pool := worker.NewPool(5)
	for _, util := range utils {
		worker := pool.GetWorker()
		worker.Go(
			func() {
				worker.Printf("#y{%s}", util)
				err := buildUtil(util)
				if err != nil {
					worker.Logf("#r{ERROR} %s: %v", util, err)
					return
				}
				err = generateReadme(util)
				if err != nil {
					worker.Logf("#r{ERROR} %s: %v", util, err)
					return
				}
				worker.Logf("#g{SUCCESS} %s", util)
			})
	}
	pool.Wait()

	generateReadme(".")
}
