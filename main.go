package main

import (
	"os"
	"os/exec"
	"path/filepath"
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

func main() {
	args := Args{}
	clap.Parse(&args)
	prefix := args.Prefix

	if args.Clean {
		cleanUp("build")
	}

	ensureDir("build")

	buildUtil := func(util string) error {
		cmd := exec.Command("go", "build", "-o", filepath.Join("build", prefix+util), filepathJoinUncleaned(".", "cmd", util))
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}

	generateReadme := func(util string) error {
		cmd := exec.Command("go", "run", filepathJoinUncleaned(".", "cmd", util), "-h")
		bs, err := cmd.Output()
		if err != nil {
			return err
		}
		if len(bs) > 0 {
			bs = append([]byte{'`', '`', '`', '\n'}, bs...)
			bs = append(bs, '\n', '`', '`', '`', '\n')
			err = os.WriteFile(filepath.Join("cmd", util, "README.md"), bs, 0644)
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
