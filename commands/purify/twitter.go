package purify

import (
	"net/url"
	"regexp"
)

type twitter struct{}

var reTwitter = regexp.MustCompile(`^https?://twitter.com/(.+?/status/\d+)`)

func (*twitter) match(u *url.URL) []string {
	return reTwitter.FindStringSubmatch(u.String())
}

func (*twitter) handle(s *Stage) string {
	return s.url.String()
}

func (*twitter) allowed(*url.URL) string {
	return "-"
}
