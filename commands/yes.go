package commands

import (
	"crypto/sha1"
	"encoding/binary"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/yes"
	"math/bits"
	"math/rand"
)

func Yes(msg *tgbotapi.Message) bool {
	token := yes.Tokenize(msg.Text)
	if token == nil {
		return false
	}

	var certain bool
	if rand.Float64() > .9 {
		certain = rand.Int()%10 >= 4
	} else {
		a := sha1.Sum([]byte(token.String()))
		certain = bits.OnesCount64(binary.BigEndian.Uint64(a[:])) > 32
	}

	sel := func(a [2]string) string {
		if certain {
			return a[0]
		}
		return a[1]
	}

	if token.Ind == "" {
		SendText(msg.Chat.ID, sel([...]string{"是", "不是"}))
	} else {
		SendText(msg.Chat.ID, sel([...]string{token.Ind, token.Obj}))
	}

	return true
}
