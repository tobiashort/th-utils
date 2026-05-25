package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/tobiashort/th-utils/lib/assert"
	"github.com/tobiashort/th-utils/lib/clap"
	"github.com/tobiashort/th-utils/lib/must"
)

type Args struct {
}

func main() {
	args := Args{}
	clap.Parse(&args)

	tmpDir := must.Do2(os.MkdirTemp("", "th-notmuch-open-*"))
	defer os.RemoveAll(tmpDir)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		filePathOld := scanner.Text()
		fileOld := must.Do2(os.Open(filePathOld))
		filePathNew := filepath.Join(tmpDir, filepath.Base(filePathOld)+".eml")
		fileNew := must.Do2(os.Create(filePathNew))
		must.Do2(io.Copy(fileNew, fileOld))
		switch runtime.GOOS {
		case "darwin":
			cmd := exec.Command("open", filePathNew)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			must.Do(cmd.Start())
		case "linux":
			cmd := exec.Command("xdg-open", filePathNew)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			must.Do(cmd.Run())
		default:
			panic("not implemented: " + runtime.GOOS)
		}
	}
	assert.Nil(scanner.Err(), "scanner error")
	time.Sleep(2 * time.Second)
}
