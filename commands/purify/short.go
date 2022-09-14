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

var reShortUrl = regexp.MustCompile(`^https?://([a-zA-Z0-9]{1,5}\.[a-zA-Z0-9]{2,3}/[a-zA-Z0-9_-]{1,10})$`)
var reShortNonHost = regexp.MustCompile(`(?i)^(` + strings.Join(nonShortHost, "|") + `)$`)

func (*short) match(u *url.URL) []string {
	if reShortNonHost.MatchString(u.Hostname()) {
		return nil
	}
	if m := reShortUrl.FindStringSubmatch(u.String()); len(m) < 2 {
		return nil
	} else if len(m[1]) > 18 {
		return nil
	} else {
		return m
	}
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
