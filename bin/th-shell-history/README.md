```
Usage:
  th-shell-history [OPTIONS]

Options:
  -f, --fish         Integration for fish shell
  -p, --powershell   Integration for powershell
  -i, --integration  Print integration code
  -h, --help         Show this help message and exit

Example:
  Fish:
    th-shell-history --integration --fish | source
    th-change-directory --integration --powershell | Out-String | Invoke-Expression

```
