package purify

import (
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strings"
)

type bilibili struct{}

var reBilibiliBV = regexp.MustCompile(`^https?://(?:www\.)?bilibili\.com/video/(BV[a-zA-Z0-9]{10,})`)
var reBilibiliAV = regexp.MustCompile(`^https?://(?:www\.)?bilibili\.com/video/(av\d+)`)

func (*bilibili) match(u *url.URL) []string {
	if m := reBilibiliBV.FindStringSubmatch(u.String()); len(m) > 0 {
		return m
	}
	return reBilibiliAV.FindStringSubmatch(u.String())
}

func (*bilibili) allowed(u *url.URL) (string, uint8) {
	var stop uint8 = 1
	if reBilibiliBV.MatchString(u.String()) {
		stop = 2
	}
	return "p:pi;t:pf", stop
}

func (b *bilibili) handle(s *Stage) *url.URL {
	s.URL.Path = "/video/" + b.bvToAv(s.matches[1])
	return s.URL
}

// The implementation of golang code is based on https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/misc/bvid_desc.md
func (*bilibili) bvToAv(s string) string {
	data := "FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf"
	arr := strings.Split(s, "")
	arr[3], arr[9] = arr[9], arr[3]
	arr[4], arr[7] = arr[7], arr[4]
	t := big.NewInt(0)
	for _, c := range arr[3:] {
		t.Mul(t, big.NewInt(58))
		t.Add(t, big.NewInt(int64(strings.IndexByte(data, c[0]))))
	}
	avid := t.Xor(t.And(t, big.NewInt(2251799813685247)), big.NewInt(23442827791579)).Int64()
	return fmt.Sprintf("av%d", avid)
}
