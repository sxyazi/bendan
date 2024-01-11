package inline

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func validateAndExtractQuery(query string) (string, string, error) {
	queryParts := strings.SplitN(query, " ", 2)
	if len(queryParts) != 2 {
		return "", "", errors.New("invalid query format")
	}
	return queryParts[0], queryParts[1], nil
}

func SendInline(inlineQueryID string, result any) *tgbotapi.Message {
	m, err := Bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		Results:       []any{result},
		CacheTime:     60,
		IsPersonal:    true,
	})
	if err != nil {
		log.Println("Error occurred while sending text:", err)
		return nil
	}
	return &m
}
