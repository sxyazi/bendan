package commands

import (
	"testing"
)

func Test_dataCenterBy(t *testing.T) {
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
		if got := dataCenterBy(d.text); d.want != got {
			t.Errorf("dataCenterBy(%s) = %v, want %v", d.text, got, d.want)
		}
	}
}
