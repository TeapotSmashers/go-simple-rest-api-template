package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sankalpmukim/todos-backend/pkg/logs"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Hello World!",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.Error("Error encoding response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func HealthZ(w http.ResponseWriter, r *http.Request) {
	// return status: ok and server_timestamp
	response := map[string]interface{}{
		"status":           "ok",
		"server_timestamp": time.Now().Unix(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.Error("Error encoding response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ReturnMail(w http.ResponseWriter, r *http.Request) {
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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logs.Error("Error encoding response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
