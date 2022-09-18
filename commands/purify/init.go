package purify

import (
	"net/url"
)

type tracks []interface {
	match(*url.URL) []string
	handle(*Stage) *url.URL
	allowed(*url.URL) (string, bool)
}

type Stage struct {
	Deep int
	Url  *url.URL

	matches []string
}

func (t tracks) test(s string) int {
	if u, err := url.Parse(s); err == nil {
		return t.Test(u)
	}
	return -1
}

func (t tracks) Test(u *url.URL) int {
	for i, v := range t {
		if len(v.match(u)) < 1 {
			continue
		}

		allowed, stop := v.allowed(u)
		removal := validExpr(u, parseExpr(allowed))

		if !stop {
			return i
		} else if len(removal) > 0 {
			return 1 << 20
		}
	}

	// test the part of url as a fragment
	if u.Fragment != "" {
		if r := t.test(u.Fragment); r >= 0 {
			return r
		}
	}

	// test for sub queries at all
	for _, q := range u.Query() {
		for _, v := range q {
			if r := t.test(v); r >= 0 {
				return r
			}
		}
	}
	return -1
}

func (t tracks) do(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	} else if u = t.Do(&Stage{Url: u}); u != nil {
		return u.String()
	}
	return s
}

func (t tracks) Do(s *Stage) *url.URL {
	s.Deep++
	if s.Deep > 5 {
		return nil // avoid infinite loop
	}

	for _, v := range t {
		matches := v.match(s.Url)
		if len(matches) < 1 {
			continue
		}

		// remove the queries that are not allowed
		qs := s.Url.Query()
		allowed, stop := v.allowed(s.Url)
		removal := validExpr(s.Url, parseExpr(allowed))
		for _, k := range removal {
			qs.Del(k)
		}
		s.Url.RawQuery = qs.Encode()

		// handle the url with called the custom handler
		if !stop {
			s.matches = matches
			s.Url = v.handle(s)
		}

		// do purify for the queries of the rest recursively
		qs = s.Url.Query()
		for _, q := range qs {
			for i, v := range q {
				q[i] = t.do(v)
			}
		}
		s.Url.RawQuery = qs.Encode()
		break
	}

	// if there fragment is still being, purify it recursively
	if s.Url.Fragment != "" {
		s.Url.Fragment = t.do(s.Url.Fragment)
	}

	return s.Url
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
