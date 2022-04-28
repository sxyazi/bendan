package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/db"
	"github.com/sxyazi/bendan/types"
	"log"
)

func unpinOldMessages(bot *tgbotapi.BotAPI, chatId int64) {
	messages, err := db.GetPinnedMessages(chatId)
	if err != nil {
		return
	}

	if len(messages) <= 10 {
		return
	}

	for _, message := range messages[10:] {
		if _, err := bot.Request(&tgbotapi.UnpinChatMessageConfig{
			ChatID:    chatId,
			MessageID: message.Id,
		}); err == nil {
			db.RemovePinnedMessage(message.Id, message.ChatId)
		}
	}
}

func Pin(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) bool {
	if msg.Text != "/pin" || msg.ReplyToMessage == nil {
		return false
	}

	_, err := db.AddPinnedMessage(&types.PinnedMessage{
		Id:     msg.ReplyToMessage.MessageID,
		ChatId: msg.Chat.ID,
	})
	if err != nil {
		sent := tgbotapi.NewMessage(msg.Chat.ID, "It seems pinned already")
		sent.ReplyToMessageID = msg.MessageID
		bot.Send(sent)
		return true
	}

	req, err := bot.Request(&tgbotapi.PinChatMessageConfig{
		ChatID:              msg.Chat.ID,
		MessageID:           msg.ReplyToMessage.MessageID,
		DisableNotification: false,
	})

	if err != nil {
		log.Println("Error pinning message:", req.Description)

		db.RemovePinnedMessage(msg.ReplyToMessage.MessageID, msg.Chat.ID)
		sent := tgbotapi.NewMessage(msg.Chat.ID, "Check if the rights are enough in the chat")
		bot.Send(sent)
	}

	unpinOldMessages(bot, msg.Chat.ID)
	return true
}
