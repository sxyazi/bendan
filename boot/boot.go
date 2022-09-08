package boot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands"
	. "github.com/sxyazi/bendan/utils"
	"net/http"
	"os"
	"time"
)

func ServePool() {
	bot := CreateBot()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	commands.Bot = bot
	bot.Request(&tgbotapi.DeleteWebhookConfig{})
	for update := range bot.GetUpdatesChan(u) {
		commands.Handle(&update)
	}
}

func ServeHook(w http.ResponseWriter, r *http.Request) {
	bot := CreateBot()
	wh, _ := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/hook/", os.Getenv("VERCEL_URL")))
	if _, err := bot.Request(wh); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	if info.LastErrorDate != 0 {
		fmt.Fprintf(w, "Telegram callback failed: %s", info.LastErrorMessage)
		return
	}

	fmt.Fprintf(w, "ok - %d", time.Now().Unix())
}
