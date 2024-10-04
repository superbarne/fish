package webserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (ws *WebServer) approveAdminFish(w http.ResponseWriter, r *http.Request) {
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

	// no toggle! if two people approve at the same time, it will be approved
	approved := r.FormValue("approved") == "true"
	fish.Approved = approved
	if fish.Approved {
		now := time.Now()
		fish.ApprovedAt = &now
	}

	// save
	if err := ws.storage.InsertFish(aquarium.ID, fish); err != nil {
		http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
		return
	}

	if fish.Approved {
		// publish
		ws.pubsub.Publish("aquarium:"+aquariumID.String(), fish)
	} else {
		// delete fish from aquarium
		ws.pubsub.Publish("aquarium:"+aquariumID.String()+":delete", fish)
	}

	http.Redirect(w, r, "/admin/aquarium/"+aquarium.ID.String(), http.StatusSeeOther)
}
