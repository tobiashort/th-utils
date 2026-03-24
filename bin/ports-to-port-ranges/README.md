```
Usage:
  ports-to-port-ranges [OPTIONS] [Ports]

Options:
  -h, --help  Show this help message and exit

Positional arguments:
  Ports       Comma separated ports. Reads from Stdin if not specified.

Example:
  $ ports-to-port-ranges 1,2,3,4,5,6,11,223,445,555,556,557
  1-6,11,223,445,555-557

```
