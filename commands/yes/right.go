package yes

import (
	"fmt"
	"regexp"
)

var reRight1 = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(对不对|行不行)\s*(?:%s+|$)`, marks))
var reRight2 = regexp.MustCompile(`\s*(.*?)\s*([对是行][吗嘛吧罢])\s*[.?。？]*\s*$`)

func matchOfRight(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ms := reRight1.FindStringSubmatch(s)
		if ms != nil {
			return &Token{Typ: TypRight, Sub: ms[1], Word: ms[2]}
		}

		ms = reRight2.FindStringSubmatch(s)
		if ms != nil {
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
