package jwt

import (
	"encoding/base64"
	"encoding/json/jsontext"
	"fmt"
	"strings"
)

func formatJson(str string) (string, error) {
	v := jsontext.Value([]byte(str))
	err := v.Indent(jsontext.WithIndent("  "))
	if err != nil { return "", err } //nofmt
	return v.String(), nil
}

func Decode(input string) (string, error) {
	parts := strings.Split(input, ".")
	if len(parts) != 3 { return "", fmt.Errorf("invalid format: %s", input) } //nofmt

	encodedHeader := parts[0]
	encodedPayload := parts[1]
	signature := parts[2]

	headerAsBytes, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil { return "", err } //nofmt
	payloadAsBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil { return "", err } //nofmt
	header, err := formatJson(string(headerAsBytes))
	if err != nil { return "", err } //nofmt
	payload, err := formatJson(string(payloadAsBytes))
	if err != nil { return "", err } //nofmt

	return fmt.Sprintf("%s\n\n%s\n\n%s", header, payload, signature), nil
}
