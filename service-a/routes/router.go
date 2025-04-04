package routes

import (
	"net/http"
	"service-a/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/cep", handlers.HandleCEP)
	return r
}
