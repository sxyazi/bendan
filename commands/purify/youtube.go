package purify

import (
	"net/url"
	"regexp"
)

type youtube struct{}

var reYoutube = regexp.MustCompile(`^https?://youtu.be/([a-zA-Z0-9_-]{10,})`)

func (*youtube) match(u *url.URL) []string {
	return reYoutube.FindStringSubmatch(u.String())
}

func (*youtube) allowed(*url.URL) (string, uint8) {
	return "t:pi", 2
}

func (*youtube) handle(s *Stage) *url.URL {
	u, err := url.Parse("https://www.youtube.com/watch?v=" + s.matches[1])
	if err != nil {
		return s.URL
	}
	if t := s.URL.Query().Get("t"); t != "" {
		u.Query().Set("t", t)
	}
	return u
}
