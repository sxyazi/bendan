package purify

import (
	collect "github.com/sxyazi/go-collection"
	"net/url"
	"regexp"
	"strings"
)

type tracks struct {
	re      []*regexp.Regexp
	cb      []func([]string, *url.URL) string
	allowed []map[string][]string
}

type Tracker struct {
	idx  int
	orig string
	url  *url.URL
}

func (k *tracks) add(re, al string, cb func([]string, *url.URL) string) {
	k.re = append(k.re, regexp.MustCompile(re))
	k.cb = append(k.cb, cb)
	k.allowed = append(k.allowed, parseExpr(al))
}

// Match returns the matched tracker
// al == "" && cb == nil, not allowed
// al == "" && cb != nil, (Y) force to use cb
// al != "" && cb == nil, (Y|N) only use al
// al != "" && cb != nil, (Y) use cb after al is applied
func (k *tracks) Match(u string) *Tracker {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil
	}

	for i, re := range k.re {
		if !re.MatchString(u) {
			continue
		}

		removal := validExpr(parsed, k.allowed[i])
		qs := parsed.Query()
		for _, r := range removal {
			qs.Del(r)
		}
		parsed.RawQuery = qs.Encode()

		if k.cb[i] != nil {
			return &Tracker{i, u, parsed}
		} else if len(removal) > 0 {
			return &Tracker{i, u, parsed}
		}
	}
	return nil
}

func (k *tracks) Handle(t *Tracker) string {
	if t == nil {
		return ""
	}

	cb := k.cb[t.idx]
	if cb == nil {
		return t.url.String()
	}
	return cb(k.re[t.idx].FindStringSubmatch(t.orig), t.url)
}

func (k *tracks) Purify(u string) string {
	if t := k.Match(u); t != nil {
		return k.Handle(t)
	}
	return u
}

var Tracks = &tracks{}

func init() {
	Tracks.add(`https?://youtu.be/([a-zA-Z0-9_-]{10,})`, ``, youtube)
	Tracks.add(`https?://b23.tv/([a-zA-Z0-9]{6,})`, ``, bilibili)

	Tracks.add(`https?://twitter.com/(.+?/status/\d+)`, `-`, nil)
	Tracks.add(`https?://(?:www\.)?bilibili\.com/video/(av\d+)`, `p:pi`, nil)

	Tracks.add(`https?://(?:www\.)?bilibili\.com/video/(BV[a-zA-Z0-9]{10,})`, `p:pi`, bilibili)
	Tracks.add(strings.Join(collect.Keys(generalParams), "|"), ``, general)
}
