package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/* returns http.Handler*/
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Post("/api/v1/signup", app.Signup)
	mux.Post("/api/v1/login", app.Login)
	mux.Get("/api/v1/all-users", app.GetUsers)
	mux.Get("/api/v1/verify-user-token", app.HandleVerifyToken)

	// Add the Prometheus metrics endpoint to the router
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
