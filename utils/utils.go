package utils

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func Value[T any](first T, rest ...any) T {
	return first
}

func LinkedName(user *tgbotapi.User) string {
	var lastName = user.LastName
	if lastName != "" {
		lastName = " " + lastName
	}

	return fmt.Sprintf(`<a href="tg://user?id=%d">%s%s</a>`, user.ID, user.FirstName, lastName)
}

func Serverless() bool {
	return os.Getenv("VERCEL") == "1"
}

func Config(name string) string {
	if Serverless() {
		return os.Getenv(strings.ToUpper(name))
	}

	file, err := os.ReadFile(".config")
	if err != nil {
		log.Fatal(err)
	}

	config := map[string]string{}
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	value, _ := config[name]
	return value
}

func CreateBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(Config("bot_token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func NeedToIgnore(bot *tgbotapi.BotAPI, text string) bool {
	matches := regexp.MustCompile(`^\s*/\w+@(\w+(?:Bot|_bot))\b`).FindStringSubmatch(text)
	if len(matches) == 0 {
		return false
	}

	if matches[1] != bot.Self.UserName {
		return true
	}
	return false
}

var reLinks = regexp.MustCompile(`https?://(?:[^/\s]+(?:\.|\b))*(/[^\s!$'()*,:;\[\]]*)?`)

func ExtractUrls(s string) []*url.URL {
	matches := reLinks.FindAllString(s, -1)
	urls := make([]*url.URL, 0, len(matches))
	occur := make(map[string]struct{}) // TODO: go-collection UniqueBy()

	for _, match := range matches {
		u, err := url.Parse(match)
		if err != nil {
			continue
		}

		escaped := u.EscapedPath()
		for strings.Contains(escaped, "//") {
			escaped = strings.Replace(escaped, "//", "/", -1)
		}

		if path, err := url.PathUnescape(escaped); err == nil {
			u.RawPath = escaped
			u.Path = path
		}

		s = u.String()
		if _, ok := occur[s]; !ok {
			occur[s] = struct{}{}
			urls = append(urls, u)
		}
	}
	return urls
}

func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 200
	t.MaxConnsPerHost = 10
	t.MaxIdleConnsPerHost = 10
	t.IdleConnTimeout = 5 * time.Minute

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: t,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

var Client = NewClient()

var reRefreshMeta = regexp.MustCompile(`(?im)<meta\s.*?http-equiv\s*=\s*['"\s]*?refresh['"\s]*?.*?>`)
var reRefreshUrl = regexp.MustCompile(`(?i);\s*URL=(.+?)['"\s]`)

func SeekLocation(u *url.URL) *url.URL {
	// Set up the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1321.0")

	// Send the request
	resp, err := Client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Check if the response is a redirect
	if l := resp.Header.Get("Location"); l != "" {
		parsed, err := url.Parse(l)
		if err == nil && parsed.String() != u.String() {
			return parsed
		}
		return nil
	}

	// Check if the response is a meta refresh
	m := reRefreshMeta.FindSubmatch(Value(io.ReadAll(resp.Body)))
	if len(m) < 1 {
		return nil
	} else if m = reRefreshUrl.FindSubmatch(m[0]); len(m) < 2 {
		return nil
	} else if parsed, err := url.Parse(string(m[1])); err != nil {
		return nil
	} else if parsed.String() != u.String() {
		return parsed
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
