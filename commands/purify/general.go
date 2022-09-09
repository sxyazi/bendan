package purify

import (
	"net/url"
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

func general(_ []string, u *url.URL) string {
	qs := u.Query()
	for p := range generalParams {
		qs.Del(p)
	}

	u.RawQuery = qs.Encode()
	return u.String()
}
