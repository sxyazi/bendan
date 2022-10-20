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
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		return false
	}

	token := yes.OkTokenize(msg.Text)
	if token == nil {
		return false
	}

	text := yes_sel([2][]string{{"行", "行的", "我觉得行"}, {"不行", "不行的", "我觉得不行"}}, token)
	ReplyText(msg, text)
	return true
}

func YesCanWill(msg *tgbotapi.Message) bool {
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		return false
	}

	token := yes.CanWillTokenize(msg.Text)
	if token == nil {
		return false
	}

	var text string
	if token.Typ == yes.TypWill {
		text = yes_sel([2][]string{{"会", "会！", "会的"}, {"不会", "不会！", "不会啊"}}, token)
	} else {
		text = yes_sel([2][]string{{"能", "能！"}, {"不能", "不能！", "不，你不能"}}, token)
	}

	ReplyText(msg, text)
	return true
}

func YesIs(msg *tgbotapi.Message) bool {
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		return false
	}

	token := yes.IsTokenize(msg.Text)
	if token == nil {
		return false
	}

	// TypIsAB/TypeHaveAB
	if token.Ind != "" {
		SendText(msg.Chat.ID, yes_sel([2][]string{{token.Ind, token.Ind + "！"}, {token.Obj, token.Obj + "！"}}, token))
		return true
	}

	var opt [2][]string
	switch token.Typ {
	case yes.TypIs: // 是X吗、应该是
		if token.Word[:1] != "是" || token.Obj == "" { // e.g. "应该是", "是吧$"
			opt = [2][]string{{"还真是"}, {"并不是"}}
		} else {
			opt = [2][]string{{"是", "是的"}, {"不是", "不是啊"}}
		}
	case yes.TypHave: // 有X吗
		if token.Word[:1] != "有" || token.Obj == "" { // e.g. "应该有", "有吧$"
			opt = [2][]string{{"还真有"}, {"并没有"}}
		} else {
			opt = [2][]string{{"有", "有的"}, {"没有", "没有啊"}}
		}
	case yes.TypIsYesNo: // 是不是X、是否X
		opt = [2][]string{{"是", "是的"}, {"不是", "不是啊"}}
	case yes.TypHaveYesNo: // 有没有X、有无X
		opt = [2][]string{{"有", "有的", "有啊"}, {"没有", "没有啊", "并没有"}}
	case yes.TypHaveSo: // 这么有X、多么有X
		opt = [2][]string{{"是的"}, {"确实有" + token.Obj, "确实是有" + token.Obj}}
	default:
		return false
	}

	ReplyText(msg, yes_sel(opt, token))
	return true
}

func YesLook(msg *tgbotapi.Message) bool {
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		return false
	}

	token := yes.LookTokenize(msg.Text)
	if token == nil {
		return false
	}

	if rand.Float64() > .9 {
		opt := []string{"看看你的", "can can need"}
		ReplyText(msg, opt[rand.Intn(len(opt))])
		return true
	}

	text := yes_sel([2][]string{{"看看", "想看"}, {"窝也想看", "想看，gkd"}}, token)
	if msg.ReplyToMessage == nil {
		SendText(msg.Chat.ID, text)
	} else {
		ReplyText(msg.ReplyToMessage, text)
	}
	return true
}
