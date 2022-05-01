package utils

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func LinkedName(user *tgbotapi.User) string {
	var lastName = user.LastName
	if lastName != "" {
		lastName = " " + lastName
	}

	return fmt.Sprintf(`<a href="tg://user?id=%d">%s%s</a>`, user.ID, user.FirstName, lastName)
}

func Serverless() bool {
	return os.Getenv("VERCEL") == "1"
}

func Config(name string) string {
	if Serverless() {
		return os.Getenv(strings.ToUpper(name))
	}

	file, err := ioutil.ReadFile(".config")
	if err != nil {
		log.Fatal(err)
	}

	config := map[string]string{}
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	value, _ := config[name]
	return value
}

func CreateBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(Config("bot_token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func NeedToIgnore(bot *tgbotapi.BotAPI, text string) bool {
	matches := regexp.MustCompile(`^\s*/\w+@(\w+(?:Bot|_bot))\b`).FindStringSubmatch(text)
	if len(matches) == 0 {
		return false
	}

	if matches[1] != bot.Self.UserName {
		return true
	}
	return false
}
