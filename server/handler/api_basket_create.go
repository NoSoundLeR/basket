package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NoSoundLeR/basket.git/server/basket"
)

type data struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

// APIBasketCreate ...
func APIBasketCreate(b basket.Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		var d data
		decoder := json.NewDecoder(r.Body)
		// if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		err := decoder.Decode(&d)
		basket, err := basket.NewBasket(ip, d.Title, d.Description, d.Value)
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
			MaxAge: 60,
		})
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
