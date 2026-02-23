package main

import (
	"bufio"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tobiashort/cfmt-go"
	"github.com/tobiashort/choose-go"
	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/clog-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	Dir       string `clap:"positional,desc='Optional directory otherwise the current working directory is used'"`
	Editor    string `clap:"desc='The path to the editor program to be opened'"`
	Recursive bool   `clap:"desc='Recursively list all files.'"`
}

type Action interface {
	Preview() string
	Execute() error
}

type MoveAction struct {
	Path    string
	NewPath string
}

func (a MoveAction) Preview() string {
	return cfmt.Sprintf("#y{%s -> %s}", a.Path, a.NewPath)
}

func (a MoveAction) Execute() error {
	err := os.MkdirAll(filepath.Dir(a.NewPath), 0755)
	if err != nil {
		return err
	}
	return os.Rename(a.Path, a.NewPath)
}

type DeleteAction struct {
	Path string
}

func (a DeleteAction) Preview() string {
	return cfmt.Sprintf("#r{%s}", a.Path)
}

func (a DeleteAction) Execute() error {
	return os.Remove(a.Path)
}

func main() {
	args := Args{}
	clap.Parse(&args)

	dir := args.Dir
	if dir == "" {
		dir = must.Do2(os.Getwd())
	}

	editor := args.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		clog.Error("No editor configured. Use EDITOR environment variable or --editor argument")
		os.Exit(1)
	}

	paths := make([]string, 0)
	if args.Recursive {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			paths = append(paths, path)
			return nil
		})
	} else {
		entries := must.Do2(os.ReadDir(dir))
		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())
			paths = append(paths, path)
		}
	}

	temp := must.Do2(os.CreateTemp("", "tmp*"))
	defer os.Remove(temp.Name())
	for _, path := range paths {
		temp.WriteString(path + "\n")
	}
	must.Do(temp.Close())

	cmd := exec.Command(editor, temp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must.Do(cmd.Run())

	newPaths := make([]string, 0)

	temp = must.Do2(os.Open(temp.Name()))
	scanner := bufio.NewScanner(temp)
	for scanner.Scan() {
		newPaths = append(newPaths, scanner.Text())
	}
	must.Do(temp.Close())

	if len(paths) != len(newPaths) {
		clog.Errors("Expected same amount of lines.")
		os.Exit(1)
	}

	actions := make([]Action, 0)

	for i := 0; i < len(paths); i++ {
		path := paths[i]
		newPath := newPaths[i]
		if path != newPath {
			if newPath == "" {
				actions = append(actions, DeleteAction{Path: path})
			} else {
				actions = append(actions, MoveAction{Path: path, NewPath: newPath})
			}
		}
	}

	if len(actions) == 0 {
		clog.Info("Nothing to do.")
		os.Exit(0)
	}

	for _, action := range actions {
		cfmt.Println(action.Preview())
	}

	if choose.YesNo("Proceed?", choose.DEFAULT_NO) {
		for _, action := range actions {
			err := action.Execute()
			if err != nil {
				clog.Error(err.Error())
			}
		}
	}
}
