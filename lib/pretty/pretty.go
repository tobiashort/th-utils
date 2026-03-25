package pretty

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/tobiashort/th-utils/lib/cfmt"
)

func Sprint(a any) string {
	sb := &strings.Builder{}
	sprint(reflect.ValueOf(a), sb, 0)
	return sb.String()
}

func sprint(v reflect.Value, sb *strings.Builder, depth int) {
	indent := strings.Repeat("  ", depth)
	switch v.Kind() {
	case reflect.Struct:
		cfmt.Fprint(sb, "#y{{\n}")
		for i := 0; i < v.NumField(); i++ {
			name := v.Type().Field(i).Name
			value := v.Field(i)
			cfmt.Fprintf(sb, "#y{  %s%s: }", indent, name)
			sprint(value, sb, depth+1)
			fmt.Fprint(sb, "\n")
		}
		cfmt.Fprintf(sb, "#y{%s\\}}", indent)
	case reflect.Map:
		cfmt.Fprint(sb, "#y{{\n}")
		keys := v.MapKeys()
		slices.SortFunc(keys, func(a, b reflect.Value) int { return strings.Compare(a.String(), b.String()) })
		for _, key := range keys {
			cfmt.Fprintf(sb, "#y{  %s%s: }", indent, key)
			sprint(v.MapIndex(key), sb, depth+1)
			fmt.Fprint(sb, "\n")
		}
		cfmt.Fprintf(sb, "#y{%s\\}}", indent)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		cfmt.Fprint(sb, "#y{[\n}")
		for i := 0; i < v.Len(); i++ {
			cfmt.Fprintf(sb, "#y{  %s%d: }", indent, i)
			sprint(v.Index(i), sb, depth+1)
			fmt.Fprint(sb, "\n")
		}
		cfmt.Fprintf(sb, "#y{%s]}", indent)
	case reflect.Interface:
		sprint(v.Elem(), sb, depth)
	default:
		cfmt.Fprintf(sb, "#b{%v}", v)
	}
}
