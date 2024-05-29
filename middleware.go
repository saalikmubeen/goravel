package goravel

import "net/http"

func (g *Goravel) SessionLoad(next http.Handler) http.Handler {
	g.InfoLog.Println("Session middleware loaded")
	return g.Session.LoadAndSave(next)
}
