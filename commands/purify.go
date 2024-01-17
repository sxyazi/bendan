package commands

import (
	"net/url"
	"strings"
	"sync"

	"github.com/google/uuid"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/purify"
	"github.com/sxyazi/bendan/utils"
)

type purifyResult struct {
	before *url.URL
	after  *url.URL
}

func purifyDo(s string) chan []*purifyResult {
	urls := utils.ExtractUrls(s)
	if len(urls) < 1 {
		return nil
	}

	todo := make([]*purifyResult, 0, len(urls))
	for _, u := range urls {
		if purify.Tracks.Test(u) {
			todo = append(todo, &purifyResult{before: u})
		}
	}
	if len(todo) < 1 {
		return nil
	}

	ch := make(chan []*purifyResult, 1)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(todo))
		for _, r := range todo {
			go func(r *purifyResult) {
				defer wg.Done()
				url := *r.before // Clone a new URL without modifying the original one
				r.after = purify.Tracks.Do(&purify.Stage{URL: &url})
			}(r)
		}

		wg.Wait()
		ch <- todo
	}()

	return ch
}

func Purify(msg *tgbotapi.Message) bool {
	ch := purifyDo(msg.Text + "\n" + msg.Caption)
	if ch == nil {
		return false
	}

	sent := ReplyText(msg, "Purifying up the URLs...")
	if sent == nil {
		return true
	}

	var text strings.Builder
	for _, r := range <-ch {
		if r.after != nil {
			text.WriteString(r.after.String())
			text.WriteByte('\n')
		}
	}
	if text.Len() < 1 {
		DeleteMessage(sent)
	} else if s := text.String(); strings.Count(s, "\n") == 1 {
		EditText(sent, "<b>Purified URL:</b> "+s)
	} else {
		EditText(sent, "<b>The URL(s) purified below:</b>\n\n"+s)
	}
	return true
}

func PurifyViaQuery(query *tgbotapi.InlineQuery) bool {
	ch := purifyDo(query.Query)
	if ch == nil {
		return false
	}

	text := query.Query
	for _, r := range <-ch {
		if r.after != nil {
			text = strings.Replace(text, r.before.String(), r.after.String(), 1)
		}
	}

	if text == query.Query {
		return false
	}

	result := tgbotapi.InlineQueryResultArticle{
		Type:  "article",
		ID:    uuid.New().String(),
		Title: utils.TruncateUTF8(text, 64),
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text: text,
		},
	}
	InlineQueryResponse(query.ID, result)
	return true
}
