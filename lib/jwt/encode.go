package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func encodeJson(str string) (string, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(str))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf.Bytes()), nil
}

func Encode(input string) (string, error) {
	parts := strings.Split(input, "\n\n")
	if len(parts) != 3 {
		return "", fmt.Errorf("Invalid input. Make sure HEADER, PAYLOAD and SIGNATURE are delimited by '\\n\\n'.")
	}

	header, err := encodeJson(parts[0])
	if err != nil {
		return "", err
	}

	payload, err := encodeJson(parts[1])
	if err != nil {
		return "", err
	}

	signature := strings.TrimSpace(parts[2])

	return fmt.Sprintf("%s.%s.%s", header, payload, signature), nil
}
