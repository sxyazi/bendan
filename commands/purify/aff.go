package purify

import (
	"net/url"
	"regexp"
	"strings"
)

// prefix matching
var affPaths = []string{
	"reg",
	"sign",
	"log",
	"aff",
}

// prefix matching
var affParams = []string{
	"aff",
	"ref",
	"code",
}

var reAffPaths = regexp.MustCompile(`(?i)\b(` + strings.Join(affPaths, "|") + `)`)
var reAffParams = regexp.MustCompile(`(?i)\b(` + strings.Join(affParams, "|") + `)`)

type aff struct{}

func (*aff) match(u *url.URL) []string {
	if u.Path == "" || u.Path == "/" {
		// root path
	} else if !reAffPaths.MatchString(u.Path) {
		return nil
	}
	return reAffParams.FindStringSubmatch(u.RawQuery)
}

func (*aff) allowed(*url.URL) (string, uint8) {
	return "", 0
}

func (*aff) handle(s *Stage) *url.URL {
	qs := s.URL.Query()
	for name := range qs {
		if reAffParams.MatchString(name) {
			qs.Del(name)
		}
	}
	s.URL.RawQuery = qs.Encode()

	// dig it down to `general`
	return Tracks.Do(&Stage{URL: s.URL})
}
