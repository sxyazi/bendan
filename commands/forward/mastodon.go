package forward

import (
	"errors"
	"fmt"
	"github.com/sxyazi/bendan/types"
	"io"
	"log"
	"net/url"
)

type mastodon struct{}

func (m *mastodon) uploadPhoto(photo io.Reader) (string, error) {
	var result struct{ Id string }
	_, err := client.R().
		SetAuthToken(Cfg.Mastodon.Token).
		SetFileReader("file", "image.jpg", photo).
		SetResult(&result).
		Post(Cfg.Mastodon.Endpoint + "/api/v2/media")

	if err != nil {
		return "", err
	} else if result.Id == "" {
		return "", errors.New("no photo id found in response")
	} else {
		return result.Id, nil
	}
}

func (m *mastodon) Create(fm *types.ForwardedMessage, photos []string) error {
	c := Cfg.Mastodon

	photoValues := url.Values{}
	for _, id := range photos {
		photoValues.Add("media_ids[]", id)
	}

	var result struct{ Id, Url string }
	resp, err := client.R().
		SetAuthToken(c.Token).
		SetHeader("Idempotency-Key", fmt.Sprintf("channel:%d:%s", fm.ChatId, fm.PhotoGroup)).
		SetFormDataFromValues(photoValues).
		SetFormData(map[string]string{
			"status":    fm.Text,
			"sensitive": "false",
		}).
		SetResult(&result).
		Post(c.Endpoint + "/api/v1/statuses")

	if err != nil {
		log.Println("Error posting to mastodon:", err)
		return err
	} else if result.Id == "" {
		log.Println("Error posting to mastodon, response:", resp.String())
		return errors.New("no post id found in response")
	} else {
		fm.TootId = result.Id
		fm.TootUrl = result.Url
		return nil
	}
}

func newMastodon() *mastodon {
	return &mastodon{}
}
