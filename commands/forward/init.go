package forward

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sxyazi/bendan/db"
	"github.com/sxyazi/bendan/types"
	. "github.com/sxyazi/bendan/utils"
	collect "github.com/sxyazi/go-collection"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"strings"
	"sync"
	"unicode/utf16"
)

var Cfg struct {
	Group struct {
		Id    string `json:"id"`
		Owner string `json:"owner"`
	} `json:"group"`
	Twitter struct {
		ConsumerKey    string `json:"consumer_key"`
		ConsumerSecret string `json:"consumer_secret"`
		UserToken      string `json:"user_token"`
		UserSecret     string `json:"user_secret"`
	}
	Mastodon struct {
		Endpoint string `json:"endpoint"`
		Token    string `json:"token"`
	} `json:"mastodon"`
	AllowedTags []string `json:"allowed_tags"`
}

func init() {
	err := json.Unmarshal([]byte(Config("forward_config")), &Cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadPhotos(bot *tgbotapi.BotAPI, fms []*types.ForwardedMessage, t *twitter, m *mastodon) ([2][]string, error) {
	// non-photo-group message
	if fms[0].PhotoId == "" {
		return [2][]string{}, nil
	}

	wg := sync.WaitGroup{}
	fails := false
	photos := [2][]string{make([]string, len(fms)), make([]string, len(fms))}
	for i, fm := range fms {
		wg.Add(1)
		i, fm := i, fm
		go func() {
			defer wg.Done()

			rc, err := DownloadFile(Value(bot.GetFileDirectURL(fm.PhotoId)))
			if err != nil {
				fails = true
				return
			}
			defer rc.Close()

			buf := new(bytes.Buffer)
			photos[0][i], _ = t.uploadPhoto(io.TeeReader(rc, buf))
			photos[1][i], _ = m.uploadPhoto(buf)
			if photos[0][i] == "" || photos[1][i] == "" {
				fails = true
			}
		}()
	}

	wg.Wait()
	if fails {
		return [2][]string{}, errors.New("failed to upload photos")
	}
	return photos, nil
}

func formatText(text string, entities []tgbotapi.MessageEntity) string {
	s := utf16.Encode([]rune(text))
	offset := 0
	for _, e := range entities {
		var rep []uint16
		switch e.Type {
		case "mention":
			username := s[e.Offset+1+offset : e.Offset+e.Length+offset]
			rep = append(utf16.Encode([]rune("https://t.me/")), username...)
		case "hashtag":
			tag := string(utf16.Decode(s[e.Offset+1+offset : e.Offset+e.Length+offset]))
			if !collect.Contains(Cfg.AllowedTags, tag) {
				rep = nil
			} else {
				continue
			}
		case "text_link":
			text := s[e.Offset+offset : e.Offset+e.Length+offset]
			rep = append(text, utf16.Encode([]rune(" ("+e.URL+")"))...)
		default:
			continue
		}

		end := s[e.Offset+e.Length+offset:]
		s = append(s[:e.Offset+offset], rep...)
		s = append(s, end...)
		offset += len(rep) - e.Length
	}
	return strings.TrimSpace(string(utf16.Decode(s)))
}

func Mark(msg *tgbotapi.Message) error {
	fm := &types.ForwardedMessage{
		Id:         msg.MessageID,
		Text:       formatText(msg.Text, msg.Entities),
		ChatId:     msg.Chat.ID,
		PhotoGroup: msg.MediaGroupID,
	}

	if fm.Text == "" {
		fm.Text = formatText(msg.Caption, msg.CaptionEntities)
	}
	if len(msg.Photo) > 0 {
		fm.PhotoId = Value(collect.Last(msg.Photo)).FileID
	}
	if fm.PhotoGroup == "" {
		fm.PhotoGroup = fmt.Sprintf("single_%d", msg.MessageID)
	}
	return db.AddForwarded(fm)
}

func Forward(bot *tgbotapi.BotAPI, fms []*types.ForwardedMessage) (*types.ForwardedMessage, error) {
	t, m := newTwitter(), newMastodon()
	photos, err := uploadPhotos(bot, fms, t, m)
	if err != nil {
		return nil, err
	}

	first := *fms[0]
	if first.TweetId == "" {
		_ = t.Create(&first, photos[0])
	}
	if first.TootId == "" {
		_ = m.Create(&first, photos[1])
	}
	if first.TweetId == "" && first.TootId == "" {
		return nil, errors.New("failed to forward")
	}

	_ = db.UpdateForwardedByGroup(fms[0].PhotoGroup, fms[0].ChatId, &bson.M{
		"tweetId":  first.TweetId,
		"tweetUrl": first.TweetUrl,
		"tootId":   first.TootId,
		"tootUrl":  first.TootUrl,
	})
	if fms[0].PhotoId != "" {
		for i, fm := range fms {
			m := bson.M{}
			if fms[0].TweetId == "" {
				m["tweetPhotoId"] = photos[0][i]
			}
			if fms[0].TootId == "" {
				m["tootPhotoId"] = photos[1][i]
			}
			_ = db.UpdateForwarded(fm.Id, fm.ChatId, m)
		}
	}
	return &first, nil
}
