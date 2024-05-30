package middleware

import "net/http"

// Web authentication
func (m *Middleware) Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !m.App.Session.Exists((r.Context()), "userID") {
			http.Error(w, http.StatusText(401), http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func (m *Middleware) AuthToken(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := m.Models.Tokens.AuthenticateToken(r)

		if err != nil {
			var payload struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}

			payload.Error = true
			payload.Message = "Invalid authentication credentials"

			// send the response
			m.App.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"data": payload,
			})

		} else {
			next.ServeHTTP(w, r)
		}

	})
}
