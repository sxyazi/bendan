package yes

import (
	"testing"
)

func TestRightTokenize(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{"对吗", "sub="},
		{"XX行不行呢", "sub=XX"},
		{"应该是吧", "sub="},
		{"和之前一样是吧", "sub=和之前一样"},
	}

	for _, d := range data {
		if got := RightTokenize(d.input); got.String() != d.want {
			t.Errorf("RightTokenize(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
