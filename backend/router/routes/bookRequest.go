package routes

import (
	"net/http"
	"olhar-biblico/auth"
	"olhar-biblico/controllers"
	"olhar-biblico/services"

	"github.com/go-chi/chi"
)

func BookRequestsRouter() http.Handler {
	r := chi.NewRouter()

	service := services.NewBookRequestService()
	handler := controllers.NewBookRequestHandler(service)

	r.Post("/create", handler.Create)
	r.Get("/", auth.Auth(handler.Find))

	return r
}
