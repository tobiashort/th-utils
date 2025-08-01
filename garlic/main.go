package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func assertNil(val any) {
	if val != nil {
		panic(val)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: garlic [dir]
 
ARGUMENTS
            [dir]   optional directory otherwise the current working directory
                    is used.
 
OPTIONS
`)
	flag.PrintDefaults()
}

func main() {
	var dir string
	var editor string

	flag.Usage = usage
	flag.StringVar(&editor, "editor", os.Getenv("EDITOR"), "editor to be used")
	flag.Parse()

	if flag.NArg() == 0 {
		var err error
		dir, err = os.Getwd()
		assertNil(err)
	} else if flag.NArg() == 1 {
		dir = flag.Arg(0)
	} else {
		usage()
		os.Exit(1)
	}

	if editor == "" {
		fmt.Fprintf(os.Stderr, "No editor configured. Use EDITOR environment variable or -editor flag.\n")
		os.Exit(1)
	}

	entries, err := os.ReadDir(dir)
	assertNil(err)

	lineFormat := fmt.Sprintf("[%%%dd] %%s", len(strconv.Itoa(len(entries))))
	linePattern := regexp.MustCompile("^(\\[\\s*[0-9]+\\]\\s)(.*)$")

	tempFile, err := os.CreateTemp("", "garlic")
	assertNil(err)
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	linesBefore := make([]string, 0)

	for index, entry := range entries {
		line := fmt.Sprintf(lineFormat, index+1, fmt.Sprint(filepath.Join(dir, entry.Name())))
		linesBefore = append(linesBefore, line)
	}

	_, err = tempFile.WriteString(strings.Join(linesBefore, "\n"))
	assertNil(err)

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()

	tempFile, err = os.Open(tempFile.Name())
	data, err := io.ReadAll(tempFile)
	assertNil(err)

	linesAfter := make([]string, 0)

	for line := range strings.Lines(string(data)) {
		line = strings.TrimSpace(line)
		linesAfter = append(linesAfter, line)
	}

	actions := make([]func(), 0)

linesBeforeLoop:
	for _, lineBefore := range linesBefore {
		lineBeforeMatches := linePattern.FindStringSubmatch(lineBefore)
		lineBeforePrefix := lineBeforeMatches[1]
		lineBeforeFileName := lineBeforeMatches[2]
		for _, lineAfter := range linesAfter {
			lineAfterMatches := linePattern.FindStringSubmatch(lineAfter)
			if lineAfterMatches != nil {
				lineAfterPrefix := lineAfterMatches[1]
				lineAfterFileName := lineAfterMatches[2]
				if lineBeforePrefix == lineAfterPrefix {
					if lineBeforeFileName != lineAfterFileName {
						fmt.Println(lineBeforeFileName, "->", lineAfterFileName)
						actions = append(actions, func() {
							err = os.Rename(lineBeforeFileName, lineAfterFileName)
							if err != nil {
								fmt.Fprint(os.Stderr, err)
							}
						})
					}
					continue linesBeforeLoop
				}
			}
		}
		fmt.Println("Delete", lineBeforeFileName)
		actions = append(actions, func() {
			err = os.RemoveAll(lineBeforeFileName)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		})
	}

	for _, lineAfter := range linesAfter {
		if !linePattern.MatchString(lineAfter) && lineAfter != "" {
			if strings.HasSuffix(lineAfter, string(os.PathSeparator)) {
				fmt.Println("Mkdir", lineAfter)
				actions = append(actions, func() {
					err = os.MkdirAll(lineAfter, 0755)
					if err != nil {
						fmt.Fprint(os.Stderr, err)
					}
				})
			} else {
				fmt.Println("Touch", lineAfter)
				actions = append(actions, func() {
					err = os.MkdirAll(filepath.Dir(lineAfter), 0755)
					if err != nil {
						fmt.Fprint(os.Stderr, err)
					} else {
						newFile, err := os.Create(lineAfter)
						if err != nil {
							fmt.Fprint(os.Stderr, err)
						} else {
							newFile.Close()
						}
					}
				})
			}
		}
	}

	if len(actions) == 0 {
		os.Exit(0)
	}

confirmation:
	fmt.Print("Apply changes? (y/N) ")
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	assertNil(err)
	answer = strings.TrimSpace(answer)
	switch answer {
	case "y":
		fallthrough
	case "Y":
		for _, action := range actions {
			action()
		}
	case "n":
		fallthrough
	case "N":
		os.Exit(1)
	default:
		goto confirmation
	}
}
