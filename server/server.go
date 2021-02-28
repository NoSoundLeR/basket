package server

import (
	"flag"
	"log"
	"net/http"

	"github.com/NoSoundLeR/basket.git/server/basket"
	"github.com/NoSoundLeR/basket.git/server/handler"
	"github.com/NoSoundLeR/basket.git/server/hub"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var bindAddr = flag.String("bind", "127.0.0.1:8080", "http service address")

// Server ...
type Server struct {
	bs *basket.Baskets
}

// NewServer ...
func NewServer() *Server {
	return &Server{
		bs: &basket.Baskets{
			Baskets: make([]*basket.Basket, 0, 50),
		},
	}
}

// Run ...
func (s *Server) Run() {
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hub := hub.NewHub()
	go hub.Run()

	r.Get("/", handler.BasketIndex(s.bs))
	r.Get("/{id}", handler.BasketGet(s.bs))
	r.Post("/api/baskets", handler.APIBasketCreate(s.bs))
	r.Put("/api/baskets/{id}", handler.APIBasketThrow(hub, s.bs))
	r.Post("/api/baskets/{id}/close", handler.APIBasketClose(hub, s.bs))
	r.HandleFunc("/ws/{id}", handler.BasketWS(hub, s.bs))

	err := http.ListenAndServe(*bindAddr, r)
	if err != nil {
		log.Fatal("ListenAndServer", err)
	}
}
