package forward

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/sxyazi/bendan/types"
)

type mastodon struct{}

func (*mastodon) uploadPhoto(photo io.Reader) (string, error) {
	var result struct{ ID string }
	_, err := client.R().
		SetAuthToken(Cfg.Mastodon.Token).
		SetFileReader("file", "image.jpg", photo).
		SetResult(&result).
		Post(Cfg.Mastodon.Endpoint + "/api/v2/media")

	if err != nil {
		return "", err
	} else if result.ID == "" {
		return "", errors.New("no photo id found in response")
	} else {
		return result.ID, nil
	}
}

func (*mastodon) Create(fm *types.ForwardedMessage, photos []string) error {
	c := Cfg.Mastodon

	photoValues := url.Values{}
	for _, id := range photos {
		photoValues.Add("media_ids[]", id)
	}

	var result struct{ ID, URL string }
	resp, err := client.R().
		SetAuthToken(c.Token).
		SetHeader("Idempotency-Key", fmt.Sprintf("channel:%d:%s", fm.ChatID, fm.PhotoGroup)).
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
	} else if result.ID == "" {
		log.Println("Error posting to mastodon, response:", resp.String())
		return errors.New("no post id found in response")
	} else {
		fm.TootID = result.ID
		fm.TootURL = result.URL
		return nil
	}
}

func newMastodon() *mastodon {
	return &mastodon{}
}
