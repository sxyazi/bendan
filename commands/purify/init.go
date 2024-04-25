package purify

import (
	"net/url"
)

type tracks []interface {
	match(*url.URL) []string

	// stop:
	//   0. no
	//   1. yes
	//   2. yes, if no disallowed params are present
	allowed(*url.URL) (string, uint8)

	handle(*Stage) *url.URL
}

type Stage struct {
	Deep int
	URL  *url.URL

	matches []string
	removal []string
}

func (t tracks) test(s string) bool {
	if u, err := url.Parse(s); err == nil {
		return t.Test(u)
	}
	return false
}

func (t tracks) Test(u *url.URL) bool {
	for _, v := range t {
		if len(v.match(u)) < 1 {
			continue
		}

		if allowed, stop := v.allowed(u); stop == 0 {
			return true
		} else if len(validExpr(u, parseExpr(allowed))) > 0 {
			return true
		}
	}

	// test the fragments as a part of url
	if u.Fragment != "" {
		if t.test(u.Fragment) {
			return true
		}
	}

	// test for all sub queries
	for _, q := range u.Query() {
		for _, v := range q {
			if t.test(v) {
				return true
			}
		}
	}
	return false
}

func (t tracks) do(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	} else if u = t.Do(&Stage{URL: u}); u != nil {
		return u.String()
	}
	return s
}

func (t tracks) Do(s *Stage) *url.URL {
	s.Deep++
	if s.Deep > 5 {
		return nil // avoid infinite loop
	} else if s.URL == nil {
		return nil // avoid nil pointer
	}

	for _, v := range t {
		s.matches = v.match(s.URL)
		if len(s.matches) < 1 {
			continue
		}

		// remove the queries that are not allowed
		qs := s.URL.Query()
		allowed, stop := v.allowed(s.URL)
		s.removal = validExpr(s.URL, parseExpr(allowed))
		for _, k := range s.removal {
			qs.Del(k)
		}
		s.URL.RawQuery = qs.Encode()

		// handle the url with the custom handler called
		if stop != 1 || (stop == 2 && len(s.removal) > 0) {
			s.URL = v.handle(s)
		}
		break
	}

	// After handling, the URL might become nil,
	// such as, if it exceeds the maximum depth (s.Deep > 5).
	if s.URL == nil {
		return nil
	}

	// do purify for the queries of the rest recursively
	qs := s.URL.Query()
	for _, q := range qs {
		for i, v := range q {
			q[i] = t.do(v)
		}
	}
	s.URL.RawQuery = qs.Encode()

	// if there fragment is still being, purify it recursively
	if s.URL.Fragment != "" {
		s.URL.Fragment = t.do(s.URL.Fragment)
	}

	return s.URL
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
