package handler

import (
	"log"
	"net/http"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/hub"
	"github.com/go-chi/chi"
)

// BasketWS ...
func BasketWS(h *hub.Hub, s basket.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		b, err := s.GetBasket(id)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		if !b.Active {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		conn, err := hub.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := hub.NewClient(h, conn, id)
		client.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}
