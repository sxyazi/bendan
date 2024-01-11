package inline

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
	"net/url"
	"strings"
	"sync"
)

func Purify(inlineQuery *tgbotapi.InlineQuery) bool {
	_, text, err := validateAndExtractQuery(inlineQuery.Query)
	if err != nil {
		return false
	}

	urls := utils.ExtractUrls(text + "\n")
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

	wg.Wait()

	var s strings.Builder
	for _, u := range todo {
		if u != nil {
			s.WriteString(u.String())
			s.WriteByte('\n')
		}
	}

	result := tgbotapi.InlineQueryResultArticle{
		Type:  "article",
		ID:    uuid.New().String(),
		Title: "Purified URL: " + s.String(),
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text: s.String(),
		},
	}

	SendInline(inlineQuery.ID, result)
	return true
}
