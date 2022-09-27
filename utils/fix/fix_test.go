package fix

import "testing"

func TestFix(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{`
package main

func main() {
	rand.Seed(33)
	fmt.Println(rand.Intn(10))
}`, `import"math/rand";import"fmt";`},
	}

	for _, d := range data {
		if got := Imports(d.input); got != d.want {
			t.Errorf("Imports(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
