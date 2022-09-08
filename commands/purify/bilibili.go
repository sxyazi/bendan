package purify

import (
	"fmt"
	"math"
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

func bilibili(m []string) string {
	id := m[1]
	if id[:2] == "BV" {
		id = bvToAv(id)
	}
	return fmt.Sprintf("https://www.bilibili.com/video/%s", id)
}
