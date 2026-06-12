package unicode

import "testing"

func TestWidth(t *testing.T) {
	w := WidthString("abcd📦")
	expected := 6
	if w != expected {
		t.Log("expected", expected, "got", w)
		t.Fail()
	}
}
