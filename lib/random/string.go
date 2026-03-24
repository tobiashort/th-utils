package random

import (
	"math/rand"
)

const (
	Symbols      = "!\"#$%&'()*+,-./:;<>=?@[\\]^_`{|}~"
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers      = "0123456789"
	Hexadecimals = Numbers + "abcdef"
)

func String(alphabet string, width int, seed int64) string {
	rnd := rand.New(rand.NewSource(seed))
	b := make([]byte, width)
	for i := range width {
		b[i] = alphabet[rnd.Intn(len(alphabet))]
	}
	return string(b)
}
