```
This tool combines the contents of one or more input files by joining their
lines in a structured, line-by-line manner. When multiple files are provided,
the tool merges corresponding lines from each file together, creating a new
output where each line consists of the combined content from the same line
position across all input files. When a single file is provided, the tool
combines the lines within that file into a single continuous joined sequence

Usage:
  th-join-lines [OPTIONS] <Files>...

Options:
  -s, --separator <Separator>  The separator how lines shall be joined. (default: " ")
  -h, --help                   Show this help message and exit

Positional arguments:
  Files                        The files which lines shall be joined. (required, can be specified multiple times)

```
