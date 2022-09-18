package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
	"net/url"
	"strings"
	"sync"
)

func Purify(msg *tgbotapi.Message) bool {
	urls := utils.ExtractUrls(msg.Text + "\n" + msg.Caption)
	if len(urls) < 1 {
		return true
	}

	todo := make([]*url.URL, 0, len(urls))
	for _, u := range urls {
		if t := purify.Tracks.Test(u); t >= 0 {
			todo = append(todo, u)
		}
	}
	if len(todo) < 1 {
		return true
	}

	wg := sync.WaitGroup{}
	wg.Add(len(todo))
	for i, u := range todo {
		go func(i int, u *url.URL) {
			defer wg.Done()
			todo[i] = purify.Tracks.Do(&purify.Stage{Url: u})
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
	if text.Len() < 1 {
		DeleteMessage(sent)
	} else if s := text.String(); strings.Count(s, "\n") == 1 {
		EditText(sent, "<b>Purified URL:</b> "+s)
	} else {
		EditText(sent, "<b>The URLs purified below:</b>\n\n"+s)
	}
	return true
}
