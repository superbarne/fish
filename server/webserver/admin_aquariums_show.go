package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (ws *WebServer) showAdminAquarium(w http.ResponseWriter, r *http.Request) {
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	aquarium, err := ws.storage.Aquarium(aquariumID)
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	fishes, err := ws.storage.Fishes(aquariumID)
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	ws.tmpl.ExecuteTemplate(w, "admin_aquarium.html", map[string]interface{}{
		"Aquarium": aquarium,
		"Fishes":   fishes,
		"Revision": ws.gitCommit,
	})
}
