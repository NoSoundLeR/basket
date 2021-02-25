package handler

import (
	"net/http"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/hub"
	"github.com/go-chi/chi"
)

// APIBasketClose ...
func APIBasketClose(h *hub.Hub, s basket.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		b, err := s.GetBasket(id)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		secretCookie, err := r.Cookie(id)
		if err != nil || secretCookie.Value != b.Secret {
			http.Error(w, http.StatusText(403), 403)
			return
		}
		b.Close()
		h.BroadcastResult(id, b.Result)
	}
}
