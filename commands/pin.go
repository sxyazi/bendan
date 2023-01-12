package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/db"
	"github.com/sxyazi/bendan/types"
	"log"
)

func unpinOldMessages(chatId int64) {
	messages, err := db.GetPinned(chatId)
	if err != nil {
		return
	}

	if len(messages) <= 10 {
		return
	}

	for _, message := range messages[10:] {
		if _, err := Bot.Request(&tgbotapi.UnpinChatMessageConfig{
			ChatID:    chatId,
			MessageID: message.Id,
		}); err == nil {
			db.RemovePinned(message.Id, message.ChatId)
		}
	}
}

func Pin(msg *tgbotapi.Message) bool {
	if msg.Text != "//pin" || msg.ReplyToMessage == nil {
		return false
	}

	err := db.AddPinned(&types.PinnedMessage{
		Id:     msg.ReplyToMessage.MessageID,
		ChatId: msg.Chat.ID,
	})
	if err != nil {
		ReplyText(msg, "It seems pinned already")
		return true
	}

	req, err := Bot.Request(&tgbotapi.PinChatMessageConfig{
		ChatID:              msg.Chat.ID,
		MessageID:           msg.ReplyToMessage.MessageID,
		DisableNotification: false,
	})

	if err != nil {
		log.Println("Error pinning message:", req.Description)

		db.RemovePinned(msg.ReplyToMessage.MessageID, msg.Chat.ID)
		ReplyText(msg, "Check if the rights are enough in the chat")
	}

	unpinOldMessages(msg.Chat.ID)
	return true
}
