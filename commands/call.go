package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
	"strings"
)

func targetOfInteraction(msg *tgbotapi.Message) *tgbotapi.User {
	// Just raised a call blandly
	if msg.ReplyToMessage == nil {
		return &tgbotapi.User{ID: msg.From.ID, FirstName: "自己"}
	}

	// Aiming at one mentioned by this bot
	target := msg.ReplyToMessage.From
	if target.ID == Bot.Self.ID && len(msg.ReplyToMessage.Entities) > 0 {
		for _, entity := range msg.ReplyToMessage.Entities {
			if entity.Type != "text_mention" || entity.User == nil {
				continue
			}

			target = entity.User
			if target.ID != msg.From.ID {
				break
			}
		}
	}

	// Aiming at its previous session
	if target.ID == msg.From.ID {
		return &tgbotapi.User{ID: msg.From.ID, FirstName: "自己"}
	}

	return target
}

func Call(msg *tgbotapi.Message) bool {
	// valid:   /call
	// invalid: //call
	if len(msg.Text) < 2 || msg.Text[0] != '/' || msg.Text[1] == '/' {
		return false
	}

	// Bot does not in the interaction
	target := targetOfInteraction(msg)
	if target.ID == Bot.Self.ID {
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

	SendText(msg.Chat.ID, message)
	return true
}
