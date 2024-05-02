package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var hushDir = filepath.Join(os.TempDir(), "/bendan/hush")

func init() {
	err := os.MkdirAll(hushDir, 0755)
	if err != nil {
		log.Println("Hush init failed: ", err)
	}
}

func readTimeFromFile(chatID string) (time.Time, error) {
	path := filepath.Join(hushDir, chatID)
	content, err := os.ReadFile(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("readTimeFromFile: %v", err)
	}

	unixTime, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("readTimeFromFile: %v", err)
	}

	return time.Unix(unixTime, 0), nil
}

func writeTimeToFile(chatID string) error {
	path := filepath.Join(hushDir, chatID)
	content := strconv.FormatInt(time.Now().Add(30*time.Minute).Unix(), 10)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Println("writeTimeToFile: ", err)
	}

	return err
}

func Hush(msg *tgbotapi.Message) bool {
	chatID := strconv.FormatInt(msg.Chat.ID, 10)

	if msg.ReplyToMessage != nil && strings.Contains(msg.Text, "说话") {
		err := os.Remove(filepath.Join(hushDir, chatID))
		if err != nil {
			log.Println("Hush Err:", err)
			ReplyText(msg, "想说，但说不出来...")
		} else {
			ReplyText(msg, "已经在说了...")
		}
		return true
	}

	expiredTime, err := readTimeFromFile(chatID)
	if err == nil && time.Now().Before(expiredTime) {
		return true
	}

	if msg.ReplyToMessage != nil && strings.Contains(msg.Text, "闭嘴") {
		err := writeTimeToFile(chatID)
		if err == nil {
			ReplyText(msg, "好吧...")
		} else {
			log.Println("Hush Err:", err)
			ReplyText(msg, "想闭，但闭不了嘴...")
		}
		return true
	}

	return false
}
