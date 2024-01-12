package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
	"strings"
)

var all = []func(*tgbotapi.Message) bool{
	ForwardMark,
	Pin,
	Whoami,
	Me,
	Eval,
	Dontworry,
	Call,
	Mark,
	Forward,
	Purify,
	YesRight,
	YesIs,
	YesCan,
	YesLook,
}

var allInlineQuery = map[string]func(*tgbotapi.InlineQuery) bool{
	"purify": PurifyByInlineQuery,
	"pu":     PurifyByInlineQuery, // purify alias: purify url (pu)
}

var Bot *tgbotapi.BotAPI

func Handle(update *tgbotapi.Update) {
	var message *tgbotapi.Message
	if update.Message != nil {
		message = update.Message
	} else if update.ChannelPost != nil {
		message = update.ChannelPost
	} else if update.InlineQuery != nil {
		command, _, err := ExtractInlineQuery(update.InlineQuery.Query)
		if err != nil {
			return
		}
		if commandFunc, ok := allInlineQuery[strings.ToLower(command)]; ok {
			commandFunc(update.InlineQuery)
			return
		}
	} else {
		return
	}

	if NeedToIgnore(Bot, message.Text) {
		return
	}

	//log.Printf("[%s] says: %s", message.From.UserName, message.Text)
	for _, command := range all {
		if command(message) {
			break
		}
	}
}
