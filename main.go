package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/tobiashort/cfmt-go"
	"github.com/tobiashort/worker-go"
)

var prefix = "th-"
var installDir = os.ExpandEnv("$HOME/.th-utils/")

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
	ensureDir("build")
	ensureDir(installDir)

	buildUtil := func(util string) error {
		cmd := exec.Command("go", "build", "-o", "build/"+prefix+util, "./"+util)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}

	installUtil := func(util string) error {
		cmd := exec.Command("cp", "build/"+prefix+util, installDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
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
			worker.Done()
		}()
	}

	pool.Wait()

	if errorSeen {
		cfmt.Printf("#r{---------------}\n")
		cfmt.Printf("#r{OVERALL FAILURE}\n")
		cfmt.Printf("#r{---------------}\n")
	} else {
		cfmt.Printf("#g{----------------%s}\n", strings.Repeat("-", len(installDir)))
		cfmt.Printf("#g{OVERALL SUCCESS} %s\n", installDir)
		cfmt.Printf("#g{----------------%s}\n", strings.Repeat("-", len(installDir)))
	}
}
