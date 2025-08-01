package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type CipherSuite struct {
	Name     string
	Security string
}

func must2[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
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
	cipherSuites := make([]CipherSuite, 0)
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			panic(err)
		}

		text = strings.TrimSpace(text)

		if text != "" {
			cipherSuites = append(cipherSuites, CipherSuite{Name: text})
		}

		if err == io.EOF {
			break
		}
	}

	for idx := range cipherSuites {
		cipherSuite := &cipherSuites[idx]
		lookupSecurity(cipherSuite)
		fmt.Printf("%s (%s)\n", cipherSuite.Name, cipherSuite.Security)
	}
}
