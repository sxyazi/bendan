package eval

import (
	"strings"
	"testing"
)

func TestNode_Eval(t *testing.T) {
	data := []struct {
		input string
		want  string
	}{
		{`1+1`, "2\n"},
		{`console.log("hello world")`, "hello world\n"},
		{`"console.log()"`, "console.log()\n"},
		{`'hello;world'`, "hello;world\n"},
		{`console.log('foo');console.log('bar')`, "foo\nbar\n"},
		{`'foo\n' + console.log('bar')`, "bar\nfoo\nundefined\n"},
	}

	for _, d := range data {
		if got := NewNode().Eval(d.input); strings.Join(got, "") != d.want {
			t.Errorf("Eval(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
