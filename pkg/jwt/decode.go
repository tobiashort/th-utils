package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func formatJson(str string) (string, error) {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(str), "", "  ")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Decode(input string) (string, error) {
	parts := strings.Split(input, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid format:", input)
	}
	encodedHeader := parts[0]
	encodedPayload := parts[1]
	signature := parts[2]
	headerAsBytes, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return "", err
	}
	payloadAsBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", err
	}
	header, err := formatJson(string(headerAsBytes))
	if err != nil {
		return "", err
	}
	payload, err := formatJson(string(payloadAsBytes))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s", header, payload, signature), nil
}
