package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var all = []func(*tgbotapi.BotAPI, *tgbotapi.Message) bool{
	Pin,
	Me,
	Call,
}

func Handle(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	log.Printf("[%s] says: %s", update.Message.From.UserName, update.Message.Text)
	for _, command := range all {
		if command(bot, update.Message) {
			break
		}
	}
}
