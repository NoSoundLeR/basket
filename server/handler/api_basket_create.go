package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/hub"
)

// APIBasketCreate ...
func APIBasketCreate(h *hub.Hub, b basket.Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		var d basket.Data
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		basket, err := basket.NewBasket(ip, d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b.AddBasket(*basket)
		res := privateInfo{
			ID: basket.ID.String(),
		}
		http.SetCookie(w, &http.Cookie{
			Name:   basket.ID.String(),
			Value:  basket.Secret,
			Path:   "/",
			MaxAge: basket.Timeout,
		})

		go func() {
			select {
			case <-time.After(time.Duration(basket.Timeout) * time.Second):
				if basket.Active {
					basket.Close()
					h.BroadcastResult(basket.ID.String(), basket.Result)
				}
			}
		}()

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
