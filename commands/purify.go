package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
	"strings"
	"sync"
)

func Purify(msg *tgbotapi.Message) bool {
	urls := utils.ExtractLinks(msg.Text + "\n" + msg.Caption)
	if len(urls) < 1 {
		return true
	}

	tickers := make([]*purify.Tracker, 0, len(urls))
	for _, url := range urls {
		if t := purify.Tracks.Match(url); t != nil {
			tickers = append(tickers, t)
		}
	}
	if len(tickers) < 1 {
		return true
	}

	wg := sync.WaitGroup{}
	results := make([]string, len(tickers))
	for i, t := range tickers {
		wg.Add(1)
		go func(i int, t *purify.Tracker) {
			defer wg.Done()
			results[i] = purify.Tracks.Handle(t)
		}(i, t)
	}

	sent := ReplyText(msg, "Purifying up the URLs...")
	if sent == nil {
		return true
	}
	wg.Wait()

	text := ""
	for _, r := range results {
		if r != "" {
			text += r + "\n"
		}
	}
	if text == "" {
		DeleteMessage(sent)
	} else if strings.Count(text, "\n") == 1 {
		EditText(sent, "<b>Purified URL:</b> "+text)
	} else {
		EditText(sent, "<b>The URLs purified below:</b>\n\n"+text)
	}
	return true
}
