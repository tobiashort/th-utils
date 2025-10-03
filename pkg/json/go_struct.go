package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func ToGoStruct(b []byte) (string, error) {
	trimmed := strings.TrimSpace(string(b))
	if len(trimmed) == 0 {
		return "", fmt.Errorf("no data")
	}
	if trimmed[0] != '{' {
		return "", fmt.Errorf("expected '{', got '%c'", trimmed[0])
	}

	var a any
	err := json.Unmarshal(b, &a)
	if err != nil {
		return "", fmt.Errorf("json unmarshall error: %w", err)
	}

	snakeToPascal := func(s string) string {
		parts := strings.Split(s, "_")
		for i := range parts {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
			}
		}
		return strings.Join(parts, "")
	}

	var render func(string, any, int, *strings.Builder) error
	render = func(name string, val any, level int, sb *strings.Builder) error {
		name = snakeToPascal(name)
		fmt.Fprint(sb, strings.Repeat("\t", level))

		if val == nil {
			fmt.Fprintf(sb, "%s any\n", name)
			return nil
		}

		switch cast := val.(type) {
		case string:
			fmt.Fprintf(sb, "%s %s\n", name, reflect.TypeOf(val))
		case bool:
			fmt.Fprintf(sb, "%s %s\n", name, reflect.TypeOf(val))
		case int:
			fmt.Fprintf(sb, "%s %s\n", name, reflect.TypeOf(val))
		case float64:
			fmt.Fprintf(sb, "%s %s\n", name, reflect.TypeOf(val))
		case []any:
			fmt.Fprintf(sb, "%s []any\n", name)
		case map[string]any:
			fmt.Fprintf(sb, "%s struct {\n", name)
			for k, v := range cast {
				render(k, v, level+1, sb)
			}
			fmt.Fprint(sb, strings.Repeat("\t", level))
			fmt.Fprint(sb, "}\n")
		default:
			return fmt.Errorf("unknown type %s", reflect.TypeOf(val))
		}

		return nil
	}

	sb := strings.Builder{}

	fmt.Fprint(&sb, "type T struct {\n")

	level := 1
	switch cast := a.(type) {
	case map[string]any:
		for k, v := range cast {
			if err := render(k, v, level, &sb); err != nil {
				return "", err
			}
		}
	default:
		return "", fmt.Errorf("unhandled type: %s", reflect.TypeOf(a))
	}

	fmt.Fprint(&sb, "}")

	return sb.String(), nil
}
