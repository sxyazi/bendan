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

func (*youtube) handle(s *Stage) *url.URL {
	if u, err := url.Parse("https://www.youtube.com/watch?v=" + s.matches[1]); err == nil {
		return u
	}
	return s.URL
}

func (*youtube) allowed(*url.URL) (string, bool) {
	return "", false
}
