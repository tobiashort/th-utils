```
Usage:
  json-query [OPTIONS] <Expr> [File]

Options:
  -h, --help  Show this help message and exit

Positional arguments:
  Expr        The expression (see https://expr-lang.org/docs/language-definition) (required)
  File        The json file to read. Reads from Stdin if not specified.

Example:
  Expr‑lang can only operate on maps, not on arrays.
  Consequently, JSON arrays such as the one below are interpreted
  as a map with a single key named "items".
  
  File array.json:
  ----------------
  [
    {
      "name": "Alex",
      "age": 21,
      "groups": ["admin", "support", "user"],
      "profile": {
        "complete": true
      }
    },
    "bla",
    {
      "name": "Judith",
      "age": 59,
      "groups": ["support", "user"],
      "profile": {
        "complete": true
      }
    },
    {
      "name": "Willow",
      "age": 12,
      "groups": ["user"],
      "profile": {
        "complete": false
      }
    }
  ]
  
  Running the query:
  -----------------
  th-json-query 'filter(items, type(#) == "map") | filter("user" in .groups) | filter(.profile.complete == true) | map(.name) | join("\n")' array.json
  
  Produces:
  ---------
  Alex
  Judith

```
