package routes

import (
	"net/http"
	"olhar-biblico/auth"
	"olhar-biblico/controllers"
	"olhar-biblico/services"

	"github.com/go-chi/chi"
)

func UsersRouter() http.Handler {
	r := chi.NewRouter()

	service := services.NewUserService()
	handler := controllers.NewUserHandler(service)

	r.Post("/login", handler.Login)
	r.Post("/create", auth.Auth(handler.Create))
	r.Put("/change-password", auth.Auth(handler.UpdatePassword))

	return r
}
