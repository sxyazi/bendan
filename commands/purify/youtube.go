package purify

import "net/url"

func youtube(m []string, u *url.URL) string {
	return "https://www.youtube.com/watch?v=" + m[1]
}
