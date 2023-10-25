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

func (*twitter) allowed(*url.URL) (string, uint8) {
	return ".*", 0
}

func (*twitter) handle(s *Stage) *url.URL {
	if len(s.matches) != 2 {
		return s.URL
	}

	newURL, err := url.Parse("https://fxtwitter.com/" + s.matches[1])
	if err != nil {
		return s.URL
	}

	newURL.RawQuery = s.URL.RawQuery
	return newURL
}
