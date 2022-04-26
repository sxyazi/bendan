package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands"
	. "github.com/sxyazi/bendan/utils"
	"log"
	"net/http"
	"os"
)

func servePool(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range bot.GetUpdatesChan(u) {
		commands.Handle(bot, &update)
	}
}

func serveHook(bot *tgbotapi.BotAPI) {
	wh, _ := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/hook", os.Getenv("VERCEL_URL")))
	if _, err := bot.Request(wh); err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	go http.ListenAndServe("0.0.0.0:80", nil)
	for update := range bot.ListenForWebhook("/hook/" + bot.Token) {
		commands.Handle(bot, &update)
	}
}

func main() {
	// Create a bot
	bot, err := tgbotapi.NewBotAPI(Config("bot_token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	if Serverless() {
		serveHook(bot)
	} else {
		servePool(bot)
	}
}
