package purify

import (
	"net/url"
	"regexp"
)

type knownTracks struct {
	re []*regexp.Regexp
	cb []func([]string) string
}

func (k *knownTracks) add(re string, cb func([]string) string) {
	k.re = append(k.re, regexp.MustCompile(re))
	k.cb = append(k.cb, cb)
}

func (k *knownTracks) Match(u string) int {
	for i, re := range k.re {
		if re.MatchString(u) {
			return i
		}
	}

	parsed, err := url.Parse(u)
	if err == nil && generalRe.MatchString(parsed.RawQuery) {
		return 1 << 10
	}
	return -1
}

func (k *knownTracks) Handle(u string, i int) string {
	if i == -1 {
		return ""
	}
	if i < len(k.re) {
		return k.cb[i](k.re[i].FindStringSubmatch(u))
	}
	// Do the general purify
	return general(u)
}

var KnownTracks = &knownTracks{}

func init() {
	KnownTracks.add("https?://youtu.be/([a-zA-Z0-9_]{10,})", youtube)
	KnownTracks.add(`https?://(?:www\.)?bilibili\.com/video/(av\d+).+`, bilibili)
	KnownTracks.add(`https?://(?:www\.)?bilibili\.com/video/(BV[a-zA-Z0-9]{10,})`, bilibili)
}
