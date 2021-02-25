package handler

import (
	"net/http"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/hub"
)

// APIBasketEnd ...
func APIBasketEnd(h *hub.Hub, s basket.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
