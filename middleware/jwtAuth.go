package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/heroku/whaler-api/models"
	"github.com/heroku/whaler-api/utils"
)

var userIDCtxKey = &contextKey{"userID"}

type contextKey struct {
	name string
}

//DEPRECATED -- REST
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/create",
			"/api/user/login",
			"/api/org/create",
			"/api/account/create",
			"/api/contact/create",
			"/schema",
			"/socket"} //temporary! not for release!
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = utils.Message(4003, "Missing auth token", true, map[string]interface{}{})
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, ".")
		if len(splitted) != 3 {
			response = utils.Message(4003, "Invalid/Malformed auth token", true, map[string]interface{}{})
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		tk := &models.AccessToken{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		isNotRefreshEndpoint := r.URL.Path != "/api/user/refresh"

		if (err != nil || !token.Valid) && isNotRefreshEndpoint {
			fmt.Print(fmt.Sprint(err))
			data := map[string]interface{}{"error": err}
			response = utils.Message(4003, "Invalid authentication token", true, data)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		fmt.Println(fmt.Sprintf("User %s", tk.UserID))
		ctx := context.WithValue(r.Context(), userIDCtxKey, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

var ParseUserIDFromToken = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		tk := &models.AccessToken{}

		jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		fmt.Println(fmt.Sprintf("User %d", tk.UserID))
		ctx := context.WithValue(r.Context(), userIDCtxKey, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(userIDCtxKey).(string)
	return id
}
