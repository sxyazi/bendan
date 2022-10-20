package yes

import (
	"fmt"
	"regexp"
)

var reCan1 = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(能不能|会不会)\s*(.*?)(?:%s+|$)`, marks))
var reCan2 = regexp.MustCompile(`\s*(.*?)\s*([能会][吗嘛吧罢])\s*[.?。？]*\s*$`)

func CanTokenize(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ms := reCan1.FindStringSubmatch(s)
		if ms != nil {
			return &Token{Typ: TypCan, Sub: ms[1], Obj: ms[3], Word: ms[2]}
		}

		ms = reCan2.FindStringSubmatch(s)
		if ms != nil {
			return &Token{Typ: TypCan, Sub: ms[1], Word: ms[2]}
		}
	}

	return nil
}
