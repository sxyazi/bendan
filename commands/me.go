package commands

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
)

func Me(msg *tgbotapi.Message) bool {
	if !strings.HasPrefix(msg.Text, "/me") {
		return false
	}

	message := strings.TrimSpace(msg.Text[3:])
	if message == "" {
		return false
	}

	SendText(msg.Chat.ID, fmt.Sprintf("%s %s！", SenderName(msg), message))
	return true
}
