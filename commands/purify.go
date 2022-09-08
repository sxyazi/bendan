package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
	"sync"
)

func Purify(msg *tgbotapi.Message) bool {
	urls := utils.ExtractLinks(msg.Text + "\n" + msg.Caption)
	if len(urls) < 1 {
		return true
	}

	purifiers := make(map[string]int, len(urls))
	for _, url := range urls {
		if i := purify.KnownTracks.Match(url); i != -1 {
			purifiers[url] = i
		}
	}
	if len(purifiers) < 1 {
		return true
	}

	sent := ReplyText(msg, "Purifying up the URLs...")
	if sent == nil {
		return true
	}

	wg := sync.WaitGroup{}
	results := make(chan string, len(purifiers))
	for url, i := range purifiers {
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()
			results <- purify.KnownTracks.Handle(url, i)
		}(url, i)
	}

	wg.Wait()
	close(results)

	text := "<b>The URLs purified below:</b>\n\n"
	for r := range results {
		text += r + "\n"
	}
	EditText(sent, text)
	return true
}
