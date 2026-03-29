```
Usage:
  change-directory [OPTIONS] [Directory]

Options:
  --clear       Clears the history
  --clean       Remove unexisting directories from history.
  --choose      Choose directory to change to
  --fish        Print fish integration code
  --powershell  Print Powershell intergration code
  -h, --help    Show this help message and exit

Positional arguments:
  Directory     

Example:
  
  Fish:
  	change-directory --fish | source
  
  Powershell:
  	th-change-directory --powershell | Out-String | Invoke-Expression

```
