package yes

import (
	"regexp"
)

var reLook = regexp.MustCompile(`^\s*看\s*看\s*(.+)\s*$`)

func LookTokenize(s string) *Token {
	parts := explode(s)
	for i := len(parts) - 1; i >= 0; i-- {
		matches := reLook.FindStringSubmatch(s)
		if len(matches) > 1 {
			return &Token{Typ: TypLook, Obj: matches[1]}
		}
	}

	return nil
}
