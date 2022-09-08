package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/sxyazi/bendan/utils"
	"log"
)

var all = []func(*tgbotapi.Message) bool{
	Pin,
	Me,
	Call,
	Purify,
}

var Bot *tgbotapi.BotAPI

func Handle(update *tgbotapi.Update) {
	if update.Message == nil || NeedToIgnore(Bot, update.Message.Text) {
		return
	}

	log.Printf("[%s] says: %s", update.Message.From.UserName, update.Message.Text)
	for _, command := range all {
		if command(update.Message) {
			break
		}
	}
}
