package purify

import (
	"net/url"
)

type tracks []interface {
	match(*url.URL) []string
	handle(*Stage) string
	allowed(*url.URL) (string, bool)
}

type Stage struct {
	idx     int
	deep    int
	orig    string
	url     *url.URL
	matches []string
}

func (t tracks) Test(u *url.URL) *Stage {
	orig := u.String()
	for i, v := range t {
		matches := v.match(u)
		if len(matches) < 1 {
			continue
		}

		allowed, stop := v.allowed(u)
		removal := validExpr(u, parseExpr(allowed))

		qs := u.Query()
		for _, r := range removal {
			qs.Del(r)
		}
		u.RawQuery = qs.Encode()

		if !stop {
			return &Stage{i, 0, orig, u, matches}
		} else if len(removal) > 0 {
			return &Stage{-1, 0, orig, u, matches}
		}
	}
	return nil
}

func (t tracks) Do(s *Stage) string {
	if s == nil {
		return ""
	} else if s.idx == -1 {
		return s.url.String()
	}

	s.deep++
	if s.deep > 5 {
		return "" // avoid infinite loop
	}

	return t[s.idx].handle(s)
}

var Tracks = tracks{}

func init() {
	Tracks = append(Tracks, &youtube{})
	Tracks = append(Tracks, &bilibili{})
	Tracks = append(Tracks, &twitter{})
	Tracks = append(Tracks, &aff{})
	Tracks = append(Tracks, &short{})
	Tracks = append(Tracks, &general{})
}
