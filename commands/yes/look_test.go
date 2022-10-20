package yes

import "testing"

func TestLookTokenize(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{"看看小裙子", "sub=, obj=小裙子"},
	}

	for _, d := range data {
		if got := LookTokenize(d.input); got.String() != d.want {
			t.Errorf("LookTokenize(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
