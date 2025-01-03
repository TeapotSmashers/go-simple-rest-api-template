package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sankalpmukim/todos-backend/internal/database"
	"github.com/sankalpmukim/todos-backend/pkg/logs"
)

func CreateUserIfNotFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			logs.Error("Error getting claims from context", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		email := claims["email"].(string)
		exists, err := database.DB.UserExists(email)
		if err != nil {
			logs.Error("Error checking if user exists", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exists {
			user := database.User{
				ID:       claims["sub"].(string),
				Name:     claims["user_metadata"].(map[string]interface{})["full_name"].(string),
				Email:    claims["email"].(string),
				Provider: claims["app_metadata"].(map[string]interface{})["provider"].(string),
			}
			if err := database.DB.CreateUser(user); err != nil {
				logs.Error("Error creating user", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			logs.Info("Created user", user)
		}
		next.ServeHTTP(w, r)
	})
}
