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

	stages := make([]*purify.Stage, 0, len(urls))
	for _, url := range urls {
		if t := purify.Tracks.Test(url); t != nil {
			stages = append(stages, t)
		}
	}
	if len(stages) < 1 {
		return true
	}

	wg := sync.WaitGroup{}
	results := make([]string, len(stages))
	for i, t := range stages {
		wg.Add(1)
		go func(i int, t *purify.Stage) {
			defer wg.Done()
			results[i] = purify.Tracks.Do(t)
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
