package yes

import (
	"fmt"
	"strings"
)

const (
	TypIs = iota
	TypHave
	TypIsEnd
	TypHaveEnd
	TypIsAB
	TypHaveAB
	TypIsYesNo
	TypHaveYesNo
	TypHaveSo

	TypOk

	TypCan
	TypWill

	TypLook

	TypUnknown
)

type Token struct {
	Typ  uint8
	Sub  string
	Obj  string
	Ind  string
	Word string
}

func (t *Token) String() string {
	if t == nil {
		return ""
	} else if t.Obj == "" {
		return fmt.Sprintf("sub=%s", t.Sub)
	} else if t.Ind == "" {
		return fmt.Sprintf("sub=%s, obj=%s", t.Sub, t.Obj)
	}
	return fmt.Sprintf("sub=%s, obj=%s, ind=%s", t.Sub, t.Obj, t.Ind)
}

func explode(s string) []string {
	ps := reClause.FindAllString(s, -1)

	// Since go's regex engine doesn't support look behind,
	// We have to do this to avoid wrong splitting for "，是", "，有", "，还是", etc.
	for i := len(ps) - 1; i > 0; i-- {
		switch {
		case strings.HasPrefix(ps[i], "是"):
		case strings.HasPrefix(ps[i], "有"):
		case strings.HasPrefix(ps[i], "还是"):
		default:
			continue
		}

		ps[i-1] += ps[i]
		ps = ps[:i]
	}
	return ps
}
