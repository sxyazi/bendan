package yes

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	TypIs = iota
	TypHave
	TypIsAB
	TypHaveAB
	TypIsYesNo
	TypHaveYesNo
	TypHaveSo

	TypRight
	TypCan
	TypLook

	TypUnknown
)

const marks = `[啊阿呀吗嘛吧呢捏罢,.?!;，。？！；]`

var reClause = regexp.MustCompile(`.+?\s*(?:[,.?!:;()，。？！：；（）]+|$)`)
var reDeterminer = regexp.MustCompile(`^(啥|甚|什么|什麽|什麼|哪个|哪样|哪)`)
var reConjunction = regexp.MustCompile(`^(虽然|但是|然而|偏偏|只是|不过|至于|那么|原来|因为|由于|因此|所以|或者|如果|假如|只要|除非|倘若|即使|要是|似乎|不如|不及|尽管|而且|况且|以免|为了|于是|然后|此外|接着|应该)`)

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

func rmRec(s string, re *regexp.Regexp) string {
	for s != "" {
		if old, r := s, re.ReplaceAllString(s, ""); r != old {
			s = r
		} else {
			break
		}
	}
	return s
}
