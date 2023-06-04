package purify

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/sxyazi/bendan/utils"
)

type short struct{}

var shortHost = []string{
	`b23\.tv`,
}

var nonShortHost = []string{
	`t\.me`,
}

var reShortHost = regexp.MustCompile(`(?i)^(` + strings.Join(shortHost, "|") + `)$`)
var reShortNonHost = regexp.MustCompile(`(?i)^(` + strings.Join(nonShortHost, "|") + `)$`)

var reShortPattern = regexp.MustCompile(`^https?://([a-zA-Z0-9]{1,5}\.[a-zA-Z0-9]{2,3}/[a-zA-Z0-9_-]{1,10})$`)

func (*short) match(u *url.URL) []string {
	host := u.Hostname()
	if reShortHost.MatchString(host) {
		return []string{u.String()}
	} else if m := reShortPattern.FindStringSubmatch(u.String()); len(m) < 2 {
		return nil
	} else if len(m[1]) > 18 {
		return nil
	} else if reShortNonHost.MatchString(u.Hostname()) {
		return nil
	} else {
		return []string{m[0]}
	}
}

func (*short) allowed(*url.URL) (string, uint8) {
	return "", 0
}

func (*short) handle(s *Stage) *url.URL {
	loc := utils.SeekLocation(s.URL)
	if loc == nil {
		loc = s.URL
	}

	return Tracks.Do(&Stage{Deep: s.Deep, URL: loc})
}
