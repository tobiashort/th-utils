```
Usage:
  jwk-rsa-to-der [OPTIONS] [JWK]

Options:
  -h, --help  Show this help message and exit

Positional arguments:
  JWK         The JSON Web Key. Reads from Stdin if not specified.

Example:
  To convert it into pem format:
  $ cat example.json | jwk-rsa-to-der | openssl rsa -inform der -RSAPublicKey_in

```
