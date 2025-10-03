package strings

import (
	"strings"

	"github.com/tobiashort/utils-go/must"
)

func Dedent(s string) string {
	builder := strings.Builder{}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if i == 0 {
			must.Do2(builder.WriteString(line))
		} else if index := strings.IndexByte(line, '|'); index > -1 {
			must.Do2(builder.WriteString(line[index+1:]))
		} else {
			must.Do2(builder.WriteString(line))
		}
		if i < len(lines)-1 {
			must.Do(builder.WriteByte('\n'))
		}
	}
	return builder.String()
}
