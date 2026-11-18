package routes

import (
	"net/http"
	"olhar-biblico/auth"
	"olhar-biblico/controllers"
	"olhar-biblico/services"

	"github.com/go-chi/chi"
)

func VisitRequestsRouter() http.Handler {
	r := chi.NewRouter()

	service := services.NewVisitRequestService()
	handler := controllers.NewVisitRequestHandler(service)

	r.Post("/create", handler.Create)
	r.Get("/", auth.Auth(handler.Find))

	return r
}
