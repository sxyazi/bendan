package yes

import (
	"fmt"
	"regexp"
)

var reWillOrNot = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*会不会\s*(.*?)(?:%s+|$)`, marks))
var reCanOrNot = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*能不能\s*(.*?)(?:%s+|$)`, marks))

func CanWillTokenize(s string) *Token {
	parts := explode(s)
	for i := len(parts) - 1; i >= 0; i-- {
		matches := reWillOrNot.FindStringSubmatch(s)
		if len(matches) > 2 {
			return &Token{Typ: TypWill, Sub: matches[1], Obj: matches[2]}
		}

		matches = reCanOrNot.FindStringSubmatch(s)
		if len(matches) > 2 {
			return &Token{Typ: TypCan, Sub: matches[1], Obj: matches[2]}
		}
	}

	return nil
}
