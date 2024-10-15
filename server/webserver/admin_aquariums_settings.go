package webserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (ws *WebServer) toggleAdminNeedApproval(w http.ResponseWriter, r *http.Request) {
	// validate id
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	// find aquarium
	aquarium, err := ws.storage.Aquarium(aquariumID)
	if err != nil {
		ws.log.Error("Failed to get aquarium", slog.String("error", err.Error()))
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	aquarium.NeedApproval = !aquarium.NeedApproval

	if err := ws.storage.InsertAquarium(aquarium); err != nil {
		ws.log.Error("Failed to save aquarium", slog.String("error", err.Error()))
		http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
}
