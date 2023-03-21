package commands

import (
	"math/rand"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/yes"
)

var reMark = regexp.MustCompile(`^[?？¿‽]+$`)

func Mark(msg *tgbotapi.Message) bool {
	if !reMark.MatchString(msg.Text) {
		return false
	}

	text := ""
	if rand.Float64() > .9 {
		text = []string{"啊？", "嗯？"}[rand.Intn(2)]
	} else {
		text = yesSel([2][]string{
			{"?", "？", "¿"},
			{msg.Text[:rand.Intn(len(msg.Text))+1]},
		}, &yes.Token{Sub: msg.Text})
	}

	ReplyText(msg, text)
	return true
}
