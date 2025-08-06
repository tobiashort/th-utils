```
Usage:
  cut-n-stitch [OPTIONS] <Delimiter> <Format>

Options:
  -h, --help  Show this help message and exit

Positional arguments:
  Delimiter   The delimiter where a given line from Stdin shall be cut. (required)
  Format      The format how the cut line shall be stitched together (required)

Example:
  $ echo "left-middle-right" | cut-n-stitch -- "-"" "{{ index . 0 }}-{{ index . 2}}"
  left-right

```
