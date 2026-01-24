package swaggerAuth

import (
	"crypto/subtle"
	"net/http"
)

// * middleware для защиты Swagger
func New(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Если credentials пустые, Swagger недоступен
			if username == "" || password == "" {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			user, pass, ok := r.BasicAuth()

			usernameMatch := subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1
			passwordMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1

			if !ok || !usernameMatch || !passwordMatch {
				w.Header().Set("WWW-Authenticate", `Basic realm="Swagger Documentation"`)
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Unauthorized"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
