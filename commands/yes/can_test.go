package yes

import "testing"

func TestCanTokenize(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{"能吗", "sub="},
		{"A能不能B？", "sub=A, obj=B"},
		{"能不能X呢？", "sub=, obj=X"},

		{"会不会降价啊", "sub=, obj=降价"},
		{"你会不会", "sub=你"},
		{"它会不会消失", "sub=它, obj=消失"},
	}

	for _, d := range data {
		if got := CanTokenize(d.input); got.String() != d.want {
			t.Errorf("CanTokenize(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
