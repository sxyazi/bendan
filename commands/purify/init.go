package purify

import (
	"net/url"
)

type tracks []interface {
	match(*url.URL) []string
	handle(*Stage) string
	allowed(*url.URL) string
}

type Stage struct {
	idx     int
	deep    int
	orig    string
	url     *url.URL
	matches []string
}

func (t tracks) Test(u string) *Stage {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil
	}

	for i, v := range t {
		matches := v.match(parsed)
		if len(matches) < 1 {
			continue
		}

		removal := validExpr(parsed, parseExpr(v.allowed(parsed)))
		qs := parsed.Query()
		for _, r := range removal {
			qs.Del(r)
		}

		parsed.RawQuery = qs.Encode()
		return &Stage{i, 0, u, parsed, matches}
	}
	return nil
}

func (t tracks) Do(s *Stage) string {
	if s == nil {
		return ""
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
