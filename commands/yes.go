package commands

import (
	"crypto/sha1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/yes"
	"math/rand"
)

func yes_sel(a [2][]string, t *yes.Token) string {
	var s []string
	if rand.Float64() > .9 {
		s = a[rand.Int()&1]
	} else {
		s = a[sha1.Sum([]byte(t.String()))[0]&1]
	}
	return s[rand.Intn(len(s))]
}

func YesOk(msg *tgbotapi.Message) bool {
	token := yes.OkTokenize(msg.Text)
	if token == nil {
		return false
	}

	text := yes_sel([2][]string{{"行", "行的", "我觉得行"}, {"不行", "不行的", "我觉得不行"}}, token)
	SendText(msg.Chat.ID, text)

	return true
}

func YesIs(msg *tgbotapi.Message) bool {
	token := yes.IsTokenize(msg.Text)
	if token == nil {
		return false
	}

	if token.Ind == "" {
		SendText(msg.Chat.ID, yes_sel([2][]string{{"是", "是的"}, {"不是"}}, token))
	} else {
		SendText(msg.Chat.ID, yes_sel([2][]string{{token.Ind}, {token.Obj}}, token))
	}

	return true
}
