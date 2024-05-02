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

var cacheDir = filepath.Join(os.TempDir(), "/bendan/hush")

func init() {
	err := os.MkdirAll(baseHushDir, 0755)
	if err != nil {
		log.Println("hush init failed: ", err)
	}
}

func readTimeFromFile(chatID string) (time.Time, error) {
	filePath := filepath.Join(filepath.Join(baseHushDir, chatID))
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("readTimeFromFile: %v", err)
	}

	unixTime, err := strconv.ParseInt(string(fileContent), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("readTimeFromFile: %v", err)
	}

	return time.Unix(unixTime, 0), nil
}

func writeTimeToFile(chatID string) error {
	path := filepath.Join(cacheDir, chatID)
	content := fmt.Sprintf("%d", time.Now().Add(30*time.Minute).Unix())

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Errorf("writeTimeToFile: %v", err)
	}

	return err
}

func Hush(msg *tgbotapi.Message) bool {
	chatID := strconv.FormatInt(msg.Chat.ID, 10)
	if msg.ReplyToMessage != nil && strings.Contains(msg.Text, "说话") {
		err := writeTimeToFile(chatID)
		if err == nil {
			err := os.Remove(filepath.Join(filepath.Join(baseHushDir, chatID)))
			if err != nil {
				log.Println("Hush Err:", err)
				ReplyText(msg, "想说，但说不出来...")
				return true
			}
			ReplyText(msg, "已经在说了...")
			return true
		}
	}

	expiredTime, err := readTimeFromFile(chatID)

	if err == nil && time.Now().Before(expiredTime) {
		return true
	}

	if msg.ReplyToMessage != nil && strings.Contains(msg.Text, "闭嘴") {
		err := writeTimeToFile(chatID)
		if err == nil {
			ReplyText(msg, "好吧...")
			return true
		} else {
			log.Println("Hush Err:", err)
			ReplyText(msg, "想闭，但闭不了嘴...")
			return false
		}
	}
	return false
}
