package webserver

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/superbarne/fish/aquarium"
	"github.com/superbarne/fish/imageprocess"
	"github.com/superbarne/fish/models"
)

func (ws *WebServer) uploadFish(w http.ResponseWriter, r *http.Request) {
	// validate id
	aquariumID, err := uuid.Parse(chi.URLParam(r, "aquariumID"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var aquarium *aquarium.Aquarium
	var ok bool
	if err := func() error {
		ws.aquariumsLock.RLock()
		defer ws.aquariumsLock.RUnlock()
		aquarium, ok = ws.aquariums[aquariumID]
		if !ok {
			return errors.New("aquarium not found")
		}
		return nil
	}(); err != nil {
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

		name := r.FormValue("name")
		if name == "" {
			name = "Boid"
		}
		fishID := uuid.New()

		tmpFilePath, err := ws.storage.SaveTmpFishImageFromRequest(aquarium.ID, fishID, file, multipartHeader)
		if err != nil {
			ws.log.Error("Failed to save image", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		// Process Image
		targetPath := ws.storage.FishImagePath(aquarium.ID, fishID)
		if err := imageprocess.ProcessImage(tmpFilePath, targetPath, ws.log); err != nil {
			ws.log.Error("Failed to process image", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		// Remove tmp file
		if err := os.Remove(tmpFilePath); err != nil {
			ws.log.Error("Failed to remove temp file", slog.String("error", err.Error()))
		}

		// Write Json with metadata about the uploaded file
		fish := &models.Fish{
			ID:         fishID,
			AquariumID: aquarium.ID,
			Name:       name,
			Filename:   fishID.String() + ".png",
			Approved:   false,
		}

		if err := ws.storage.SaveFishMetadata(aquariumID, fish); err != nil {
			ws.log.Error("Failed to save fish metadata", slog.String("error", err.Error()))
			http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
			return
		}

		aquarium.AddFish(fish)

		http.Redirect(w, r, "/aquarium/"+aquariumID.String(), http.StatusSeeOther)
		return
	}

	ws.tmpl.ExecuteTemplate(w, "upload.html", map[string]interface{}{
		"ID": aquariumID.String(),
	})
}
