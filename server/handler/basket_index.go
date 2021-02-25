package handler

import (
	"net/http"

	"github.com/NoSoundLeR/basket.git/server/basket"
)

// BasketIndex ...
func BasketIndex(b basket.Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}
}
