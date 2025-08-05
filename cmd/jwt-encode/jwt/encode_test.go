package jwt_test

import (
	"testing"

	jwt_decode "github.com/tobiashort/jwt-decode/jwt"
	jwt_encode "github.com/tobiashort/jwt-encode/jwt"
)

var TestDecodeEncodeExamples = []string{
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwicm9sZSI6IkFkbWluIiwiaWF0IjoxNTE2MjM5MDIyfQ.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwicGF5bG9hZCI6InRlc3RpbmcifQ.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjM0NTY3ODkwIiwicm9sZXMiOlsiQWRtaW4iLCJVc2VyIl0sImlhdCI6MTUxNjIzOTAyMn0.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
}

func TestDecodeEncode(t *testing.T) {
	for _, example := range TestDecodeEncodeExamples {
		t.Run(example, func(t *testing.T) {
			decoded, err := jwt_decode.Decode(example)
			if err != nil {
				t.Error(err)
			}
			encoded, err := jwt_encode.Encode(decoded)
			if err != nil {
				t.Error(err)
			}
			if example != encoded {
				t.Errorf("Not equal:\n%s\n%s", example, encoded)
			}
		})
	}
}
