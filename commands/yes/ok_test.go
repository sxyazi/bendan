package yes

import "testing"

func TestOkTokenize(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{"XX行不行呢", "sub=XX"},
	}

	for _, d := range data {
		if got := OkTokenize(d.input); got.String() != d.want {
			t.Errorf("OkTokenize(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
