package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tobiashort/cfmt-go"
	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/worker-go"
)

type Args struct {
	Prefix          string `clap:"default-value='th-',description='the prefix each binary will be given'"`
	InstallDir      string `clap:"default-value='$HOME/.th-utils/',description='install directory where tool are going to be installed'"`
	GenerateReadmes bool   `clap:"short=r,long=readmes,description='generates README.md for each tool'"`
}

func ensureDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
}

func main() {
	args := Args{}
	clap.Parse(&args)
	prefix := args.Prefix
	installDir := os.ExpandEnv(args.InstallDir)

	ensureDir("build")
	ensureDir(installDir)

	buildUtil := func(util string) error {
		cmd := exec.Command("go", "build", "-o", filepath.Join("build", prefix+util), "."+string(filepath.Separator)+util)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}

	installUtil := func(util string) error {
		cmd := exec.Command("cp", filepath.Join("build", prefix+util), installDir)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}

	generateReadme := func(util string) error {
		cmd := exec.Command("go", "run", "."+string(filepath.Separator)+util, "-h")
		bs, err := cmd.Output()
		if err != nil {
			return err
		}
		if len(bs) > 0 {
			bs = append([]byte{'`', '`', '`', '\n'}, bs...)
			bs = append(bs, '\n', '`', '`', '`', '\n')
			err = os.WriteFile(filepath.Join(".", util, "README.md"), bs, 0644)
			if err != nil {
				return err
			}
		}
		return nil
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	utils := make([]string, 0)
	for _, entry := range entries {
		match := entry.IsDir()
		match = match && entry.Name() != "build"
		match = match && entry.Name() != "vendor"
		match = match && !strings.HasPrefix(entry.Name(), ".")

		if match {
			utils = append(utils, entry.Name())
		}
	}

	errorSeen := false
	pool := worker.NewPool(5)
	for _, util := range utils {
		worker := pool.GetWorker()
		go func() {
			worker.Printf("#y{%s}", util)
			if args.GenerateReadmes {
				err := generateReadme(util)
				if err != nil {
					errorSeen = true
					worker.Logf("#r{ERROR} %s: %v", util, err)
				} else {
					worker.Logf("#g{SUCCESS} %s", util)
				}
			} else {
				err := buildUtil(util)
				if err != nil {
					errorSeen = true
					worker.Logf("#r{ERROR} %s: %v", util, err)
				} else {
					err := installUtil(util)
					if err != nil {
						errorSeen = true
						worker.Logf("#r{ERROR} %s: %v", util, err)
					} else {
						worker.Logf("#g{SUCCESS} %s", util)
					}
				}
			}
			worker.Done()
		}()
	}

	pool.Wait()

	if args.GenerateReadmes {
		generateReadme(".")
	}

	if errorSeen {
		cfmt.Printf("#r{--------------------------------------------------------------------------------}\n")
		cfmt.Printf("#r{ERROR}\n")
		cfmt.Printf("#r{--------------------------------------------------------------------------------}\n")
	} else {
		cfmt.Printf("#g{--------------------------------------------------------------------------------}\n")
		cfmt.Printf("#g{SUCCESS} %s\n", installDir)
		cfmt.Printf("#g{--------------------------------------------------------------------------------}\n")
	}
}
