package webserver

import (
	"log/slog"
	"net/http"
)

func (ws *WebServer) listAdminAquariums(w http.ResponseWriter, r *http.Request) {
	aquariums, err := ws.storage.Aquariums()
	if err != nil {
		ws.log.Error("Failed to get aquariums", slog.String("error", err.Error()))
		http.Error(w, "Failed to get aquariums", http.StatusInternalServerError)
		return
	}

	ws.tmpl.ExecuteTemplate(w, "admin_aquariums.html", map[string]interface{}{
		"Aquariums": aquariums,
		"Revision":  ws.gitCommit,
	})
}
