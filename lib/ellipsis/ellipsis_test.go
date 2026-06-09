package ellipsis

import (
	"fmt"
	"testing"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
)

func TestAnsiStart(t *testing.T) {
	for l := range 24 {
		t.Run(fmt.Sprintf("l:%d", l), func(t *testing.T) {
			f := cfmt.Formatter{ForceColors: true}
			text := f.Sprint("#c{This is cyan.}#r{ This is red.}#y{ This is yellow.}#p{ This is purple.}")
			text = Ellipsis(text, l, "...", PosStart)
			f.Print(text)
			if len(ansi.Strip(text)) != max(l, 3) {
				t.Log(len(ansi.Strip(text)))
				t.Fail()
			}
		})
	}
}

func TestAnsiEnd(t *testing.T) {
	for l := range 24 {
		t.Run(fmt.Sprintf("l:%d", l), func(t *testing.T) {
			f := cfmt.Formatter{ForceColors: true}
			text := f.Sprint("#c{This is cyan.}#r{ This is red.}#y{ This is yellow.}#p{ This is purple.}")
			text = Ellipsis(text, l, "...", PosEnd)
			f.Print(text)
			if len(ansi.Strip(text)) != max(l, 3) {
				t.Log(len(ansi.Strip(text)))
				t.Fail()
			}
		})
	}
}

func TestAnsiCenter(t *testing.T) {
	for l := range 24 {
		t.Run(fmt.Sprintf("l:%d", l), func(t *testing.T) {
			f := cfmt.Formatter{ForceColors: true}
			text := f.Sprint("#c{This is cyan.}#r{ This is red.}#y{ This is yellow.}#p{ This is purple.}")
			text = Ellipsis(text, l, "...", PosCenter)
			f.Print(text)
			if len(ansi.Strip(text)) != max(l, 3) {
				t.Log(len(ansi.Strip(text)))
				t.Fail()
			}
		})
	}
}
