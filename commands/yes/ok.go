package yes

import (
	"fmt"
	"regexp"
)

var reOkOrNot = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(行\s*不\s*行)\s*(?:%s+|$)`, marks))

func OkTokenize(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ms := reOkOrNot.FindStringSubmatch(s)
		if len(ms) > 1 {
			return &Token{Typ: TypOk, Sub: ms[1], Word: regexp.MustCompile(`\s+`).ReplaceAllString(ms[2], "")}
		}
	}

	return nil
}
