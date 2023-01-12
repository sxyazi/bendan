package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Whoami(msg *tgbotapi.Message) bool {
	if msg.Text != "//whoami" {
		return false
	}

	ReplyText(msg, fmt.Sprintf("User ID: <code>%d</code>\nChat ID: <code>%d</code>", msg.From.ID, msg.Chat.ID))
	return true
}
