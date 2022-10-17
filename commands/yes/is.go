package yes

import (
	"fmt"
	"regexp"
	"strings"
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

const marks = `[啊阿呀吗嘛吧呢捏,.?!;，。？！；]`

var reAOrB = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*是\s*(.+?)\s*%s*还是\s*(.+?)(?:%s+|$)`, marks, marks))
var reYesOrNo = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*是不是\s*(.*?)(?:%s+|$)`, marks))
var reYes = regexp.MustCompile(`\s*(.*)\s*是\s*(.+?)\s*[吗嘛吧?!？！]+`)
var reClause = regexp.MustCompile(`.+?\s*(?:[,.?!:;()，。？！：；（）]+|$)`)
var reDeterminer = regexp.MustCompile(`^(啥|甚|什么|什麽|什麼|哪个|哪样|哪)`)
var reConjunction = regexp.MustCompile(`^(虽然|但是|然而|偏偏|只是|不过|至于|那么|原来|因为|由于|因此|所以|或者|如果|假如|只要|除非|倘若|即使|要是|似乎|不如|不及|尽管|而且|况且|以免|为了|于是|然后|此外|接着)`)

func explode(s string) []string {
	parts := reClause.FindAllString(s, -1)

	// Since go's regex engine doesn't support look behind,
	// We have to do this to avoid wrong splitting for "，是", "，还是", etc.
	for i := len(parts) - 1; i > 0; i-- {
		switch {
		case strings.HasPrefix(parts[i], "是"):
		case strings.HasPrefix(parts[i], "还是"):
		default:
			continue
		}

		parts[i-1] += parts[i]
		parts = parts[:i]
	}
	return parts
}

func match(s string) *Token {
	parts := explode(s)
	for i := len(parts) - 1; i >= 0; i-- {
		parts[i] = strings.TrimSpace(parts[i])
		for strings.Contains(parts[i], "  ") {
			parts[i] = strings.Replace(parts[i], "  ", " ", -1)
		}

		matches := reAOrB.FindStringSubmatch(parts[i])
		if len(matches) > 3 {
			return &Token{Typ: 2, Sub: matches[1], Obj: matches[2], Ind: matches[3]}
		}

		matches = reYesOrNo.FindStringSubmatch(parts[i])
		if len(matches) > 2 {
			return &Token{Typ: 1, Sub: matches[1], Obj: matches[2]}
		}

		matches = reYes.FindStringSubmatch(parts[i])
		if len(matches) > 2 &&
			!reDeterminer.MatchString(matches[1]) &&
			!reDeterminer.MatchString(matches[2]) {
			return &Token{Typ: 0, Sub: matches[1], Obj: matches[2]}
		}
	}
	return nil
}

func IsTokenize(s string) *Token {
	token := match(s)
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
		token.Typ != 1 /* Since it("是不是") expects a clear Yes or No */ {
		return nil
	}

	return token
}
