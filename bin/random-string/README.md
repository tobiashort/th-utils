```
Usage:
  random-string [OPTIONS]

Options:
  -w, --width <Width>        Width in characters (default: 12)
  -l, --lowercase            Use lowercase characters
  -u, --uppercase            Use uppercase characters
  -n, --numbers              Use numbers
  -s, --symbols              Use symbols
  -x, --hex                  Use hexadecimals as the alphabet
  -a, --alphabet <Alphabet>  The custom alphabet to be used
  -c, --amount <Amount>      The amount of strings to be generated (default: 1)
  -h, --help                 Show this help message and exit

Example:
  $ rand-string -lun -w 12 -c 100
  Generates 100 twelve character wide random strings with
  lowercase and uppercase characters and numbers.
  
  $ rand-string -a "abc123" -w 10 -c 50
  Generates 50 ten character wide random strings that match
  the given alphabet.

```
