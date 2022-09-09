package purify

import (
	"fmt"
	"github.com/sxyazi/bendan/utils"
	"math"
	"net/url"
	"strings"
)

func bvToAv(b string) string {
	t := "fZodR9XQDSUm21yCkr6zBqiveYah8bt4xsWpHnJE7jL5VG3guMTKNPAwcF"
	tr := make(map[rune]int, len(t))
	s := []int{11, 10, 3, 8, 4, 6}
	for i, c := range t {
		tr[c] = i
	}

	var r int
	for i := 0; i < 6; i++ {
		r += tr[rune(b[s[i]])] * int(math.Pow(58, float64(i)))
	}
	return fmt.Sprintf("av%d", (r-8728348608)^177451812)
}

func b23(u string) string {
	s := utils.SeekLocation(u)
	if s == "" || strings.Contains(s, "b23.tv") {
		return ""
	}
	return Tracks.Purify(s)
}

func bilibili(m []string, u *url.URL) string {
	id := m[1]
	if strings.Contains(m[0], "b23.tv") {
		return b23(m[0])
	}

	if id[:2] == "BV" {
		id = bvToAv(id)
	}

	u.Path = fmt.Sprintf("/video/%s", id)
	return u.String()
}
