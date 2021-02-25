package handler

import (
	"html/template"
	"net/http"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/go-chi/chi"
)

// BasketGet ...
func BasketGet(s basket.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		basket, err := s.GetBasket(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var isOwner bool
		secretCookie, err := r.Cookie(id)
		if err == nil {
			isOwner = secretCookie.Value == basket.Secret
		}

		res := publicInfo{
			ID:          basket.ID.String(),
			IsOwner:     isOwner,
			Active:      basket.Active,
			Title:       basket.Title,
			Description: basket.Description,
			Result:      basket.Result,
			Count:       len(basket.Vars),
		}

		t, err := template.ParseFiles("./public/id.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, res)
	}
}
