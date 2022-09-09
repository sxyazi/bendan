package purify

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
)

type bilibili struct{}

var reBilibiliBV = regexp.MustCompile(`^https?://(?:www\.)?bilibili\.com/video/(BV[a-zA-Z0-9]{10,})`)
var reBilibiliAV = regexp.MustCompile(`^https?://(?:www\.)?bilibili\.com/video/(av\d+)`)

func (b *bilibili) match(u *url.URL) []string {
	if m := reBilibiliBV.FindStringSubmatch(u.String()); len(m) > 0 {
		return m
	}
	return reBilibiliAV.FindStringSubmatch(u.String())
}

func (b *bilibili) handle(s *Stage) string {
	id := s.matches[1]
	if id[:2] == "BV" {
		id = b.bvToAv(id)
	}

	s.url.Path = fmt.Sprintf("/video/%s", id)
	return s.url.String()
}

func (b *bilibili) allowed(*url.URL) string {
	return "p:pi"
}

func (b *bilibili) bvToAv(s string) string {
	t := "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
	tr := make(map[rune]int, len(t))
	n := []int{11, 10, 3, 8, 4, 6}
	for i, c := range t {
		tr[c] = i
	}

	var r int
	for i := 0; i < 6; i++ {
		r += tr[rune(s[n[i]])] * int(math.Pow(58, float64(i)))
	}
	return fmt.Sprintf("av%d", (r-8728348608)^177451812)
}
