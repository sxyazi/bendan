package commands

import (
	"crypto/sha1"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/yes"
)

func yesSel(a [2][]string, t *yes.Token) string {
	var s []string
	if rand.Float64() > .9 {
		s = a[rand.Int()&1]
	} else {
		s = a[sha1.Sum([]byte(t.String()))[0]&1]
	}
	return s[rand.Intn(len(s))]
}

func YesRight(msg *tgbotapi.Message) bool {
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		return false
	}

	token := yes.RightTokenize(msg.Text)
	if token == nil {
		return false
	}

	var opt [2][]string
	switch []rune(token.Word)[0] {
	case '对':
		opt = [2][]string{{"对", "对的", "emm好像挺对的"}, {"No", "不对", "不大对"}}
	case '是':
		opt = [2][]string{{"是", "是的", "是啊"}, {"No", "不是", "应该不是"}}
	case '有':
		opt = [2][]string{{"有", "有的", "有啊"}, {"No", "没有", "没吧"}}
	case '行':
		opt = [2][]string{{"行", "行啊", "我觉得行"}, {"不行", "不太行", "我觉得不行"}}
	default: // e.g. "应该是", "应该有", etc.
		switch true {
		case strings.Contains(token.Word, "对"): // "X应该对吧"
			opt = [2][]string{{"y", "yyy", "对的", "挺对的", "没毛病"}, {"不对", "不对啊", "不大对", "明显错了", "肯定不对啊"}}
		case strings.Contains(token.Word, "是"): // "X应该是吧", "X应该是X吧"
			opt = [2][]string{{"好像是", "应该是", "还真是", "草，还真是"}, {"不是啊", "并不是", "显然不是", "我倒希望是"}}
		case strings.Contains(token.Word, "有"): // "X应该有吧", "X应该有X吧"
			opt = [2][]string{{"好像有", "应该有", "还真有", "草，还真有"}, {"没有啊", "并没有", "然而并没有", "我倒希望有"}}
		default: // "X应该行吧"
			opt = [2][]string{{"y", "yyy", "可以", "肯定行啊", "我觉得行"}, {"不行", "不太行", "应该不行", "肯定不行", "我觉得不行"}}
		}
	}

	ReplyText(msg, yesSel(opt, token))
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
		SendText(msg.Chat.ID, yesSel([2][]string{{token.Ind, token.Ind + "！"}, {token.Obj, token.Obj + "！"}}, token))
		return true
	}

	var opt [2][]string
	switch token.Typ {
	case yes.TypIs: // 是X吗
		opt = [2][]string{{"是", "是的"}, {"不是", "不是啊"}}
	case yes.TypHave: // 有X吗
		opt = [2][]string{{"有", "有的"}, {"没有", "没有啊"}}
	case yes.TypIsYesNo: // 是不是X、是否X
		opt = [2][]string{{"是", "是的"}, {"不是", "不是啊"}}
	case yes.TypHaveYesNo: // 有没有X、有无X
		opt = [2][]string{{"有", "有的", "有啊"}, {"没有", "没有啊", "并没有"}}
	case yes.TypHaveSo: // 这么有X、多么有X
		opt = [2][]string{{"是的"}, {"确实有" + token.Obj, "确实是有" + token.Obj}}
	default:
		return false
	}

	ReplyText(msg, yesSel(opt, token))
	return true
}

func YesCan(msg *tgbotapi.Message) bool {
	if msg.ForwardFromChat != nil && msg.ForwardFromChat.IsChannel() {
		return false
	}

	token := yes.CanTokenize(msg.Text)
	if token == nil {
		return false
	}

	var text string
	switch []rune(token.Word)[0] {
	case '能':
		text = yesSel([2][]string{{"能", "能！"}, {"不能", "不能！", "不，你不能"}}, token)
	case '会':
		text = yesSel([2][]string{{"会", "会！", "会的"}, {"不会", "不会啊", "不会的！"}}, token)
	}

	ReplyText(msg, text)
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

	text := yesSel([2][]string{{"看看", "想看"}, {"窝也想看", "想看，gkd"}}, token)
	if msg.ReplyToMessage == nil {
		SendText(msg.Chat.ID, text)
	} else {
		ReplyText(msg.ReplyToMessage, text)
	}
	return true
}
