package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendText(chat int64, text string) *tgbotapi.Message {
	m, err := Bot.Send(tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: chat},
		Text:                  text,
		ParseMode:             tgbotapi.ModeHTML,
		Entities:              nil,
		DisableWebPagePreview: true,
	})
	if err != nil {
		log.Println("Error occurred while replying text:", err)
		return nil
	}
	return &m
}

func ReplyText(msg *tgbotapi.Message, text string) *tgbotapi.Message {
	m, err := Bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           msg.Chat.ID,
			ReplyToMessageID: msg.MessageID,
		},
		Text:                  text,
		ParseMode:             tgbotapi.ModeHTML,
		Entities:              nil,
		DisableWebPagePreview: true,
	})
	if err != nil {
		log.Println("Error occurred while replying text:", err)
		return nil
	}
	return &m
}

func EditText(sent *tgbotapi.Message, text string) bool {
	_, err := Bot.Request(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    sent.Chat.ID,
			MessageID: sent.MessageID,
		},
		Text:                  text,
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	})
	if err != nil {
		log.Println("Error occurred while editing text:", err)
	}
	return err == nil
}

func DeleteMessage(sent *tgbotapi.Message) bool {
	_, err := Bot.Request(tgbotapi.DeleteMessageConfig{ChatID: sent.Chat.ID, MessageID: sent.MessageID})
	return err == nil
}

func RespondInlineQuery(id string, result any) *tgbotapi.APIResponse {
	resp, err := Bot.Request(tgbotapi.InlineConfig{
		InlineQueryID: id,
		Results:       []any{result},
		CacheTime:     60,
		IsPersonal:    true,
	})
	if err != nil {
		log.Println("Error occurred while responding inline query:", err)
		return nil
	}
	return resp
}
