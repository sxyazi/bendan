package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/commands/forward"
	"github.com/sxyazi/bendan/db"
	"github.com/sxyazi/bendan/types"
	"log"
)

func ForwardMark(msg *tgbotapi.Message) bool {
	if !msg.IsAutomaticForward {
		return false
	} else if msg.ForwardFromChat == nil {
		return true
	} else if fmt.Sprint(msg.Chat.ID) != forward.Cfg.Group.Id {
		return true
	}

	group := msg.MediaGroupID
	if group == "" {
		group = fmt.Sprintf("single_%d", msg.MessageID)
	}
	if fm, err := db.GetForwardedGroupOne(group, msg.Chat.ID); err != nil {
		return true
	} else if fm != nil {
		_ = forward.Mark(msg)
		return true
	}

	sent := ReplyText(msg, "It's ready to sync to other platforms. Reply to this message with any option to operate:\n\n[<b>c</b>]reate")
	if sent == nil {
		return true
	}

	rm := &types.RepliedMessage{
		Id:        sent.MessageID,
		ChatId:    sent.Chat.ID,
		RepliedTo: msg.MessageID,
	}
	if err := db.AddReplied(rm); err == nil {
		_ = forward.Mark(msg)
	} else {
		DeleteMessage(sent)
	}
	return true
}

func Forward(msg *tgbotapi.Message) bool {
	if msg.ReplyToMessage == nil {
		return false
	} else if msg.ReplyToMessage.From.ID != Bot.Self.ID {
		return false
	} else if msg.Text != "c" && msg.Text != "d" {
		return false
	} else if fmt.Sprint(msg.Chat.ID) != forward.Cfg.Group.Id {
		return false
	} else if fmt.Sprint(msg.From.ID) != forward.Cfg.Group.Owner {
		return false
	}

	rm, _ := db.GetReplied(msg.ReplyToMessage.MessageID, msg.Chat.ID)
	if rm == nil {
		return false
	}

	DeleteMessage(msg)
	fms, _ := db.GetForwarded(rm.RepliedTo, rm.ChatId)
	if len(fms) < 1 {
		EditText(msg.ReplyToMessage, "The original message is not reachable, try to update this message once on Telegram.\n\n[<b>c</b>]reate")
		return true
	}

	if msg.Text == "c" && (fms[0].TweetId == "" || fms[0].TootId == "") {
		if fm, err := forward.Forward(Bot, fms); err != nil {
			log.Println("Failed to forward:", err)
			EditText(msg.ReplyToMessage, "Failed to forward this message, to continue with the option(s) below:\n\n[<b>c</b>]reate")
		} else {
			EditText(msg.ReplyToMessage, fmt.Sprintf(
				`Twitter: <a href="%s">%s</a>`+"\n"+
					`Mastodon: <a href="%s">%s</a>`,
				fm.TweetUrl, fm.TweetId,
				fm.TootUrl, fm.TootId,
			))
		}
	} else if msg.Text == "d" && (fms[0].TweetId != "" || fms[0].TootId != "") {
		//TODO
	}
	return true
}
