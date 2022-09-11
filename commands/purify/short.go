package purify

import (
	"github.com/sxyazi/bendan/utils"
	"net/url"
	"regexp"
)

type short struct{}

var reShort = regexp.MustCompile(`^https?://[a-zA-Z0-9]{1,5}\.[a-zA-Z0-9]{2,3}/[a-zA-Z0-9_-]{4,8}$`)

func (*short) match(u *url.URL) []string {
	return reShort.FindStringSubmatch(u.String())
}

func (*short) handle(s *Stage) string {
	loc := utils.SeekLocation(s.url)
	if loc == nil {
		return ""
	}

	if t := Tracks.Test(loc); t != nil {
		t.deep = s.deep
		return Tracks.Do(t)
	}
	return loc.String()
}

func (*short) allowed(*url.URL) (string, bool) {
	return "", false
}
