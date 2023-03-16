package commands

import (
	"testing"
)

func Test_DC_Query(t *testing.T) {
	var data = []struct {
		text string
		want string
	}{
		{"PavelDurovs", "4"},
		{"TelegramTips", "4"},
		{"TelegramTipsAR", "1"},
		{"0", ""},
		{"", ""},
	}

	for _, d := range data {
		if got := dcQuery(d.text); d.want != got {
			t.Errorf("dcQuery(%s) = %v, want %v", d.text, got, d.want)
		}
	}
}
