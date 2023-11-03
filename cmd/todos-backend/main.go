package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sankalpmukim/todos-backend/internal/initialize"
	"github.com/sankalpmukim/todos-backend/internal/middleware"
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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(middleware.SupabaseTokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(middleware.CreateUserIfNotFound)
		r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := jwtauth.FromContext(r.Context())
			// w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["email"])))
			if err != nil {
				logs.Error("Error getting claims from context", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			email := claims["email"].(string)

			if email == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "Unauthorized",
				})
				return
			}

			response := map[string]interface{}{
				"email": email,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})
	})

	// Listen on port 3000
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}
	logs.Info("Starting server on port " + PORT)
	http.ListenAndServe(":"+PORT, r)
}
