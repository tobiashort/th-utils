package ellipsis

import (
	"fmt"
	"testing"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
	"github.com/tobiashort/th-utils/lib/unicode"
)

func TestAnsiStart(t *testing.T) {
	for w := range 24 {
		t.Run(fmt.Sprintf("w:%d", w), func(t *testing.T) {
			f := cfmt.Formatter{ForceColors: true}
			text := f.Sprint("#c{This i📦s cyan.📦}#r{ This i📦s red.📦}#y{ This i📦s yellow📦.}#p{ This is 📦purple.📦}")
			text = Ellipsis(text, w, "...", PosStart)
			f.Println(text)
			actual := unicode.WidthString(ansi.Strip(text))
			expected := max(w, 3)
			if actual != expected {
				t.Log("expected", expected, "got", actual)
				t.Fail()
			}
		})
	}
}

func TestAnsiEnd(t *testing.T) {
	for w := range 24 {
		t.Run(fmt.Sprintf("w:%d", w), func(t *testing.T) {
			f := cfmt.Formatter{ForceColors: true}
			text := f.Sprint("#c{This i📦s cyan.📦}#r{ This i📦s red.📦}#y{ This i📦s yellow📦.}#p{ This is 📦purple.📦}")
			text = Ellipsis(text, w, "...", PosEnd)
			f.Println(text)
			actual := unicode.WidthString(ansi.Strip(text))
			expected := max(w, 3)
			if actual != expected {
				t.Log("expected", expected, "got", actual)
				t.Fail()
			}
		})
	}
}

func TestAnsiCenter(t *testing.T) {
	for w := range 24 {
		t.Run(fmt.Sprintf("w:%d", w), func(t *testing.T) {
			f := cfmt.Formatter{ForceColors: true}
			text := f.Sprint("#c{This i📦s cyan.📦}#r{ This i📦s red.📦}#y{ This i📦s yellow📦.}#p{ This is 📦purple.📦}")
			text = Ellipsis(text, w, "...", PosCenter)
			f.Print(text)
			actual := unicode.WidthString(ansi.Strip(text))
			expected := max(w, 3)
			if actual != expected {
				t.Log("expected", expected, "got", actual)
				t.Fail()
			}
		})
	}
}
