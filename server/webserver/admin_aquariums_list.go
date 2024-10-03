package webserver

import (
	"net/http"
)

func (ws *WebServer) listAdminAquariums(w http.ResponseWriter, r *http.Request) {
	aquariums, err := ws.storage.Aquariums()
	if err != nil {
		http.Error(w, "Failed to get aquariums", http.StatusInternalServerError)
		return
	}

	ws.tmpl.ExecuteTemplate(w, "admin_aquariums.html", map[string]interface{}{
		"Aquariums": aquariums,
	})
}
