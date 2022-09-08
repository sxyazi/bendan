package purify

import "fmt"

func twitter(m []string) string {
	return fmt.Sprintf("https://twitter.com/%s", m[1])
}
