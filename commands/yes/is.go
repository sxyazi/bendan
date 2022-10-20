package yes

import (
	"fmt"
	"regexp"
	"strings"
)

const marks = `[啊阿呀吗嘛吧呢捏,.?!;，。？！；]`

var reClause = regexp.MustCompile(`.+?\s*(?:[,.?!:;()，。？！：；（）]+|$)`)
var reDeterminer = regexp.MustCompile(`^(啥|甚|什么|什麽|什麼|哪个|哪样|哪)`)
var reConjunction = regexp.MustCompile(`^(虽然|但是|然而|偏偏|只是|不过|至于|那么|原来|因为|由于|因此|所以|或者|如果|假如|只要|除非|倘若|即使|要是|似乎|不如|不及|尽管|而且|况且|以免|为了|于是|然后|此外|接着)`)

var reAOrB = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*([是有])\s*(.+?)\s*%s*还是\s*(.+?)(?:%s+|$)`, marks, marks))
var reYesOrNo = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(是不是|是否|有没有|有木有|有无)\s*(.*?)(?:%s+|$)`, marks))
var reHaveSo = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*(?:这\s*么|那\s*么|多\s*么)\s*有\s*(.*?)(?:%s+|$)`, marks))
var reYes = regexp.MustCompile(fmt.Sprintf(`\s*(.*)\s*([是有])\s*(.+?)\s*%s*[吗嘛吧?!？！]+`, marks))

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
		if s[:1] == "有" {
			typ = TypHaveYesNo
		}
	case 2:
		typ = TypIsYes
		if s == "有" {
			return TypHaveYes
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
		if len(ms) > 3 {
			return &Token{Typ: typeOfIs(0, ms[2]), Sub: ms[1], Obj: ms[3], Ind: ms[4]}
		}

		ms = reYesOrNo.FindStringSubmatch(ps[i])
		if len(ms) > 2 {
			return &Token{Typ: typeOfIs(1, ms[2]), Sub: ms[1], Obj: ms[3]}
		}

		ms = reHaveSo.FindStringSubmatch(ps[i])
		if len(ms) > 2 {
			return &Token{Typ: TypHaveSo, Sub: ms[1], Obj: ms[2]}
		}

		ms = reYes.FindStringSubmatch(ps[i])
		if len(ms) > 2 &&
			!reDeterminer.MatchString(ms[1]) &&
			!reDeterminer.MatchString(ms[2]) {
			return &Token{Typ: typeOfIs(2, ms[2]), Sub: ms[1], Obj: ms[3]}
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

	rmRec := func(s string, re *regexp.Regexp) string {
		for s != "" {
			if old, r := s, re.ReplaceAllString(s, ""); r != old {
				s = r
			} else {
				break
			}
		}
		return s
	}

	// remove conjunctions
	token.Sub = rmRec(token.Sub, reConjunction)

	// remove determiners
	token.Obj = rmRec(token.Obj, reDeterminer)
	token.Ind = rmRec(token.Ind, reDeterminer)

	// All the objects are determiners that have undetermined, so we can't do the options, just ignore them.
	if token.Obj == "" && token.Ind == "" &&
		token.Typ != TypIsYesNo &&
		token.Typ != TypHaveYesNo &&
		token.Typ != TypHaveSo /* Since it("是否", "有无", "这么有") expects a clear Yes or No */ {
		return nil
	}

	return token
}
