package router

import (
	"bullshape/confs"
	"bullshape/models"
	u "bullshape/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//List of endpoints that doesn't require auth
var NotAuthEndPoint = map[string]string{
	"/api/companies/": "GET",
	"/api/user":       "POST",
	"/api/user/login": "POST",
}

func jwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for endpoint, method := range NotAuthEndPoint {
			if method == r.Method && endpoint == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tk := &models.Token{}
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "No cookie... Unauthorized", http.StatusUnauthorized)
			return
		}

		if cookie != nil {
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(confs.TokenPass), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		} else if true {
			tokenHeader := r.Header.Get("Authorization")

			if tokenHeader == "" {
				w.WriteHeader(http.StatusForbidden)
				u.HttpError(w, http.StatusInternalServerError, errors.New("Missing auth token"))
				return
			}
			//The token normally comes in format `Bearer {token-body}`
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				u.HttpError(w, http.StatusForbidden, errors.New("Invalid/Malformed auth token"))
				return
			}

			tokenPart := splitted[1]

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})

			if err != nil {
				u.HttpError(w, http.StatusForbidden, errors.New("Malformed authentication token"))
				return
			}

			if !token.Valid {
				u.HttpError(w, http.StatusForbidden, errors.New("Token is not valid"))
				return
			}
		}

		fmt.Println("User ", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
