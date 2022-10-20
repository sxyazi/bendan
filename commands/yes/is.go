package yes

import (
	"fmt"
	"regexp"
	"strings"
)

var reAOrB = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*([是有])\s*(.+?)\s*%s*还是\s*(.+?)(?:%s+|$)`, marks, marks))
var reYesOrNo = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(是不是|是否|有没有|有木有|有无)\s*(.*?)(?:%s+|$)`, marks))
var reHaveSo = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(这\s*么|那\s*么|多\s*么)\s*有\s*(.*?)(?:%s+|$)`, marks))
var reYes = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*((?:应该|我猜|其实|确实|大概)?[是有])\s*(.+?)\s*%s*[吗嘛吧罢?!？！]+`, marks))

func typeOfIs(i int, s string) uint8 {
	var typ uint8 = TypUnknown
	switch i {
	case 0:
		typ = TypIsAB
		if s == "有" {
			return TypHaveAB
		}
	case 1:
		typ = TypIsYesNo
		if strings.Contains(s, "有") {
			typ = TypHaveYesNo
		}
	case 2:
		typ = TypIs
		if strings.Contains(s, "有") {
			return TypHave
		}
	}
	return typ
}

func matchOfIs(s string) *Token {
	ps := explode(s)
	for i := len(ps) - 1; i >= 0; i-- {
		ps[i] = strings.TrimSpace(ps[i])
		for strings.Contains(ps[i], "  ") {
			ps[i] = strings.Replace(ps[i], "  ", " ", -1)
		}

		ms := reAOrB.FindStringSubmatch(ps[i])
		if ms != nil {
			return &Token{Typ: typeOfIs(0, ms[2]), Sub: ms[1], Obj: ms[3], Ind: ms[4], Word: ms[2]}
		}

		ms = reYesOrNo.FindStringSubmatch(ps[i])
		if ms != nil {
			return &Token{Typ: typeOfIs(1, ms[2]), Sub: ms[1], Obj: ms[3], Word: ms[2]}
		}

		ms = reHaveSo.FindStringSubmatch(ps[i])
		if ms != nil {
			return &Token{Typ: TypHaveSo, Sub: ms[1], Obj: ms[3], Word: regexp.MustCompile(`\s*`).ReplaceAllString(ms[2], "")}
		}

		ms = reYes.FindStringSubmatch(ps[i])
		if ms != nil && !reDeterminer.MatchString(ms[1]) && !reDeterminer.MatchString(ms[3]) {
			return &Token{Typ: typeOfIs(2, ms[2]), Sub: ms[1], Obj: ms[3], Word: ms[2]}
		}
	}
	return nil
}

func IsTokenize(s string) *Token {
	token := matchOfIs(s)
	if token == nil {
		return nil
	} else if strings.HasSuffix(token.Sub, "但") {
		return nil // ignore "但是"
	}

	// remove conjunctions
	token.Sub = rmRec(token.Sub, reConjunction)

	// remove determiners
	token.Obj = rmRec(token.Obj, reDeterminer)
	token.Ind = rmRec(token.Ind, reDeterminer)

	// All the objects are determiners that have undetermined, so we can't do the options, just ignore them.
	if token.Obj == "" && token.Ind == "" &&
		(token.Typ == TypIs || token.Typ == TypHave ||
			token.Typ == TypIsAB || token.Typ == TypHaveAB) {
		return nil /* Since it not expects a clear Yes or No, thus need one less non-determiner object. */
	}

	return token
}
