package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	CreateAllRoutes()
}

type handler struct {
	*chi.Mux
}

func NewHandler() Handler {
	return &handler{
		Mux: chi.NewMux(),
	}
}

func (h *handler) CreateAllRoutes() {
	h.Post("/in", post())
}

func post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("This should be the end!")

		w.WriteHeader(200)
	}
}
