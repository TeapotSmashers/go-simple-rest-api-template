package auth

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

var SupabaseTokenAuth *jwtauth.JWTAuth

func Initialize() {
	SECRET := os.Getenv("SUPABASE_SECRET")
	SupabaseTokenAuth = jwtauth.New("HS256", []byte(SECRET), nil)
}

// type JWTUser struct {
// 	Email string `json:"email"`
// 	// Add other fields if necessary
// }
// func IsValidSupabaseToken(token, url, apikey string) (*JWTUser, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url+"/auth/v1/user", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	req.Header.Add("apikey", apikey)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == 200 {
// 		bodyBytes, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}

// 		var user JWTUser
// 		err = json.Unmarshal(bodyBytes, &user)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &user, nil
// 	}

// 	return nil, fmt.Errorf("invalid token or other error occurred")
// }
