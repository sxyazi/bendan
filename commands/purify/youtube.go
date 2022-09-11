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

func (*youtube) handle(s *Stage) string {
	return "https://www.youtube.com/watch?v=" + s.matches[1]
}

func (*youtube) allowed(*url.URL) (string, bool) {
	return "", false
}
