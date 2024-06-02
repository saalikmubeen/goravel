package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/saalikmubeen/goravel-demo-app/models"
)

// CheckRemember checks for a remember_me token in the request and logs the user in if it's valid
func (m *Middleware) CheckRememberMe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !m.App.Session.Exists(r.Context(), "userID") { // user is not logged in (no session)

			// check for the remember_me cookie
			cookie, err := r.Cookie(fmt.Sprintf("_%s_remember_me", m.App.AppName))

			if err != nil { // no remember_me cookie, so on to the next middleware
				// no cookie, so on to the next middleware
				next.ServeHTTP(w, r)
			} else {
				// we found a remember_me cookie, so check it
				key := cookie.Value // the value of the cookie (it has two parts: user_id|remember_me_token)
				var u models.User

				if len(key) > 0 { // cookie contains data, so validate it

					split := strings.Split(key, "|")
					uid, hash := split[0], split[1]
					id, _ := strconv.Atoi(uid)

					// check if the hash(remember_me_token) is valid and exists in the database
					validHash := u.CheckForRememberToken(id, hash)
					if !validHash {

						// if the hash is not valid, delete the remember_me cookie
						m.deleteRememberMeCookie(w, r)
						m.App.Session.Put(r.Context(), "error", "You've been logged out from another device")
						next.ServeHTTP(w, r)
					} else {

						// valid hash, so log the user in.
						user, _ := u.Get(id)
						m.App.Session.Put(r.Context(), "userID", user.ID)
						m.App.Session.Put(r.Context(), "remember_me_token", hash)
						next.ServeHTTP(w, r)
					}
				} else {
					// key length is zero, so it's probably a leftover cookie (user has not closed browser)

					// delete the remember_me cookie
					m.deleteRememberMeCookie(w, r)
					next.ServeHTTP(w, r)
				}

			}
		} else {
			// user is logged in
			next.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) deleteRememberMeCookie(w http.ResponseWriter, r *http.Request) {

	// renew the session token
	_ = m.App.Session.RenewToken(r.Context())
	// delete the cookie
	newCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_remember_me", m.App.AppName),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-100 * time.Hour),
		HttpOnly: true,
		Domain:   m.App.Session.Cookie.Domain,
		MaxAge:   -1,
		Secure:   m.App.Session.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &newCookie)

	// log the user out
	m.App.Session.Remove(r.Context(), "userID")
	m.App.Session.Remove(r.Context(), "remember_me_token")
	m.App.Session.Destroy(r.Context())

	_ = m.App.Session.RenewToken(r.Context())
}
