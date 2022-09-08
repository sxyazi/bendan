package purify

import (
	"testing"
)

func TestBvToAv(t *testing.T) {
	if s := bvToAv("BV1ZP4y1o76Y"); s != "av900297685" {
		t.Error(s)
	}
}
