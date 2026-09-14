package vercel

import (
	"crypto/hmac"
	"fmt"
	"net/http"

	"github.com/sxyazi/bendan/commands"
	. "github.com/sxyazi/bendan/utils"
)

func HookHandler(w http.ResponseWriter, r *http.Request) {
	secret := Config("webhook_secret")
	if !hmac.Equal([]byte(r.Header.Get("X-Telegram-Bot-Api-Secret-Token")), []byte(secret)) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	bot := CreateBot()
	update, err := bot.HandleUpdate(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	commands.Bot = bot
	commands.Handle(update)
	fmt.Fprint(w, "request processed")
}
