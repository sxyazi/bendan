package yes

import (
	"fmt"
	"regexp"
)

var reWillOrNot = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(会\s*不\s*会)\s*(.*?)(?:%s+|$)`, marks))
var reCanOrNot = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(能\s*不\s*能)\s*(.*?)(?:%s+|$)`, marks))

func CanWillTokenize(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ms := reWillOrNot.FindStringSubmatch(s)
		if len(ms) > 2 {
			return &Token{Typ: TypWill, Sub: ms[1], Obj: ms[3], Word: regexp.MustCompile(`\s+`).ReplaceAllString(ms[2], "")}
		}

		ms = reCanOrNot.FindStringSubmatch(s)
		if len(ms) > 2 {
			return &Token{Typ: TypCan, Sub: ms[1], Obj: ms[3], Word: regexp.MustCompile(`\s+`).ReplaceAllString(ms[2], "")}
		}
	}

	return nil
}
