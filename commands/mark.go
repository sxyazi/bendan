package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"regexp"
)

var reMark = regexp.MustCompile(`^[?？¿‽]+$`)

func Mark(msg *tgbotapi.Message) bool {
	if !reMark.MatchString(msg.Text) {
		return false
	}

	ReplyText(msg, msg.Text[:rand.Intn(len(msg.Text))+1])
	return true
}
