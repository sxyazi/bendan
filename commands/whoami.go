package commands

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/utils"
	"io"
	"log"
	"net/http"
	"regexp"
)

// This DC query method is based on decoding the user profile photo file id.
// The implementation of golang code is based on this blog: https://woomai.me/talk/telegram-determine-dc-by-file-id/.
func dcQueryByFileId(fileId string) int {
	rleDecode := func(binStr string) []int {
		result := make([]int, 0)
		isPreviousZero := false
		for _, bin := range binStr {
			dec := int(bin)
			if dec == 0 {
				isPreviousZero = true
				continue
			}
			if isPreviousZero {
				zeroArray := make([]int, dec)
				for i := range zeroArray {
					zeroArray[i] = 0
				}
				result = append(result, zeroArray...)
				isPreviousZero = false
			} else {
				result = append(result, dec)
			}
		}
		return result
	}
	decodeFileId, err := base64.RawURLEncoding.DecodeString(fileId)
	if err != nil {
		log.Printf("DC query by file id: Error occurred while decoding file id: %s, err: %v\n", fileId, err)
		return -1
	}
	if data := rleDecode(string(decodeFileId[:])); len(data) >= 4 {
		if dc := data[4]; dc >= 1 && dc <= 5 {
			return dc
		}
	}
	return -1
}

// This DC query method is based on checking the user profile photo.
// It may not work if the user doesn't have a username or profile photo or their photo is hidden.
func dcQueryByUsername(username string) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://t.me/%s", username), nil)
	if err != nil {
		log.Printf("DC query by username: Error occurred while creating HTTP request: %v\n", err)
		return ""
	}
	resp, err := utils.Client.Do(req)
	if err != nil {
		log.Printf("DC query by username: Error occurred while sending HTTP request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("DC query by username: Error occurred while reading response body: %v\n", err)
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

	dc := -1
	photos, err := Bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: msg.From.ID})
	if err != nil {
		log.Printf("Whoami: Error occurred while get user profile photos: %v\n", err)
	} else if len(photos.Photos) > 0 { // Check if profile photos exists
		dc = dcQueryByFileId(photos.Photos[0][0].FileID)
	}
	text := fmt.Sprintf("User ID: <code>%d</code>\nChat ID: <code>%d</code>", msg.From.ID, msg.Chat.ID)
	if dc != -1 {
		text += fmt.Sprintf("\nUser DC: <code>%d</code>", dc)
	}
	ReplyText(msg, text)
	return true
}
