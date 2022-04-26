package vercel

import (
	"github.com/sxyazi/bendan/boot"
	. "github.com/sxyazi/bendan/utils"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != Config("refresh_key") {
		w.WriteHeader(http.StatusFound)
		w.Header().Set("Location", "https://github.com/sxyazi/bendan")
		return
	}

	boot.ServeHook(w, r)
}
