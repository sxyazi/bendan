package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
)

var viaQuery = []func(*tgbotapi.InlineQuery) bool{
	PurifyViaQuery,
}

var viaMessage = []func(*tgbotapi.Message) bool{
	ForwardMark,
	Pin,
	Whoami,
	Me,
	Eval,
	Dontworry,
	Call,
	Forward,
	Purify,
	Hush,
	Mark,
	YesRight,
	YesIs,
	YesCan,
	YesLook,
}

var Bot *tgbotapi.BotAPI

func Handle(update *tgbotapi.Update) {
	if update.InlineQuery != nil {
		HandleQuery(update.InlineQuery)
		return
	}

	var message *tgbotapi.Message
	if update.Message != nil {
		message = update.Message
	} else if update.ChannelPost != nil {
		message = update.ChannelPost
	} else {
		return
	}

	if NeedToIgnore(Bot, message.Text) {
		return
	}

	//log.Printf("[%s] says: %s", message.From.UserName, message.Text)
	for _, f := range viaMessage {
		if f(message) {
			break
		}
	}
}

func HandleQuery(query *tgbotapi.InlineQuery) {
	for _, f := range viaQuery {
		if f(query) {
			break
		}
	}
}
