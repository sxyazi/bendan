package commands

import (
	"encoding/base64"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "github.com/sxyazi/bendan/utils"
	collect "github.com/sxyazi/go-collection"
)

// This DC query method is based on decoding the user profile photo file id.
// The implementation of golang code is based on this blog: https://woomai.me/talk/telegram-determine-dc-by-file-id/
func dataCenterBy(fileID string) int {
	i, zero := 0, false
	for _, c := range string(Value(base64.RawURLEncoding.DecodeString(fileID))) {
		if int(c) == 0 {
			zero = true
			continue
		}
		if zero {
			i += int(c)
			if i > 4 {
				break
			}
			zero = false
		} else {
			if i == 4 {
				return int(c)
			}
			i++
		}
	}
	return -1
}

func Whoami(msg *tgbotapi.Message) bool {
	if msg.Text != "//whoami" {
		return false
	}

	// Basic
	text := fmt.Sprintf("User ID: <code>%d</code>\nChat ID: <code>%d</code>", msg.From.ID, msg.Chat.ID)

	// Data center
	photo, _ := collect.First(Value(Bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: msg.From.ID})).Photos)
	if p, ok := collect.First(photo); ok {
		text += fmt.Sprintf("\nUser DC: <code>%d</code>", dataCenterBy(p.FileID))
	}

	ReplyText(msg, text)
	return true
}
