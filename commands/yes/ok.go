package yes

import (
	"fmt"
	"regexp"
)

var reOkOrNot = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*行不行\s*(?:%s+|$)`, marks))

func OkTokenize(s string) *Token {
	parts := explode(s)
	for i := len(parts) - 1; i >= 0; i-- {
		matches := reOkOrNot.FindStringSubmatch(s)
		if len(matches) > 1 {
			return &Token{Typ: 3, Sub: matches[1]}
		}
	}

	return nil
}
