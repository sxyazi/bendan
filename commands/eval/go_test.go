package eval

import (
	"strings"
	"testing"
)

func TestGo_Eval(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{"1+1", "2"},
		{`fmt.Println("10")`, "10\n"},
		{`
func add(a, b int) int {
	return a + b
}
func main() {
	fmt.Print(add(1, 2))
}
		`, "3"},
	}

	for _, d := range data {
		if got := NewGo().Eval(d.input); strings.Join(got, "") != d.want {
			t.Errorf("Eval(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
