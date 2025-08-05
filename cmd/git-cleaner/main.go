package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tobiashort/worker-go"
)

var flagRemoveLocalBranches bool

type ExecutionResult struct {
	path   string
	output string
	err    error
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2[T any](v T, err error) T {
	must(err)
	return v
}

func findGitRepositories() []string {
	var paths []string
	wd := must2(os.Getwd())
	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == ".git" {
			paths = append(paths, filepath.Dir(path))
		}
		return nil
	})
	return paths
}

func run(cmd *exec.Cmd) (output string, err error) {
	bs, err := cmd.CombinedOutput()
	output = string(bs)
	if err != nil {
		return output, err
	}
	return output, nil
}

func runGit(path string, args ...string) ExecutionResult {
	args = append([]string{"-C", path}, args...)
	cmd := exec.Command("git", args...)
	output, err := run(cmd)
	return ExecutionResult{path, output, err}
}

func gitResetHard(path string) ExecutionResult {
	return runGit(path, "reset", "--hard", "--recurse-submodules")
}

func gitListBranches(path string) (branches []string, executionResult ExecutionResult) {
	executionResult = runGit(path, "branch", "--no-color")
	if executionResult.err != nil {
		return nil, executionResult
	}
	for line := range strings.SplitSeq(executionResult.output, "\n") {
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "* ", "", 1)
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches, executionResult
}

func gitCheckoutMain(path string) ExecutionResult {
	branches, executionResult := gitListBranches(path)
	if executionResult.err != nil {
		return executionResult
	}
	var branch string
	if slices.Contains(branches, "master") {
		branch = "master"
	} else if slices.Contains(branches, "main") {
		branch = "main"
	}
	if branch == "" {
		return ExecutionResult{
			executionResult.path,
			executionResult.output,
			fmt.Errorf("no master/main branch found"),
		}
	}
	return runGit(path, "checkout", branch)
}

func gitClean(path string) ExecutionResult {
	return runGit(path, "clean", "-fd")
}

func gitPull(path string) ExecutionResult {
	return runGit(path, "pull", "-p")
}

func gitRemoveLocalBranches(path string) ExecutionResult {
	branches, executionResult := gitListBranches(path)
	if executionResult.err != nil {
		return executionResult
	}
	for _, branch := range branches {
		if branch == "master" || branch == "main" {
			continue
		}
		executionResult = runGit(path, "branch", "-D", branch)
		if executionResult.err != nil {
			return executionResult
		}
	}
	return ExecutionResult{path, "", nil}
}

func cleanGitRepository(path string, worker worker.Worker) {
	var executionResult ExecutionResult
	worker.Printf(path)
	executionResult = gitCheckoutMain(path)
	if executionResult.err != nil {
		goto errorCase
	}
	executionResult = gitResetHard(path)
	if executionResult.err != nil {
		goto errorCase
	}
	executionResult = gitClean(path)
	if executionResult.err != nil {
		goto errorCase
	}
	executionResult = gitPull(path)
	if executionResult.err != nil {
		goto errorCase
	}
	if flagRemoveLocalBranches {
		executionResult = gitRemoveLocalBranches(path)
		if executionResult.err != nil {
			goto errorCase
		}
	}
	worker.Logf("[#g{DONE}] %s\n", path)
	worker.Done()
	return
errorCase:
	worker.Logf("[#r{ERROR}] %s\n", path)
	worker.Logf("%s", executionResult.err)
	worker.Logf(executionResult.output)
	worker.Done()
	return
}

func main() {
	flag.BoolVar(&flagRemoveLocalBranches, "remove-local-branches", false, "removes all local branches")
	flag.Parse()
	gitRepositories := findGitRepositories()
	pool := worker.NewPool(5)
	for _, path := range gitRepositories {
		worker := pool.GetWorker()
		go cleanGitRepository(path, worker)
	}
	pool.Wait()
}
