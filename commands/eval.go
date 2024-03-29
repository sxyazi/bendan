package commands

import (
	"bytes"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/eval"
)

var reEval = regexp.MustCompile(`(?mi)^//\s*(go|golang|js|javascript|node|nodejs)[\s\n]+([\s\S]+)`)

func Eval(msg *tgbotapi.Message) bool {
	matches := reEval.FindStringSubmatch(msg.Text)
	if len(matches) < 3 {
		return false
	}

	result := make(chan []string, 1)
	go func() {
		switch strings.ToLower(matches[1]) {
		case "go", "golang":
			result <- eval.NewGo().Eval(matches[2])
		case "js", "javascript", "node", "nodejs":
			result <- eval.NewNode().Eval(matches[2])
		default:
			result <- []string{"Unknown language"}
		}
	}()

	sent := ReplyText(msg, "Evaluating...")
	if sent == nil {
		return true
	}

	var buf bytes.Buffer
	for _, s := range <-result {
		if s == "" {
			continue
		}
		buf.WriteString(`<code>`)
		buf.WriteString(strings.ReplaceAll(s, "<", "&lt;"))
		buf.WriteString(`</code>`)
	}

	if buf.Len() == 0 {
		buf.WriteString("No output")
	}

	EditText(sent, buf.String())
	return true
}
