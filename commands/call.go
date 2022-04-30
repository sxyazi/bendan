package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
	"strings"
)

func Call(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) bool {
	// valid:   /call
	// invalid: //call
	if len(msg.Text) < 2 || msg.Text[0] != '/' || msg.Text[1] == '/' {
		return false
	}

	// Target of the interaction
	var target *tgbotapi.User
	if msg.ReplyToMessage != nil {
		target = msg.ReplyToMessage.From
	}
	if target == nil || target.ID == msg.From.ID {
		target = &tgbotapi.User{
			ID:        msg.From.ID,
			FirstName: "自己",
		}
	}

	// Bot does not in the interaction
	if target.ID == bot.Self.ID {
		return false
	}

	var message string
	params := strings.Fields(msg.Text[1:])
	switch len(params) {
	case 1:
		message = fmt.Sprintf(`%s %s了 %s ！`, LinkedName(msg.From), params[0], LinkedName(target))
	case 2:
		message = fmt.Sprintf(`%s %s %s %s！`, LinkedName(msg.From), params[0], LinkedName(target), params[1])
	default:
		return false
	}

	sent := tgbotapi.NewMessage(msg.Chat.ID, message)
	sent.ParseMode = tgbotapi.ModeHTML
	bot.Send(sent)
	return true
}
