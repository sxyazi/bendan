package commands

import (
	"github.com/google/uuid"
	"net/url"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
)

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
	if text.Len() < 1 {
		DeleteMessage(sent)
	} else if s := text.String(); strings.Count(s, "\n") == 1 {
		EditText(sent, "<b>Purified URL:</b> "+s)
	} else {
		EditText(sent, "<b>The URLs purified below:</b>\n\n"+s)
	}
	return true
}

func PurifyByInlineQuery(inlineQuery *tgbotapi.InlineQuery) bool {
	urls := utils.ExtractUrls(inlineQuery.Query + "\n")
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
			inlineQuery.Query = strings.Replace(inlineQuery.Query, u.String(), purify.Tracks.Do(&purify.Stage{URL: u}).String(), -1)
		}(i, u)
	}
	wg.Wait()

	result := tgbotapi.InlineQueryResultArticle{
		Type:  "article",
		ID:    uuid.New().String(),
		Title: "Message after purified the URL(s): \n" + inlineQuery.Query,
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text:      "<b>Message after purified the URL(s): </b> \n" + inlineQuery.Query,
			ParseMode: "HTML",
		},
	}
	InlineQueryResponse(inlineQuery.ID, result)
	return true
}
