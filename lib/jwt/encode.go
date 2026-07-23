package jwt

import (
	"encoding/base64"
	"encoding/json/jsontext"
	"fmt"
	"strings"
)

func encodeJson(str string) (string, error) {
	v := jsontext.Value([]byte(str))
	err := v.Compact()
	if err != nil { return "", err} //nofmt
	return base64.RawURLEncoding.EncodeToString(v), nil
}

func Encode(input string) (string, error) {
	parts := strings.Split(input, "\n\n")
	if len(parts) != 3 { return "", fmt.Errorf("Invalid input. Make sure HEADER, PAYLOAD and SIGNATURE are delimited by '\\n\\n'.") } //nofmt
	header, err := encodeJson(parts[0])
	if err != nil { return "", err } //nofmt
	payload, err := encodeJson(parts[1])
	if err != nil { return "", err } //nofmt
	signature := strings.TrimSpace(parts[2])
	return fmt.Sprintf("%s.%s.%s", header, payload, signature), nil
}
