package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gomarkdown/markdown"
	"github.com/tobiashort/clap-go"
	"github.com/tobiashort/utils-go/must"
)

type Args struct {
	File string `clap:"mandatory,positional,description='The markdown file to preview'"`
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}

func main() {
	args := Args{}
	clap.Parse(&args)

	listener := must.Do2(net.Listen("tcp", "127.0.0.1:0"))
	fmt.Println("Listening on:", listener.Addr().String())
	if err := openBrowser("http://" + listener.Addr().String()); err != nil {
		fmt.Printf("Unable to open browser: %v", err)
	}

	must.Do(
		http.Serve(
			listener,
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					md := must.Do2(os.ReadFile(args.File))
					html := markdown.ToHTML(md, nil, nil)
					must.Do2(res.Write(html))
				})))
}
