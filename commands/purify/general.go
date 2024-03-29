package purify

import (
	"net/url"
	"regexp"
	"strings"
)

// prefix matching
var generalParams = []string{
	"utm_source",
	"utm_medium",
	"utm_term",
	"utm_content",
	"utm_campaign",
	"utm_referrer",
	"yclid",
	"gclid",
	"fbclid",
	"_openstat",
	"fb_action_ids",
	"fb_comment_id",
	"fb_action_types",
	"fb_ref",
	"fb_source",
	"action_object_map",
	"action_type_map",
	"action_ref_map",
	"spm_id_from",
	"spm",
	"from_source",
	"from_spmid",
	"vd_source",
	"_trans_",
}

var reGeneral = regexp.MustCompile(`(?i)\b(` + strings.Join(generalParams, "|") + `)`)

type general struct{}

func (*general) match(u *url.URL) []string {
	return reGeneral.FindStringSubmatch(u.RawQuery)
}

func (*general) allowed(*url.URL) (string, uint8) {
	return "", 0
}

func (*general) handle(s *Stage) *url.URL {
	qs := s.URL.Query()
	for name := range qs {
		if reGeneral.MatchString(name) {
			qs.Del(name)
		}
	}

	s.URL.RawQuery = qs.Encode()
	return s.URL
}
