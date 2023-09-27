package router

import (
	"bullshape/confs"
	"bullshape/models"
	"bullshape/utils"
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

func (s *Server) jwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log := s.Logger
		for endpoint, method := range NotAuthEndPoint {
			if method == r.Method && endpoint == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tk := &models.Token{}
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			log.Error("No cookie... Unauthorized")
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
				log.Error("Missing auth token")
				w.WriteHeader(http.StatusForbidden)
				u.HttpError(w, http.StatusInternalServerError, errors.New("Missing auth token"))
				return
			}
			//The token normally comes in format `Bearer {token-body}`
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				log.Error("Invalid/Malformed auth token")
				u.HttpError(w, http.StatusForbidden, errors.New("Invalid/Malformed auth token"))
				return
			}

			tokenPart := splitted[1]

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})

			if err != nil {
				log.Error("Malformed authentication token")
				u.HttpError(w, http.StatusForbidden, errors.New("Malformed authentication token"))
				return
			}

			if !token.Valid {
				log.Error("Token is not valid")
				u.HttpError(w, http.StatusForbidden, errors.New("Token is not valid"))
				return
			}
		}

		log.Info("User ", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.Logger.Error("Recover. Error: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				u.HttpError(w, http.StatusInternalServerError, errors.New("System hanged"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var HttpLogger = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := "req-" + utils.NewUUIDV4()
		logUser := ""
		user := r.Context().Value("user")
		if user != nil {
			logUser = "user(ID)=" + fmt.Sprint(user)
		}
		r.Header.Add("request_id", reqID)
		r.Header.Add("logged_user", logUser)

		next.ServeHTTP(w, r)

	})
}
