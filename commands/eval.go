package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/eval"
	"regexp"
	"strings"
)

var reEval = regexp.MustCompile(`(?mi)^//\s*(go|golang)[\s\n]+(.+)$`)

func Eval(msg *tgbotapi.Message) bool {
	matches := reEval.FindStringSubmatch(msg.Text)
	if len(matches) < 3 {
		return false
	}

	result := make(chan string, 1)
	go func() {
		switch strings.ToLower(matches[1]) {
		case "go", "golang":
			result <- eval.NewGo().Eval(matches[2])
		default:
			result <- "Unknown language"
		}
	}()

	sent := ReplyText(msg, "Evaluating...")
	if sent == nil {
		return true
	}

	EditText(sent, <-result)
	return true
}
