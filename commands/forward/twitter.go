package forward

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sxyazi/bendan/types"
	. "github.com/sxyazi/bendan/utils"
)

var client = resty.New()

type twitter struct{}

func (*twitter) buildHeader(method, path string, params map[string]string) string {
	c := Cfg.Twitter
	vals := url.Values{}
	vals.Add("oauth_consumer_key", c.ConsumerKey)
	vals.Add("oauth_nonce", RandString(48))
	vals.Add("oauth_signature_method", "HMAC-SHA1")
	vals.Add("oauth_timestamp", strconv.Itoa(int(time.Now().Unix())))
	vals.Add("oauth_version", "1.0")
	vals.Add("oauth_token", c.UserToken)
	for k, v := range params {
		vals.Add(k, v)
	}

	base := strings.ToUpper(method) + "&" +
		url.QueryEscape(path) + "&" +
		url.QueryEscape(strings.Replace(vals.Encode(), "+", "%20", -1))
	hash := hmac.New(sha1.New, []byte(c.ConsumerSecret+"&"+c.UserSecret))
	hash.Write([]byte(base))

	header := "OAuth "
	header += fmt.Sprintf(`oauth_consumer_key="%s",`, c.ConsumerKey)
	header += fmt.Sprintf(`oauth_nonce="%s",`, vals.Get("oauth_nonce"))
	header += fmt.Sprintf(`oauth_signature="%s",`, url.QueryEscape(base64.StdEncoding.EncodeToString(hash.Sum(nil))))
	header += `oauth_signature_method="HMAC-SHA1",`
	header += fmt.Sprintf(`oauth_timestamp="%s",`, vals.Get("oauth_timestamp"))
	header += `oauth_version="1.0",`
	header += fmt.Sprintf(`oauth_token="%s"`, c.UserToken)
	return header
}

func (t *twitter) uploadPhoto(photo io.Reader) (string, error) {
	var result struct {
		MediaIDString string `json:"media_id_string"`
	}

	auth := t.buildHeader("POST", "https://upload.twitter.com/1.1/media/upload.json", nil)
	_, err := client.R().
		SetHeader("Authorization", auth).
		SetFileReader("media", "image.jpg", photo).
		SetResult(&result).
		Post("https://upload.twitter.com/1.1/media/upload.json")

	if err != nil {
		return "", err
	} else if result.MediaIDString == "" {
		return "", errors.New("no photo id found in response")
	} else {
		return result.MediaIDString, nil
	}
}

func (t *twitter) Get(id string) (string, error) {
	auth := t.buildHeader("GET", "https://api.twitter.com/2/tweets/"+id, nil)

	var result struct {
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}
	_, err := client.R().
		SetHeader("Authorization", auth).
		SetResult(&result).
		Get("https://api.twitter.com/2/tweets/" + id)

	if err != nil {
		return "", err
	} else if result.Data.ID == "" {
		return "", errors.New("no tweet found")
	} else {
		return result.Data.Text, nil
	}
}

func (t *twitter) Create(fm *types.ForwardedMessage, photos []string) error {
	type media struct {
		MediaIds []string `json:"media_ids"`
	}
	var body struct {
		Text  string `json:"text"`
		Media *media `json:"media,omitempty"`
	}

	body.Text = fm.Text
	if len(photos) > 0 {
		body.Media = &media{MediaIds: photos}
	}

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	auth := t.buildHeader("POST", "https://api.twitter.com/2/tweets", nil)
	resp, err := client.R().
		SetHeader("Authorization", auth).
		SetBody(body).
		SetResult(&result).
		Post("https://api.twitter.com/2/tweets")

	if err != nil {
		log.Println("Error posting to twitter:", err)
		return err
	} else if result.Data.ID == "" {
		log.Println("Error posting to twitter, response:", resp.String())
		return errors.New("no tweet id found in response")
	} else {
		fm.TweetID = result.Data.ID
		fm.TweetURL = "https://twitter.com/_/status/" + result.Data.ID
		return nil
	}
}

func newTwitter() *twitter {
	return &twitter{}
}
