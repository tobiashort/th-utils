package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/cfmt"
	"github.com/tobiashort/th-utils/lib/choose"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
	"github.com/tobiashort/th-utils/lib/orderedmap"
)

type Args struct {
	Clear      bool   `clap:"short=,desc='Clears the history'"`
	Clean      bool   `clap:"short=,desc='Remove unexisting directories from history.'"`
	Choose     bool   `clap:"short=,desc='Choose directory to change to'"`
	Fish       bool   `clap:"short=,desc='Print fish integration code'"`
	Powershell bool   `clap:"short=,desc='Print Powershell intergration code'"`
	Directory  string `clap:"positional"`
}

func run() int {
	args := Args{}
	clap.Example(`
Fish:
	change-directory --fish | source

Powershell:
	th-change-directory --powershell | Out-String | Invoke-Expression
`)
	clap.Parse(&args)

	var cacheDir = must.Do2(os.UserCacheDir())
	var historyFilePath = filepath.Join(cacheDir, "cd_history")
	historyFile := must.Do2(os.OpenFile(historyFilePath, os.O_CREATE|os.O_RDONLY, 0600))
	defer historyFile.Close()

	paths := orderedmap.NewOrderedMap[string, struct{}]()
	scanner := bufio.NewScanner(historyFile)
	for scanner.Scan() {
		line := scanner.Text()
		paths.Put(line, struct{}{})
	}
	assert.Nil(scanner.Err(), "scanner error")

	if args.Directory != "" {
		directory := must.Do2(filepath.Abs(args.Directory))
		paths.Del(directory)
		if _, err := os.Stat(directory); err == nil {
			paths.Put(directory, struct{}{})
		}
		fmt.Print(directory)
	} else if args.Fish {
		fmt.Print(`function cd
    set dir (th-change-directory $argv); or return
	builtin cd $dir
end

function j
    set dir (th-change-directory --choose); or return
    builtin cd $dir
end
`)
	} else if args.Powershell {
		fmt.Print(`function global:cd {
    param([string]$Path)
    $dir = th-change-directory $Path
    if (-not $dir) { return }
    Set-Location $dir
}
function global:j {
    $dir = th-change-directory --choose
    if (-not $dir) { return }
    Set-Location $dir
}
`)
	} else if args.Clear {
		must.Do(os.Remove(historyFilePath))
		return 0
	} else if args.Clean {
		for _, p := range paths.Keys() {
			if _, err := os.Stat(p); err != nil {
				paths.Del(p)
			}
		}
	} else if args.Choose {
		ps := paths.Keys()
		slices.Reverse(ps)
		formatter := cfmt.Formatter{ForceColors: true}
		sortFunc := func(o1, o2 choose.Option, search string) int {
			sl := strings.ToLower(search)
			if sl == "" {
				return o1.Index - o2.Index
			}
			o1l := strings.ToLower(o1.Value)
			o2l := strings.ToLower(o2.Value)
			b1 := filepath.Base(o1l)
			b2 := filepath.Base(o2l)
			if strings.Contains(b1, sl) && strings.Contains(b2, sl) {
				len1 := strings.Count(o1l, string(filepath.Separator))
				len2 := strings.Count(o2l, string(filepath.Separator))
				if len1 == len2 {
					return o1.Index - o2.Index
				} else if o1.Index < 10 || o2.Index < 10 {
					return o1.Index - o2.Index
				} else {
					return len1 - len2
				}
			} else if strings.Contains(b1, sl) {
				return -1
			} else if strings.Contains(b2, sl) {
				return 1
			} else {
				len1 := strings.Count(o1l, string(filepath.Separator))
				len2 := strings.Count(o2l, string(filepath.Separator))
				if len1 == len2 {
					if c := strings.Compare(o1l, o2l); c == 0 {
						return o1.Index - o2.Index
					} else {
						return c
					}
				} else {
					return len1 - len2
				}
			}
		}
		chooser := choose.Chooser{
			Writer:    os.Stderr,
			Formatter: formatter,
			SortFunc:  sortFunc,
		}
		option, ok := chooser.One("Change directory:", choose.ToOptions(ps))
		if ok {
			paths.Del(option.Value)
			paths.Put(option.Value, struct{}{})
			fmt.Print(option.Value)
		} else {
			return 1
		}
	} else {
		return 0
	}

	must.Do(os.WriteFile(historyFilePath, []byte(strings.Join(paths.Keys(), "\n")), 0600))
	return 0
}

func main() {
	os.Exit(run())
}
