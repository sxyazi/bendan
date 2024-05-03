package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

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

	if msg.ReplyToMessage != nil && reHush.MatchString(msg.Text) {
		if err := os.WriteFile(path, nil, 0644); err != nil {
			log.Println("Hush Err:", err)
			ReplyText(msg, "想闭，但闭不了嘴。。。")
		} else {
			ReplyText(msg, "好吧。。。")
		}
		return true
	}

	if msg.ReplyToMessage != nil && reUnHush.MatchString(msg.Text) {
		if _, err := os.Lstat(path); err == nil {
			if err := os.Remove(path); err != nil {
				log.Println("Hush Err:", err)
				ReplyText(msg, "想说，但说不出来。。。")
			} else {
				ReplyText(msg, "哈？我又可以说话了吗？")
			}
		} else {
			ReplyText(msg, "哈？你想让我说什么？")
		}
		return true
	}

	if info, err := os.Lstat(path); err == nil && time.Now().Before(info.ModTime().Add(30*time.Minute)) {
		return true
	}

	return false
}
