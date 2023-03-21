package commands

import (
	"fmt"
	"strings"
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
)

func targetOfInteraction(msg *tgbotapi.Message) string {
	// Raise a call blandly
	r := msg.ReplyToMessage
	if r == nil {
		return "自己"
	} else if r.SenderChat != nil && msg.SenderChat != nil && r.SenderChat.ID == msg.SenderChat.ID {
		return "自己"
	} else if r.From.ID == msg.From.ID {
		return "自己"
	}

	// Aiming at one mentioned by this bot
	if r.From.ID == Bot.Self.ID {
		for _, e := range r.Entities {
			switch e.Type {
			case "text_mention":
				if e.User == nil {
					continue
				}
				if e.User.ID != msg.From.ID && e.User.ID != Bot.Self.ID {
					return SenderName(&tgbotapi.Message{From: e.User})
				}
			case "text_link":
				if !strings.HasPrefix(e.URL, "tg://resolve?domain=") {
					continue
				}
				if msg.SenderChat == nil || msg.SenderChat.UserName != e.URL[20:] {
					s := utf16.Encode([]rune(r.Text))
					return fmt.Sprintf(`<a href="%s">%s</a>`, e.URL,
						string(utf16.Decode(s[e.Offset:e.Offset+e.Length])))
				}
			}
		}
		return ""
	}

	return SenderName(r)
}

func Call(msg *tgbotapi.Message) bool {
	// valid:   /call
	// invalid: //call
	if len(msg.Text) < 2 || msg.Text[0] != '/' || msg.Text[1] == '/' {
		return false
	}

	// Bot does not in the interaction
	target := targetOfInteraction(msg)
	if target == "" {
		return false
	}

	var message string
	params := strings.Fields(msg.Text[1:])
	switch len(params) {
	case 1:
		message = fmt.Sprintf(`%s %s了 %s ！`, SenderName(msg), params[0], target)
	case 2:
		message = fmt.Sprintf(`%s %s %s %s！`, SenderName(msg), params[0], target, params[1])
	default:
		return false
	}

	SendText(msg.Chat.ID, message)
	return true
}
