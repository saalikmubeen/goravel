package goravel

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (g *Goravel) SessionLoad(next http.Handler) http.Handler {
	g.InfoLog.Println("Session middleware loaded")
	return g.Session.LoadAndSave(next)
}

func (g *Goravel) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(g.config.cookie.secure)

	// Exempt the API routes from CSRF protection
	// This tells the CSRF middleware to not apply CSRF protection and not to
	// check for the CSRF token being present in the request for the routes that
	// match the pattern.
	// So it won't check to ensure that the request made to this route is a POST request
	// with a form field named "_csrf" containing the CSRF token.
	// For api routes, we manually check and verify the CSRF token in the handler function.
	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   g.config.cookie.domain,
	})
	return csrfHandler
}
