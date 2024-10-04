package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (ws *WebServer) deleteAdminFish(w http.ResponseWriter, r *http.Request) {
	// validate id
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	fishID, err := uuid.Parse(chi.URLParam(r, "fishID"))
	if err != nil {
		http.Redirect(w, r, "/admin/aquarium/"+aquariumID.String(), http.StatusSeeOther)
		return
	}

	// find aquarium
	aquarium, err := ws.storage.Aquarium(aquariumID)
	if err != nil {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	// find fish
	fish, err := ws.storage.Fish(aquarium.ID, fishID)
	if err != nil {
		http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
		return
	}

	// delete fish
	err = ws.storage.DeleteFish(aquarium.ID, fishID)
	if err != nil {
		http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
		return
	}

	// pubsub
	ws.pubsub.Publish("aquarium:"+aquariumID.String()+":delete", fish)

	http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
}
