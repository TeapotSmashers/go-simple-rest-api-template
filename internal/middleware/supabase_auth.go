package middleware

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

var SupabaseTokenAuth *jwtauth.JWTAuth

func Initialize() {
	SECRET := os.Getenv("SUPABASE_SECRET")
	SupabaseTokenAuth = jwtauth.New("HS256", []byte(SECRET), nil)
}
