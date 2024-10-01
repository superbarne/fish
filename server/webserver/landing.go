package webserver

import (
	"net/http"
)

func (ws *WebServer) getLandingPage(w http.ResponseWriter, r *http.Request) {
	ws.tmpl.ExecuteTemplate(w, "landing.html", nil)
}
