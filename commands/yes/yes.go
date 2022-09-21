package yes

import (
	"fmt"
	"regexp"
	"strings"
)

type Token struct {
	Sub string
	Obj string
	Ind string
}

func (t *Token) String() string {
	if t == nil {
		return ""
	} else if t.Ind == "" {
		return fmt.Sprintf("sub=%s, obj=%s", t.Sub, t.Obj)
	}
	return fmt.Sprintf("sub=%s, obj=%s, ind=%s", t.Sub, t.Obj, t.Ind)
}

const marks = `[啊阿呀吗嘛吧呢捏,.?!;，。？！；]`

var reYesOr = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*是\s*(.+?)\s*%s*还是\s*(.+?)(?:%s+|$)`, marks, marks))
var reYesOrNo = regexp.MustCompile(fmt.Sprintf(`\s*(.*?)\s*是不是\s*(.*?)(?:%s+|$)`, marks))
var reYes = regexp.MustCompile(`\s*(.*)\s*是\s*(.+?)\s*[吗嘛吧?!？！]+`)
var reClause = regexp.MustCompile(`.+?\s*(?:[,.?!;，。？！；]+|$)`)
var reIgnore = regexp.MustCompile(`^(啥|甚|什么|什麽|什麼|哪个|哪样)`)

func Tokenize(s string) *Token {
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

	for i := len(parts) - 1; i >= 0; i-- {
		parts[i] = strings.TrimSpace(parts[i])
		for strings.Contains(parts[i], "  ") {
			parts[i] = strings.Replace(parts[i], "  ", " ", -1)
		}

		matches := reYesOr.FindStringSubmatch(parts[i])
		if len(matches) > 3 {
			return &Token{Sub: matches[1], Obj: matches[2], Ind: matches[3]}
		}

		matches = reYesOrNo.FindStringSubmatch(parts[i])
		if len(matches) > 2 {
			return &Token{Sub: matches[1], Obj: matches[2]}
		}

		matches = reYes.FindStringSubmatch(parts[i])
		if len(matches) > 2 &&
			!reIgnore.MatchString(matches[1]) &&
			!reIgnore.MatchString(matches[2]) {
			return &Token{Sub: matches[1], Obj: matches[2]}
		}
	}
	return nil
}
