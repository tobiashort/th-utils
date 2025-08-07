package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tobiashort/clap-go"
)

func must2[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

type Args struct {
}

type CipherSuite struct {
	Name     string
	Security string
}

func lookupSecurity(cipherSuite *CipherSuite) {
	url := fmt.Sprintf("https://ciphersuite.info/api/cs/%s/", cipherSuite.Name)
	client := http.Client{Timeout: 10 * time.Second}
	resp := must2(client.Get(url))

	if resp.StatusCode == http.StatusNotFound {
		cipherSuite.Security = "unknown"
		return
	}

	data := must2(io.ReadAll(resp.Body))
	dict := make(map[string]interface{})
	json.Unmarshal(data, &dict)
	security := dict[cipherSuite.Name].(map[string]interface{})["security"].(string)
	cipherSuite.Security = security
}

func main() {
	args := Args{}
	clap.Description("Reads cipher suites line by line from Stdin and checks them")
	clap.Parse(&args)

	cipherSuites := make([]CipherSuite, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		cipherSuites = append(cipherSuites, CipherSuite{Name: text})
	}

	for idx := range cipherSuites {
		cipherSuite := &cipherSuites[idx]
		lookupSecurity(cipherSuite)
		fmt.Printf("%s (%s)\n", cipherSuite.Name, cipherSuite.Security)
	}
}
