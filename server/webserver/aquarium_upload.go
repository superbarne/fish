package webserver

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/superbarne/fish/imageprocess"
	"github.com/superbarne/fish/models"
)

func (ws *WebServer) uploadAquariumFish(w http.ResponseWriter, r *http.Request) {
	// validate id
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// find aquarium
	aquarium, err := ws.storage.Aquarium(aquariumID)
	if err != nil {
		ws.log.Error("Failed to get aquarium", slog.String("error", err.Error()))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Get the file from the request
		file, multipartHeader, err := r.FormFile("image")
		if err != nil {
			ws.log.Error("Failed to get image from request", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}
		defer file.Close()

		// is file a image
		if multipartHeader.Header.Get("Content-Type") != "image/png" && multipartHeader.Header.Get("Content-Type") != "image/jpeg" && multipartHeader.Header.Get("Content-Type") != "image/jpg" {
			ws.log.Error("File is not a image", slog.String("content-type", multipartHeader.Header.Get("Content-Type")))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		fishID := uuid.New()
		tmpFilePath, err := ws.storage.SaveTmpFishImageFromRequest(aquarium.ID, fishID, file, multipartHeader)
		if err != nil {
			ws.log.Error("Failed to save image", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		defer func() {
			// Remove tmp file
			if err := os.Remove(tmpFilePath); err != nil {
				ws.log.Error("Failed to remove temp file", slog.String("error", err.Error()))
			}
		}()

		// Process Image
		targetPath, err := ws.storage.FishImagePath(aquarium.ID, fishID)
		if err != nil {
			ws.log.Error("Failed to get fish image path", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}
		if err := imageprocess.ProcessImage(tmpFilePath, targetPath, ws.log); err != nil {
			ws.log.Error("Failed to process image", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		name := r.FormValue("name")
		if name == "" {
			name = "Boid"
		}

		// Write Json with metadata about the uploaded file
		fish := &models.Fish{
			ID:         fishID,
			AquariumID: aquarium.ID,
			Name:       name,
			Filename:   fishID.String() + ".png",
			Approved:   !aquarium.NeedApproval, // if need approval true, set fish approved value to false
		}

		if err := ws.storage.InsertFish(aquariumID, fish); err != nil {
			ws.log.Error("Failed to save fish", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		ws.pubsub.Publish("aquarium:"+aquariumID.String(), fish)

		http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
		return
	}

	ws.tmpl.ExecuteTemplate(w, "upload.html", map[string]interface{}{
		"ID":       aquarium.ID.String(),
		"Revision": ws.gitCommit,
	})
}
