package inline

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

var all = map[string]func(*tgbotapi.InlineQuery) bool{
	"purify": Purify,
}

var Bot *tgbotapi.BotAPI

func Handle(update *tgbotapi.Update) {
	if update.InlineQuery != nil {
		command, _, err := validateAndExtractQuery(update.InlineQuery.Query)
		if err != nil {
			return
		}
		if commandFunc, ok := all[strings.ToLower(command)]; ok {
			commandFunc(update.InlineQuery)
			return
		}
	} else {
		return
	}
}
