package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/groupby-go"
)

type Args struct {
	Editor string   `clap:"description='The path to the editor program to be opened'"`
	RgArgs []string `clap:"positional,description='Additional rg command line arguments'"`
}

func assertNil(val any) {
	if val != nil {
		panic(val)
	}
}

func assertNotNil(val any, format string, args ...any) {
	if val == nil {
		panic(fmt.Errorf(format, args...))
	}
}

func main() {
	var editor string

	args := Args{}
	clap.Parse(&args)

	editor = args.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		fmt.Fprintf(os.Stderr, "No editor configured. Use EDITOR environment variable or --editor argument.\n")
		os.Exit(1)
	}

	rgArgs := []string{"--line-number"}
	rgArgs = append(rgArgs, args.RgArgs...)
	cmd := exec.Command("rg", rgArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if string(out) == "" {
			fmt.Println("No matches.")
			os.Exit(1)
		}
		fmt.Fprint(os.Stderr, string(out))
		panic(err)
	}

	stateBefore := string(out)
	stateBefore = strings.TrimSpace(stateBefore)
	stateBeforeLines := strings.Split(stateBefore, "\n")

	tmp, err := os.CreateTemp("", "riplace")
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, bytes.NewBufferString(stateBefore))
	assertNil(err)

	tmp.Close()

	cmd = exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()

	data, err := os.ReadFile(tmp.Name())
	assertNil(err)

	stateAfter := string(data)
	stateAfter = strings.TrimSpace(stateAfter)
	stateAfterLines := strings.Split(stateAfter, "\n")

	if len(stateBeforeLines) != len(stateAfterLines) {
		panic(fmt.Errorf("line count missmatch, before %d, after %d", len(stateBeforeLines), len(stateAfterLines)))
	}

	type Change struct {
		File    string
		Line    int
		Content string
	}

	var changes = make([]Change, 0)

	pattern := regexp.MustCompile("^([^:]+):(\\d+):(.*)$")

	for stateBeforeLineIdx := range stateBeforeLines {
		stateBeforeLine := stateBeforeLines[stateBeforeLineIdx]
		stateBeforeLineMatches := pattern.FindStringSubmatch(stateBeforeLine)
		assertNotNil(stateBeforeLineMatches, "%s", stateBeforeLine)
		stateBeforeLineFile := stateBeforeLineMatches[1]
		stateBeforeLineNumber, err := strconv.Atoi(stateBeforeLineMatches[2])
		assertNil(err)
		stateBeforeLineContent := stateBeforeLineMatches[3]

		stateAfterLine := stateAfterLines[stateBeforeLineIdx]
		stateAfterLineMatches := pattern.FindStringSubmatch(stateAfterLine)
		assertNotNil(stateAfterLineMatches, "%s", stateAfterLine)
		stateAfterLineFile := stateAfterLineMatches[1]
		stateAfterLineNumber, err := strconv.Atoi(stateAfterLineMatches[2])
		assertNil(err)
		stateAfterLineContent := stateAfterLineMatches[3]

		if stateBeforeLineFile != stateAfterLineFile && stateBeforeLineNumber != stateAfterLineNumber {
			panic(fmt.Errorf("mismatch: %s, %s", stateBeforeLine, stateAfterLine))
		}

		if stateBeforeLineContent == stateAfterLineContent {
			continue
		}

		changes = append(changes, Change{
			File:    stateBeforeLineFile,
			Line:    stateBeforeLineNumber,
			Content: stateAfterLineContent,
		})
	}

	if len(changes) == 0 {
		fmt.Println("No changes.")
		os.Exit(0)
	}

	changesGroupedByFile := groupby.GroupBy(changes, func(a, b Change) bool { return a.File == b.File })

	tabwriter := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	fmt.Fprintf(tabwriter, "FILE\tLINES\n")
	for _, changes := range changesGroupedByFile {
		file := changes[0].File
		nLines := len(changes)
		fmt.Fprintf(tabwriter, "%s\t%d\n", file, nLines)
	}
	tabwriter.Flush()

ask:
	fmt.Print("Apply changes? (y/n) ")
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	assertNil(err)
	ans = strings.TrimSpace(ans)
	switch ans {
	case "n":
		fallthrough
	case "N":
		os.Exit(0)
	case "y":
		fallthrough
	case "Y":
		break
	default:
		goto ask
	}

	for _, changes := range changesGroupedByFile {
		file := changes[0].File
		data, err := os.ReadFile(file)
		assertNil(err)
		lines := strings.Split(string(data), "\n")
		for _, change := range changes {
			lines[change.Line-1] = change.Content
		}
		err = os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0644)
		assertNil(err)
	}
}
