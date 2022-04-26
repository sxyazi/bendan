package vercel

import (
	"github.com/sxyazi/bendan/boot"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	boot.ServeHook(w, r)
}
