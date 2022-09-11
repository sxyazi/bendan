package purify

import (
	"github.com/sxyazi/bendan/utils"
	"net/url"
	"regexp"
	"strings"
)

type short struct{}

var nonShortHost = []string{
	`t\.me`,
}

var reShortUrl = regexp.MustCompile(`^https?://[a-zA-Z0-9]{1,5}\.[a-zA-Z0-9]{2,3}/[a-zA-Z0-9_-]{4,8}$`)
var reShortNonHost = regexp.MustCompile(`(?i)^(` + strings.Join(nonShortHost, "|") + `)$`)

func (*short) match(u *url.URL) []string {
	if reShortNonHost.MatchString(u.Hostname()) {
		return nil
	}
	return reShortUrl.FindStringSubmatch(u.String())
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
