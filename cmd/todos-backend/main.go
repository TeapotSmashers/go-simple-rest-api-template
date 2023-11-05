package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sankalpmukim/todos-backend/internal/handlers"
	"github.com/sankalpmukim/todos-backend/internal/initialize"
	"github.com/sankalpmukim/todos-backend/internal/middleware"
	"github.com/sankalpmukim/todos-backend/internal/routers"
	"github.com/sankalpmukim/todos-backend/pkg/logs"
)

func main() {
	if err := initialize.InitAll(); err != nil {
		fmt.Println("Error occured during initialization", err)
		panic(err)
	}

	// Get the value of the DEBUG environment variable
	DEBUG := os.Getenv("DEBUG")
	if DEBUG != "true" {
		// cannot use logs package here because it
		// doesn't print to the console.
		fmt.Printf("DEBUG = %v\n", DEBUG)
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/", handlers.HelloWorld)
	r.Get("/healthz", handlers.HealthZ)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(middleware.SupabaseTokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(middleware.CreateUserIfNotFound)
		r.Get("/auth", handlers.ReturnMail)
		r.Mount("/todos", routers.Todos)
	})

	// Listen on port 3000
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	logs.Info("Starting server on port " + PORT)
	http.ListenAndServe(":"+PORT, r)
}
