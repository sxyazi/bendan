package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"testing"
)

func Test_NeedToIgnore(t *testing.T) {
	data := []struct {
		string
		bool
	}{
		{"RustRss_bot", false},
		{"@RustRss_bot", false},
		{"/sub@username", false},

		{"/sub@RustRss_bot", true},
		{"/sub@RustRssBot", true},
		{" /sub@RustRssBot ", true},

		{"/sub@RustRssBot xx", true},
		{" /sub@RustRssBot xx ", true},
	}

	bot := &tgbotapi.BotAPI{Self: tgbotapi.User{UserName: "@bendan_bot"}}
	for _, d := range data {
		if NeedToIgnore(bot, d.string) != d.bool {
			t.Errorf("NeedToIgnore(%s) = %v, want %v", d.string, NeedToIgnore(bot, d.string), d.bool)
		}
	}
}
