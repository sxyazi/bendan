package commands

import (
	"net/url"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
)

var allowedPreviewDomain = []string{
	"fxtwitter.com",
}

func Purify(msg *tgbotapi.Message) bool {
	urls := utils.ExtractUrls(msg.Text + "\n" + msg.Caption)
	if len(urls) < 1 {
		return false
	}

	todo := make([]*url.URL, 0, len(urls))
	for _, u := range urls {
		if purify.Tracks.Test(u) {
			todo = append(todo, u)
		}
	}
	if len(todo) < 1 {
		return false
	}

	wg := sync.WaitGroup{}
	wg.Add(len(todo))
	for i, u := range todo {
		go func(i int, u *url.URL) {
			defer wg.Done()
			todo[i] = purify.Tracks.Do(&purify.Stage{URL: u})
		}(i, u)
	}

	sent := ReplyText(msg, "Purifying up the URLs...")
	if sent == nil {
		return true
	}
	wg.Wait()

	var text strings.Builder
	for _, u := range todo {
		if u != nil {
			text.WriteString(u.String())
			text.WriteByte('\n')
		}
	}

	if s := text.String(); len(allowedPreviewDomain) > 0 {
		for _, domain := range allowedPreviewDomain {
			if strings.Contains(s, domain) {
				EditTextWithWebPagePreview(sent, "<b>Purified URL:</b> "+s)
				return true
			}
		}
	}

	if text.Len() < 1 {
		DeleteMessage(sent)
	} else if s := text.String(); strings.Count(s, "\n") == 1 {
		EditText(sent, "<b>Purified URL:</b> "+s)
	} else {
		EditText(sent, "<b>The URLs purified below:</b>\n\n"+s)
	}
	return true
}
