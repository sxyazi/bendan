package purify

import (
	"net/url"
	"strconv"
	"strings"
)

func expPi(s string) bool {
	i, err := strconv.Atoi(s)
	return err == nil && i > 0
}
func expPf(s string) bool {
	i, err := strconv.ParseFloat(s, 64)
	return err == nil && i > 0
}
func expNi(s string) bool {
	i, err := strconv.Atoi(s)
	return err == nil && i < 0
}
func expNf(s string) bool {
	i, err := strconv.ParseFloat(s, 64)
	return err == nil && i < 0
}

func parseExpr(e string) map[string][]string {
	if e == "" {
		return nil // all allowed
	} else if e == "-" {
		return make(map[string][]string, 0) // none allowed
	}

	params := strings.Split(e, ";")
	allowed := make(map[string][]string, len(params))
	for _, p := range params {
		if p == "" {
			continue
		}
		if s := strings.SplitN(p, ":", 2); len(s) == 2 {
			allowed[s[0]] = strings.Split(s[1], ",")
		}
	}
	return allowed
}

func validExpr(u *url.URL, allowed map[string][]string) (removal []string) {
	if allowed == nil {
		return // all allowed
	}

	for k, v := range u.Query() {
		if len(v) != 1 {
			removal = append(removal, k)
			continue
		}

		rules, ok := allowed[k]
		if !ok || len(rules) < 1 {
			removal = append(removal, k)
			continue // not allowed
		}

		for _, rule := range rules {
			if rule == "pi" && expPi(v[0]) {
				continue
			} else if rule == "pf" && expPf(v[0]) {
				continue
			} else if rule == "ni" && expNi(v[0]) {
				continue
			} else if rule == "nf" && expNf(v[0]) {
				continue
			} else if rule == v[0] {
				continue
			}

			removal = append(removal, k)
		}
	}
	return
}
