package vercel

import (
	"fmt"
	"github.com/sxyazi/bendan/commands"
	. "github.com/sxyazi/bendan/utils"
	"net/http"
)

func HookHandler(w http.ResponseWriter, r *http.Request) {
	bot := CreateBot()
	update, err := bot.HandleUpdate(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	commands.Handle(bot, update)
	fmt.Fprint(w, "request processed")
}
