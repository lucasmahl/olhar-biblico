package router

import (
	"net/http"
	"olhar-biblico/router/routes"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Config() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Use o middleware Throttle para limitar o tráfego para x requisições por segundo
	r.Use(middleware.Throttle(50))
	// Use o middleware Timeout para cancelar a requisição se ela demorar mais do que x segundos
	r.Use(middleware.Timeout(120 * time.Second))

	// Adicione o middleware tollboothLimitHandler
	r.Use(tollboothLimitHandler)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Mount("/visit-requests", routes.VisitRequestsRouter())
	r.Mount("/book-requests", routes.BookRequestsRouter())
	r.Mount("/users", routes.UsersRouter())

	return r
}

// Função que limita a quantidade de acessos por segundo, por IP
func tollboothLimitHandler(next http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(50, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second}) //50 acessos por segundo

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpError := tollbooth.LimitByKeys(lmt, []string{r.RemoteAddr})
		if httpError != nil {
			http.Error(w, httpError.Message, httpError.StatusCode)
			return
		}
		next.ServeHTTP(w, r)
	})
}
