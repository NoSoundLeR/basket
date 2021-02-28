package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/hub"
	"github.com/go-chi/chi"
)

type newVariant struct {
	Value string `json:"value"`
}

// APIBasketThrow ...
func APIBasketThrow(h *hub.Hub, b basket.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ip := strings.Split(r.RemoteAddr, ":")[0]
		var v newVariant
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		basket, err := b.GetBasket(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		prevLen := basket.GetSize()
		if ok := basket.Add(ip, v.Value); !ok {
			http.Error(w, "could not add value", http.StatusInternalServerError)
			return
		}
		curLen := basket.GetSize()
		if curLen > prevLen {
			h.BroadcastThrow(id, strconv.Itoa(curLen))
		}
	}
}
