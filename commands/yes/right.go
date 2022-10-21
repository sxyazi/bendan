package yes

import (
	"fmt"
	"regexp"
	"strings"
)

var reRight1 = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(对不对|行不行)\s*(?:%s+|$)`, marks))
var reRight2 = regexp.MustCompile(`\s*(.*?)\s*((?:应该|我猜|其实|确实|大概)?[对是有行])\s*(.*?)\s*[吗嘛吧罢]+\s*[.?。？]*\s*$`)

func matchOfRight(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ms := reRight1.FindStringSubmatch(s)
		if ms != nil {
			return &Token{Typ: TypRight, Sub: ms[1], Word: ms[2]}
		}

		ms = reRight2.FindStringSubmatch(s)
		if ms == nil {
			continue
		} else if strings.HasSuffix(ms[2], "是") || strings.HasSuffix(ms[2], "有") {
			return &Token{Typ: TypRight, Sub: ms[1], Obj: ms[3], Word: ms[2]}
		} else if ms[3] == "" { // Since the object of "对" and "行" shouldn't be present.
			return &Token{Typ: TypRight, Sub: ms[1], Word: ms[2]}
		}
	}

	return nil
}

func RightTokenize(s string) *Token {
	token := matchOfRight(s)
	if token == nil {
		return nil
	}

	token.Sub = rmRec(token.Sub, reConjunction)
	if reDeterminer.MatchString(token.Sub) {
		return nil
	}
	return token
}
