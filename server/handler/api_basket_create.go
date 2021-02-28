package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/hub"
)

type data struct {
	ProtectionLevel string `json:"protectionLevel"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Value           string `json:"value"`
}

const timeout = 600

// APIBasketCreate ...
func APIBasketCreate(h *hub.Hub, b basket.Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		var d data
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		protectionLevel, err := strconv.Atoi(d.ProtectionLevel)
		if err != nil {
			http.Error(w, "wrong protection level", 400)
			return
		}
		basket, err := basket.NewBasket(ip, protectionLevel, d.Title, d.Description, d.Value)
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
			MaxAge: timeout,
		})

		go func() {
			select {
			case <-time.After(timeout * time.Second):
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
