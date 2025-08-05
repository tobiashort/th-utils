package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
)

type JWK struct {
	KeyType  string `json:"kty"`
	Modulus  string `json:"n"`
	Exponent string `json:"e"`
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: jwk-rsa-to-der [JWK]
Reads from STDIN if JWK not provided as an argument.
		
Flags:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
To convert it into pem format:
$ cat example.json | jwk-rsa-to-der | openssl rsa -inform der -RSAPublicKey_in
`)
	os.Exit(1)
}

func encodeLength(data []byte) []byte {
	length := len(data)
	if length <= 0x7F {
		return []byte{byte(length)}
	}
	lengthAsBytes := big.NewInt(int64(length)).Bytes()
	lengthLength := len(lengthAsBytes)
	if lengthLength > 127 {
		panic(fmt.Sprint("lengthLength too big:", lengthLength))
	}
	lengthEncoded := []byte{}
	lengthEncoded = append(lengthEncoded, byte(128+lengthLength))
	lengthEncoded = append(lengthEncoded, lengthAsBytes...)
	return lengthEncoded
}

func encodeInteger(data []byte) []byte {
	encoded := []byte{}
	encoded = append(encoded, 0x02)
	encoded = append(encoded, encodeLength(data)...)
	encoded = append(encoded, data...)
	return encoded
}

func encodeSequence(data []byte) []byte {
	encoded := []byte{}
	encoded = append(encoded, 0x30)
	encoded = append(encoded, encodeLength(data)...)
	encoded = append(encoded, data...)
	return encoded
}

func main() {
	help := flag.Bool("h", false, "print help")
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	if flag.NArg() > 1 {
		printUsage()
		return
	}
	input := ""
	if flag.NArg() == 1 {
		input = flag.Arg(0)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(data))
	}
	var jwk JWK
	err := json.Unmarshal([]byte(input), &jwk)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input not JSON format:", err)
		os.Exit(1)
		return
	}
	if strings.ToLower(jwk.KeyType) != "rsa" {
		fmt.Fprintf(os.Stderr, "JWK's key type (kty) is not RSA, but '%s'\n", jwk.KeyType)
		os.Exit(1)
		return
	}
	modulus, err := base64.RawURLEncoding.DecodeString(jwk.Modulus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "modulus not base64 format:", err)
		os.Exit(1)
		return
	}
	if (modulus[0] & 0x80) == 0x80 {
		// High order bit set, needs padding.
		modulus = append([]byte{0x00}, modulus...)
	}
	exponent, err := base64.RawStdEncoding.DecodeString(jwk.Exponent)
	if err != nil {
		fmt.Fprintln(os.Stderr, "exponent not base64 format:", err)
		os.Exit(1)
		return
	}
	if (exponent[0] & 0x80) == 0x80 {
		// High order bit set, needs padding.
		exponent = append([]byte{0x00}, exponent...)
	}
	sequence := []byte{}
	sequence = append(sequence, encodeInteger(modulus)...)
	sequence = append(sequence, encodeInteger(exponent)...)
	_, err = os.Stdout.Write(encodeSequence(sequence))
	if err != nil {
		panic(err)
	}
}
