package commands

import (
	"testing"
)

func Test_DC_Query_By_Username(t *testing.T) {
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
		if got := dcQueryByUsername(d.text); d.want != got {
			t.Errorf("dcQuery(%s) = %v, want %v", d.text, got, d.want)
		}
	}
}

func Test_DC_Query_By_File_Id(t *testing.T) {
	var data = []struct {
		text string
		want int
	}{
		{"AgACAgEAAxUAAWJGy2LS9jhZRw6R2qgXlZqrRo5oAAKqpzEbUr8gNO_REKKp8zC8AQADAgADYwADIwQ", 1},
		{"AgACAgIAAxUAAWJGzJ9B4LKaPV5cjXOdO0xiQb-7AAI8uTEby3YxSvsJRGGMR2jcAQADAgADYwADIwQ", 2},
		{"AgACAgMAAxUAAWJGza3gyRKTX-Eg5BN_n0iMHpxNAAKppzEbshEJBC3NV1bsBwi3AQADAgADYwADIwQ", 3},
		{"AgACAgQAAxUAAWJGyzZFag7qs544PxHGvabVQ6wUAAIZuTEbcDIYUqVQrWEZVR2JAQADAgADYwADIwQ", 4},
		{"BQACAgUAAx0CVTVl8wACQbJiRlY8fOQVb8cEg_eiBU3A3DZ6ZwACwgYAAu5bOFaKGWHnjrb_biME", 5},
		{"", -1},
	}

	for _, d := range data {
		if got := dcQueryByFileId(d.text); d.want != got {
			t.Errorf("dcQuery(%s) = %v, want %v", d.text, got, d.want)
		}
	}
}
