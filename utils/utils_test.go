package utils

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	collect "github.com/sxyazi/go-collection"
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

func Test_ExtractLinks(t *testing.T) {
	var data = []struct {
		text string
		want []string
	}{
		{"https://www.google.com", []string{"https://www.google.com"}},
		{"**[link](https://t.me/)**", []string{"https://t.me/"}},
		{"**[link](https://www.bilibili.com/video/av900297685?arg1=val1)**", []string{"https://www.bilibili.com/video/av900297685?arg1=val1"}},
		{"http://localhost间中 https://127.0.0.1/@bendan_bot简中", []string{"http://localhost", "https://127.0.0.1/@bendan_bot%E7%AE%80%E4%B8%AD"}},
	}

	for _, d := range data {
		var got []string
		for _, u := range ExtractUrls(d.text) {
			got = append(got, u.String())
		}

		if !collect.Same(got, d.want) {
			t.Errorf("ExtractUrls(%s) = %v, want %v", d.text, got, d.want)
		}
	}
}
