package yes

import (
	"fmt"
	"strings"
)

const (
	TypIsYes = iota
	TypHaveYes
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
	Typ uint8
	Sub string
	Obj string
	Ind string
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
	parts := reClause.FindAllString(s, -1)

	// Since go's regex engine doesn't support look behind,
	// We have to do this to avoid wrong splitting for "，是", "，有", "，还是", etc.
	for i := len(parts) - 1; i > 0; i-- {
		switch {
		case strings.HasPrefix(parts[i], "是"):
		case strings.HasPrefix(parts[i], "有"):
		case strings.HasPrefix(parts[i], "还是"):
		default:
			continue
		}

		parts[i-1] += parts[i]
		parts = parts[:i]
	}
	return parts
}
