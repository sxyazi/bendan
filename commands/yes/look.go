package yes

import (
	"regexp"
)

var reLook = regexp.MustCompile(`^\s*(看\s*看)\s*(.+)\s*$`)

func LookTokenize(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ms := reLook.FindStringSubmatch(s)
		if ms != nil {
			return &Token{Typ: TypLook, Obj: ms[2], Word: "看看"}
		}
	}

	return nil
}
