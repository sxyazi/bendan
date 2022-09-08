package purify

import (
	collect "github.com/sxyazi/go-collection"
	"net/url"
	"regexp"
	"strings"
)

var generalParams = map[string]struct{}{
	"utm_source":        {},
	"utm_medium":        {},
	"utm_term":          {},
	"utm_content":       {},
	"utm_campaign":      {},
	"utm_referrer":      {},
	"yclid":             {},
	"gclid":             {},
	"fbclid":            {},
	"_openstat":         {},
	"fb_action_ids":     {},
	"fb_comment_id":     {},
	"fb_action_types":   {},
	"fb_ref":            {},
	"fb_source":         {},
	"action_object_map": {},
	"action_type_map":   {},
	"action_ref_map":    {},
	"spm_id_from":       {},
	"spm":               {},
	"from_source":       {},
	"from_spmid":        {},
	"vd_source":         {},
}

var generalRe = regexp.MustCompile(`(?i)` + strings.Join(collect.Keys(generalParams), "|"))

func general(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return ""
	}

	qs := parsed.Query()
	for p := range generalParams {
		qs.Del(p)
	}
	parsed.RawQuery = qs.Encode()
	return parsed.String()
}
