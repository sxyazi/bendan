package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
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

var Bot *tgbotapi.BotAPI

func Handle(update *tgbotapi.Update) {
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
	for _, command := range all {
		if command(message) {
			break
		}
	}
}
