package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/utils"
	"io"
	"log"
	"net/http"
	"regexp"
)

// This DC query method is based on checking the user profile photos.
// It may not work if the user doesn't have a username or profile photos or their photos is hidden.
func dcQuery(username string) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://t.me/%s", username), nil)
	if err != nil {
		log.Printf("DC Query: Error occurred while creating HTTP request: %v\n", err)
		return ""
	}
	resp, err := utils.Client.Do(req)
	if err != nil {
		log.Printf("DC Query: Error occurred while sending HTTP request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("DC Query: Error occurred while reading response body: %v\n", err)
		return ""
	}
	reTelegramCdn := regexp.MustCompile(`cdn(\d)\.(?:telegram-cdn\.org/)`)
	if match := reTelegramCdn.FindStringSubmatch(string(body[:])); len(match) == 2 {
		return match[1]
	}
	return ""
}

func Whoami(msg *tgbotapi.Message) bool {
	if msg.Text != "//whoami" {
		return false
	}

	dc := dcQuery(msg.From.UserName)
	text := fmt.Sprintf("User ID: <code>%d</code>\nChat ID: <code>%d</code>", msg.From.ID, msg.Chat.ID)
	if dc != "" {
		text += fmt.Sprintf("\nUser DC: <code>%s</code>", dc)
	}
	ReplyText(msg, text)
	return true
}
