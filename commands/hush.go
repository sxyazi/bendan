package commands

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const hushDuration = 30 * time.Minute

var hushDir = filepath.Join(os.TempDir(), "/bendan/hush")
var reHush = regexp.MustCompile("别说话|闭嘴|安静")
var reUnHush = regexp.MustCompile("说话")

func init() {
	if err := os.MkdirAll(hushDir, 0755); err != nil {
		log.Println("Hush init failed:", err)
	}
}

func Hush(msg *tgbotapi.Message) bool {
	path := filepath.Join(hushDir, strconv.FormatInt(msg.Chat.ID, 10))
	repliesToBot := msg.ReplyToMessage != nil && msg.ReplyToMessage.From.ID == Bot.Self.ID

	if repliesToBot && reHush.MatchString(msg.Text) {
		if err := os.WriteFile(path, nil, 0644); err == nil {
			ReplyText(msg, "😭")
		} else {
			log.Println("Hush failed:", err)
			ReplyText(msg, "不要！")
		}
		return true
	}

	if repliesToBot && reUnHush.MatchString(msg.Text) {
		if err := os.Remove(path); err == nil || errors.Is(err, os.ErrNotExist) {
			ReplyText(msg, "好耶！")
		} else {
			log.Println("UnHush failed:", err)
			ReplyText(msg, "不要！")
		}
		return true
	}

	info, err := os.Lstat(path)
	return err == nil && time.Now().Before(info.ModTime().Add(hushDuration))
}
