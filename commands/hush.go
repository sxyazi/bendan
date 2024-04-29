package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var hushFilePath = filepath.Join(os.TempDir(), "hush")

func readTimeFromFile() (time.Time, error) {
	fileContent, err := os.ReadFile(hushFilePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("readTimeFromFile: %v", err)
	}

	unixTime, err := strconv.ParseInt(string(fileContent), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("readTimeFromFile: %v", err)
	}

	timestamp := time.Unix(unixTime, 0)
	return timestamp, nil
}

func writeTimeToFile() error {
	after30Minutes := time.Now().Add(30 * time.Minute)
	fileContent := strconv.Itoa(int(after30Minutes.Unix()))
	err := os.WriteFile(hushFilePath, []byte(fileContent), 0644)

	if err != nil {
		return fmt.Errorf("writeTimeToFile: %v", err)
	}

	return nil
}

func Hush(msg *tgbotapi.Message) bool {
	if msg.ReplyToMessage != nil && strings.Contains(msg.Text, "说话") {
		err := writeTimeToFile()
		if err == nil {
			err := os.Remove(hushFilePath)
			if err != nil {
				ReplyText(msg, "想说，但说不出来...")
				return true
			}
			ReplyText(msg, "已经在说了...")
			return true
		}
	}

	expiredTime, err := readTimeFromFile()

	if err == nil && time.Now().Before(expiredTime) {
		return true
	}

	if msg.ReplyToMessage != nil && strings.Contains(msg.Text, "闭嘴") {
		err := writeTimeToFile()
		if err == nil {
			ReplyText(msg, "好吧...")
			return true
		}
		return false
	}
	return false
}
